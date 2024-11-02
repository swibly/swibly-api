package repository

import (
	"encoding/json"
	"errors"

	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/pagination"
	"github.com/swibly/swibly-api/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type componentRepository struct {
	db *gorm.DB

	userRepo UserRepository
}

type ComponentRepository interface {
	Create(createModel *dto.ComponentCreation) error
	Update(componentID uint, updateModel *dto.ComponentUpdate) error

	Get(issuerID uint, componentModel *model.Component) (*dto.ComponentInfo, error)
	GetPublic(issuerID uint, page, perPage int, freeOnly bool) (*dto.Pagination[dto.ComponentInfo], error)
	GetOwned(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error)
	GetByOwnerID(issuerID, ownerID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error)
	GetHoldersByID(issuerID, userID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error)
	GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error)

	Search(issuerID uint, search *dto.SearchComponent, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error)

	Buy(issuerID, componentID uint) error
	Sell(issuerID, componentID uint) error

	SafeDelete(componentID uint) error
	Restore(componentID uint) error
	UnsafeDelete(componentID uint) error
	ClearTrash(userID uint) error
}

var (
	ErrInsufficientArkhoins     = errors.New("insufficient arkhoins")
	ErrComponentNotFound        = errors.New("component not found or trashed")
	ErrComponentNotTrashed      = errors.New("component is not trashed")
	ErrComponentAlreadyTrashed  = errors.New("component is already trashed")
	ErrComponentNotPublic       = errors.New("component is not public")
	ErrComponentAlreadyPublic   = errors.New("component is already public")
	ErrComponentNotOwned        = errors.New("component is not owned")
	ErrComponentAlreadyOwned    = errors.New("component is already owned")
	ErrComponentOwnerCannotBuy  = errors.New("component owner cannot buy their own component")
	ErrComponentOwnerCannotSell = errors.New("component owner cannot sell their own component")
)

func NewComponentRepository(userRepo UserRepository) ComponentRepository {
	return &componentRepository{db.Postgres, userRepo}
}

func (cr *componentRepository) baseComponentQuery(issuerID uint) *gorm.DB {
	return cr.db.Table("components c").
		Select(`
			c.id as id,
			c.created_at as created_at,
			c.updated_at as updated_at,
			c.deleted_at as deleted_at,
			c.name as name,
			c.description as description,
			c.content as content,
			c.price as price,
      c.budget as budget,
			co.id AS owner_id,
			u.id AS owner_id,
			u.first_name AS owner_first_name,
			u.last_name AS owner_last_name,
			u.username AS owner_username,
      u.profile_picture AS owner_profile_picture,
			u.verified AS owner_verified,
			COALESCE((
				SELECT COUNT(*)
				FROM component_holders ch
				WHERE ch.component_id = c.id
			), 0) AS holders,
      COALESCE((
				SELECT SUM(ch.price_paid)
				FROM component_holders ch
				WHERE ch.component_id = c.id
			), 0) AS total_sells,
			EXISTS (
				SELECT 1 
				FROM component_publications cp 
				WHERE cp.component_id = c.id
			) AS is_public,
			EXISTS (
				SELECT 1
				FROM component_holders ch
				WHERE ch.component_id = c.id AND ch.user_id = ?
			) AS bought,
			(SELECT ch.price_paid FROM component_holders ch WHERE ch.component_id = c.id AND ch.user_id = ?) AS paid_price,
      (SELECT ch.price_paid FROM component_holders ch WHERE ch.component_id = c.id AND ch.user_id = ?) AS sell_price
		`, issuerID, issuerID, issuerID).
		Joins("JOIN component_owners co ON co.component_id = c.id").
		Joins("JOIN users u ON co.user_id = u.id")
}

func (cr *componentRepository) paginateComponents(query *gorm.DB, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	paginationResult, err := pagination.Generate[dto.ComponentInfoJSON](query, page, perPage)
	if err != nil {
		return nil, err
	}

	componentInfoList := make([]*dto.ComponentInfo, 0, len(paginationResult.Data))
	for _, componentInfoJSON := range paginationResult.Data {
		componentInfo, err := convertToComponentInfo(componentInfoJSON)
		if err != nil {
			return nil, err
		}

		componentInfoList = append(componentInfoList, &componentInfo)
	}

	return &dto.Pagination[dto.ComponentInfo]{
		Data:         componentInfoList,
		TotalRecords: paginationResult.TotalRecords,
		TotalPages:   paginationResult.TotalPages,
		CurrentPage:  paginationResult.CurrentPage,
		NextPage:     paginationResult.NextPage,
		PreviousPage: paginationResult.PreviousPage,
	}, nil
}

func convertToComponentInfo(jsonInfo *dto.ComponentInfoJSON) (dto.ComponentInfo, error) {
	var content any
	err := json.Unmarshal(jsonInfo.Content, &content)
	if err != nil {
		return dto.ComponentInfo{}, err
	}

	return dto.ComponentInfo{
		ID:                  jsonInfo.ID,
		CreatedAt:           jsonInfo.CreatedAt,
		UpdatedAt:           jsonInfo.UpdatedAt,
		DeletedAt:           jsonInfo.DeletedAt,
		Name:                jsonInfo.Name,
		Description:         jsonInfo.Description,
		Content:             content,
		Budget:              jsonInfo.Budget,
		Price:               jsonInfo.Price,
		PaidPrice:           jsonInfo.PaidPrice,
		SellPrice:           jsonInfo.SellPrice,
		OwnerID:             jsonInfo.OwnerID,
		OwnerFirstName:      jsonInfo.OwnerFirstName,
		OwnerLastName:       jsonInfo.OwnerLastName,
		OwnerUsername:       jsonInfo.OwnerUsername,
		OwnerProfilePicture: jsonInfo.OwnerProfilePicture,
		OwnerVerified:       jsonInfo.OwnerVerified,
		IsPublic:            jsonInfo.IsPublic,
		Holders:             jsonInfo.Holders,
		Bought:              jsonInfo.Bought,
		TotalSells:          jsonInfo.TotalSells,
	}, nil
}

func refund(tx *gorm.DB, componentID uint) error {
	var holders []model.ComponentHolder
	if err := tx.Where("component_id = ?", componentID).Find(&holders).Error; err != nil {
		tx.Rollback()
		return err
	}

	totalRefund := uint64(0)
	for _, holder := range holders {
		var user model.User
		if err := tx.Where("id = ?", holder.UserID).First(&user).Error; err != nil {
			tx.Rollback()
			return err
		}

		user.Arkhoin += uint64(holder.PricePaid)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			return err
		}

		totalRefund += uint64(holder.PricePaid)

		if err := tx.Delete(&holder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	var owner model.User
	if err := tx.Joins("JOIN component_owners co ON co.user_id = users.id").Where("co.component_id = ?", componentID).First(&owner).Error; err != nil {
		tx.Rollback()
		return err
	}

	owner.Arkhoin -= totalRefund
	if owner.Arkhoin < 0 {
		owner.Arkhoin = 0
	}

	if err := tx.Save(&owner).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (cr *componentRepository) Create(createModel *dto.ComponentCreation) error {
	tx := cr.db.Begin()

	contentJSON, err := json.Marshal(createModel.Content)
	if err != nil {
		return err
	}

	component := &model.Component{
		Name:        createModel.Name,
		Description: createModel.Description,
		Content:     string(contentJSON),
		Budget:      createModel.Budget,
		Price:       createModel.Price,
	}

	if err := tx.Create(component).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&model.ComponentOwner{
		ComponentID: component.ID,
		UserID:      &createModel.OwnerID,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if createModel.Public {
		if err := tx.Create(&model.ComponentPublication{ComponentID: component.ID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (cr *componentRepository) Update(componentID uint, updateModel *dto.ComponentUpdate) error {
	if cr.db.Where("id = ?", componentID).First(&model.Component{}).Error == gorm.ErrRecordNotFound {
		return ErrComponentNotFound
	}

	tx := cr.db.Begin()

	if updateModel.Public != nil {
		switch *updateModel.Public {
		case true:
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ComponentPublication{ComponentID: componentID}).Error; err != nil {
				if errors.Is(err, gorm.ErrCheckConstraintViolated) || errors.Is(err, gorm.ErrDuplicatedKey) {
					return ErrComponentAlreadyPublic
				}

				tx.Rollback()
				return err
			}
		case false:
			if err := tx.Where("component_id = ?", componentID).Unscoped().Delete(&model.ComponentPublication{}).Error; err != nil {
				tx.Rollback()
				return err
			}

			refund(tx, componentID)
		}
	}

	updates := make(map[string]interface{})

	if updateModel.Name != nil {
		updates["name"] = *updateModel.Name
	}

	if updateModel.Description != nil {
		updates["description"] = *updateModel.Description
	}

	if updateModel.Content != nil {
		contentJSON, err := json.Marshal(updateModel.Content)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates["content"] = string(contentJSON)
	}

	if updateModel.Budget != nil {
		updates["budget"] = *updateModel.Budget
	}

	if updateModel.Price != nil {
		updates["price"] = *updateModel.Price
	}

	if err := tx.Model(&model.Component{}).Where("id = ?", componentID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (cr *componentRepository) Get(issuerID uint, componentModel *model.Component) (*dto.ComponentInfo, error) {
	var component model.Component
	var componentOwner model.ComponentOwner
	var componentPublication model.ComponentPublication
	var holders []model.ComponentHolder

	if err := cr.db.Unscoped().Where("id = ?", componentModel.ID).First(&component).Error; err != nil {
		return nil, err
	}

	if err := cr.db.Where("component_id = ?", component.ID).First(&componentOwner).Error; err != nil {
		return nil, err
	}

	isPublic := false
	if err := cr.db.Where("component_id = ?", component.ID).First(&componentPublication).Error; err == nil {
		isPublic = true
	}

	if err := cr.db.Where("component_id = ?", component.ID).Find(&holders).Error; err != nil {
		return nil, err
	}

	var owner dto.UserInfoLite
	if componentOwner.UserID != nil {
		ownerProfile, err := cr.userRepo.Get(&model.User{ID: *componentOwner.UserID})
		if err != nil {
			return nil, err
		}

		owner = dto.UserInfoLite{
			ID:             ownerProfile.ID,
			FirstName:      ownerProfile.FirstName,
			LastName:       ownerProfile.LastName,
			Username:       ownerProfile.Username,
			ProfilePicture: ownerProfile.ProfilePicture,
			Verified:       ownerProfile.Verified,
		}
	} else {
		owner = dto.UserInfoLite{
			ID:             0,
			FirstName:      "",
			LastName:       "",
			Username:       "",
			ProfilePicture: "",
			Verified:       false,
		}
	}

	var paidPrice *int
	var sellPrice *int
	bought := false
	var holder model.ComponentHolder
	if err := cr.db.Where("component_id = ? AND user_id = ?", component.ID, issuerID).First(&holder).Error; err == nil {
		paidPrice = &holder.PricePaid
		bought = true
		sellPrice = utils.ToPtr(int(*paidPrice))
	}

	var totalHolders int64
	if err := cr.db.Model(&model.ComponentHolder{}).Where("component_id = ?", component.ID).Count(&totalHolders).Error; err != nil {
		return nil, err
	}

	var content string
	if err := cr.db.Model(&model.Component{}).Unscoped().Select("content").Where("id = ?", component.ID).Scan(&content).Error; err != nil {
		return nil, err
	}

	var contentData any
	if err := json.Unmarshal([]byte(content), &contentData); err != nil {
		return nil, err
	}

	var totalSells int64
	if err := cr.db.Model(&model.ComponentHolder{}).Where("component_id = ?", component.ID).Select("COALESCE(SUM(price_paid), 0)").Scan(&totalSells).Error; err != nil {
		return nil, err
	}

	componentInfo := &dto.ComponentInfo{
		ID:                  component.ID,
		CreatedAt:           component.CreatedAt,
		UpdatedAt:           component.UpdatedAt,
		DeletedAt:           component.DeletedAt,
		Name:                component.Name,
		Description:         component.Description,
		Content:             contentData,
		OwnerID:             owner.ID,
		OwnerFirstName:      owner.FirstName,
		OwnerLastName:       owner.LastName,
		OwnerUsername:       owner.Username,
		OwnerProfilePicture: owner.ProfilePicture,
		OwnerVerified:       owner.Verified,
		Budget:              component.Budget,
		Price:               component.Price,
		PaidPrice:           paidPrice,
		SellPrice:           sellPrice,
		IsPublic:            isPublic,
		Holders:             totalHolders,
		Bought:              bought,
		TotalSells:          totalSells,
	}

	return componentInfo, nil
}

func (cr *componentRepository) GetPublic(issuerID uint, page, perPage int, freeOnly bool) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(issuerID).
		Where("EXISTS (SELECT 1 FROM component_publications cp WHERE cp.component_id = c.id)")

	if freeOnly {
		query = query.Where("c.price = 0")
	}

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(ownerID).
		Unscoped().
		Where("c.deleted_at IS NOT NULL").
		Where(`u.id = ?`, ownerID)

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) GetByOwnerID(issuerID, ownerID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(issuerID).
		Where("co.user_id = ?", ownerID)

	if onlyPublic {
		query = query.Where("EXISTS (SELECT 1 FROM component_publications cp WHERE cp.component_id = c.id)")
	}

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) GetOwned(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(issuerID).
		Joins("LEFT JOIN component_holders ch ON ch.component_id = c.id").
		Joins("LEFT JOIN component_owners coo ON coo.component_id = c.id").
		Where("ch.user_id = ? OR co.user_id = ?", userID, userID)

	if onlyPublic {
		query = query.Where("EXISTS (SELECT 1 FROM component_publications cp WHERE cp.component_id = c.id)")
	}

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) GetHoldersByID(issuerID, componentID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(issuerID).
		Joins("JOIN component_holders ch ON ch.component_id = c.id").
		Where("c.id = ?", componentID)

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) Search(issuerID uint, search *dto.SearchComponent, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	query := cr.baseComponentQuery(issuerID).
		Where("deleted_at IS NULL").
		Joins("JOIN component_publications cp on cp.component_id = c.id")

	orderDirection := "DESC"
	if search.OrderAscending {
		orderDirection = "ASC"
	}

	if search.Name != nil {
		query = query.
			Where(`(
        regexp_like(c.name, ?, 'i') OR
        regexp_like(c.description, ?, 'i') OR
        regexp_like(u.first_name, ?, 'i') OR
        regexp_like(u.last_name, ?, 'i') OR
        regexp_like(u.username, ?, 'i')
      )`,
				utils.RegexPrepareName(*search.Name),
				utils.RegexPrepareName(*search.Name),
				utils.RegexPrepareName(*search.Name),
				utils.RegexPrepareName(*search.Name),
				utils.RegexPrepareName(*search.Name),
			)

		// TODO: Create ranking system
	}

	if search.FollowedUsersOnly {
		query = query.Joins("JOIN followers f ON f.following_id = users.id").
			Where("f.follower_id = ?", issuerID)
	}

	if search.OrderAlphabetic {
		query = query.Order("c.name " + orderDirection)
	} else if search.OrderCreationDate {
		query = query.Order("c.created_at " + orderDirection)
	} else if search.OrderModifiedDate {
		query = query.Order("c.updated_at " + orderDirection)
	} else if search.MostHolders {
		query = query.Order("(SELECT COUNT(*) FROM component_holders ch WHERE ch.component_id = c.id) " + orderDirection)
	} else {
		query = query.Order("c.created_at " + orderDirection)
	}

	return cr.paginateComponents(query, page, perPage)
}

func (cr *componentRepository) Buy(issuerID, componentID uint) error {
	var component model.Component
	var user model.User
	var owner model.User
	var componentOwner model.ComponentOwner

	err := cr.db.Where("id = ?", componentID).First(&component).Error
	if err != nil {
		return err
	}

	err = cr.db.Where("id = ?", issuerID).First(&user).Error
	if err != nil {
		return err
	}

	if user.Arkhoin < uint64(component.Price) {
		return ErrInsufficientArkhoins
	}

	var holder model.ComponentHolder
	err = cr.db.Where("component_id = ? AND user_id = ?", componentID, issuerID).First(&holder).Error
	if err == nil {
		return ErrComponentAlreadyOwned
	}

	err = cr.db.Where("component_id = ?", componentID).First(&componentOwner).Error
	if err != nil {
		return err
	}

	if componentOwner.UserID != nil {
		err = cr.db.Where("id = ?", *componentOwner.UserID).First(&owner).Error
		if err != nil {
			return err
		}

		if *componentOwner.UserID == issuerID {
			return ErrComponentOwnerCannotBuy
		}

		owner.Arkhoin += uint64(component.Price)
		err = cr.db.Save(&owner).Error
		if err != nil {
			return err
		}
	}

	user.Arkhoin -= uint64(component.Price)
	err = cr.db.Save(&user).Error
	if err != nil {
		return err
	}

	newHolder := model.ComponentHolder{
		ComponentID: componentID,
		UserID:      issuerID,
		PricePaid:   component.Price,
	}

	err = cr.db.Create(&newHolder).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *componentRepository) Sell(issuerID, componentID uint) error {
	var component model.Component
	var user model.User
	var holder model.ComponentHolder
	var componentOwner model.ComponentOwner

	err := cr.db.Where("id = ?", componentID).First(&component).Error
	if err != nil {
		return err
	}

	err = cr.db.Where("component_id = ? AND user_id = ?", componentID, issuerID).First(&holder).Error
	if err != nil {
		return ErrComponentNotOwned
	}

	err = cr.db.Where("id = ?", issuerID).First(&user).Error
	if err != nil {
		return err
	}

	err = cr.db.Where("component_id = ?", componentID).First(&componentOwner).Error
	if err != nil {
		return err
	}

	if componentOwner.UserID != nil && *componentOwner.UserID == issuerID {
		return ErrComponentOwnerCannotSell
	}

	refund := int(holder.PricePaid)
	user.Arkhoin += uint64(refund)

	err = cr.db.Save(&user).Error
	if err != nil {
		return err
	}

	err = cr.db.Delete(&holder).Error
	if err != nil {
		return err
	}

	return nil
}

func (cr *componentRepository) SafeDelete(componentID uint) error {
	tx := cr.db.Begin()
	var component model.Component

	err := tx.Unscoped().Where("id = ?", componentID).First(&component).Error
	if err != nil {
		return err
	}

	if !component.DeletedAt.Valid {
		refund(tx, componentID)

		if err := tx.Delete(&component).Error; err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}

	tx.Rollback()
	return ErrComponentAlreadyTrashed
}

func (cr *componentRepository) Restore(componentID uint) error {
	var component model.Component

	err := cr.db.Unscoped().Where("id = ?", componentID).First(&component).Error
	if err != nil {
		return err
	}

	if component.DeletedAt.Valid {
		return cr.db.Unscoped().Model(&component).Update("deleted_at", nil).Error
	}

	return ErrComponentNotTrashed
}

func (cr *componentRepository) UnsafeDelete(componentID uint) error {
	var component model.Component

	err := cr.db.Unscoped().Where("id = ?", componentID).First(&component).Error
	if err != nil {
		return err
	}

	if !component.DeletedAt.Valid {
		return ErrComponentNotTrashed
	}

	tx := cr.db.Begin()

	if err := tx.Unscoped().Where("component_id = ?", componentID).Delete(&model.ComponentHolder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("component_id = ?", componentID).Delete(&model.ComponentPublication{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("component_id = ?", componentID).Delete(&model.ComponentOwner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&component).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (cr *componentRepository) ClearTrash(userID uint) error {
	tx := cr.db.Begin()

	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM component_owners co
				WHERE co.component_id = component_holders.component_id
				AND co.user_id = ?
			)
		`, userID).
		Delete(&model.ComponentHolder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM component_owners co
				WHERE co.component_id = component_publications.component_id
				AND co.user_id = ?
			)
		`, userID).
		Delete(&model.ComponentPublication{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM component_owners co
				WHERE co.component_id = component_owners.component_id
				AND co.user_id = ?
			)
		`, userID).
		Delete(&model.ComponentOwner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM component_owners co
				WHERE co.component_id = components.id
				AND co.user_id = ?
			)
		`, userID).
		Delete(&model.Component{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
