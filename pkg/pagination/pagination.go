package pagination

import (
	"fmt"
	"math"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"gorm.io/gorm"
)

func Generate[T any](query *gorm.DB, page, pageSize int) (*dto.Pagination[T], error) {
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to count total records: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	data := make([]*T, 0, pageSize)

	if err := query.Limit(limit).Offset(offset).Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	pagination := &dto.Pagination[T]{
		Data:         data,
		TotalRecords: int(totalRecords),
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     page + 1,
		PreviousPage: page - 1,
	}

	if pagination.NextPage > totalPages {
		pagination.NextPage = -1
	}

	if pagination.PreviousPage < 1 {
		pagination.PreviousPage = -1
	}

	return pagination, nil
}
