package translations

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Translation struct {
	Hello                   string `yaml:"hello"`
	InvalidAPIKey           string `yaml:"invalid_api_key"`
	MaximumAPIKey           string `yaml:"maximum_api_key"`
	RequirePermissionAPIKey string `yaml:"require_permission_api_key"`

	CategoryAuth      string `yaml:"category_auth"`
	CategoryFollowers string `yaml:"category_followers"`
	CategoryProject   string `yaml:"category_project"`
	CategoryComponent string `yaml:"category_component"`

	InternalServerError string `yaml:"internal_server_error"`
	Unauthorized        string `yaml:"unauthorized"`
	InvalidBody         string `yaml:"invalid_body"`

	NoAPIKeyFound   string `yaml:"no_api_key_found"` // Used in queries for getting the permissions of keys
	APIKeyDestroyed string `yaml:"api_key_destroyed"`
	APIKeyUpdated   string `yaml:"api_key_updated"`

	AuthDuplicatedUser   string `yaml:"auth_duplicated_user"`
	AuthUserDeleted      string `yaml:"auth_user_deleted"`
	AuthUserUpdated      string `yaml:"auth_user_updated"`
	AuthWrongCredentials string `yaml:"auth_wrong_credentials"`

	NotificationWelcomeUserRegister        string `yaml:"notification_welcome_user_register"`
	NotificationNewLoginDetected           string `yaml:"notification_new_login_detected"`
	NotificationUserFollowedYou            string `yaml:"notification_user_followed_you"`
	NotificationNewProjectCreated          string `yaml:"notification_new_project_created"`
	NotificationUserClonedYourProject      string `yaml:"notification_user_cloned_your_project"`
	NotificationYourProjectPublished       string `yaml:"notification_your_project_published"`
	NotificationYourProjectFavorited       string `yaml:"notification_your_project_favorited"`
	NotificationDeletedProjectFromTrash    string `yaml:"notification_deleted_project_from_trash"`
	NotificationRestoredProjectFromTrash   string `yaml:"notification_restored_project_from_trash"`
	NotificationAddedUserToProject         string `yaml:"notification_added_user_to_project"`
	NotificationRemovedUserFromProject     string `yaml:"notification_removed_user_from_project"`
	NotificationAddedYouToProject          string `yaml:"notification_added_you_to_project"`
	NotificationRemovedYouFromProject      string `yaml:"notification_removed_you_from_project"`
	NotificationUserLeftProject            string `yaml:"notification_user_left_project"`
	NotificationNewComponentCreated        string `yaml:"notification_new_component_created"`
	NotificationYourComponentPublished     string `yaml:"notification_your_component_published"`
	NotificationDeletedComponentFromTrash  string `yaml:"notification_deleted_component_from_trash"`
	NotificationRestoredComponentFromTrash string `yaml:"notification_restored_component_from_trash"`
	NotificationYourComponentBought        string `yaml:"notification_your_component_bought"`
	NotificationYouBoughtComponent         string `yaml:"notification_you_bought_component"`

	NotificationInvalid        string `yaml:"notification_invalid"`
	NotificationAlreadyRead    string `yaml:"notification_already_read"`
	NotificationNotRead        string `yaml:"notification_not_read"`
	NotificationNotAssigned    string `yaml:"notification_not_assigned"`
	NotificationMarkedAsRead   string `yaml:"notification_marked_as_read"`
	NotificationMarkedAsUnread string `yaml:"notification_marked_as_unread"`

	SearchIncorrect string `yaml:"search_incorrect"`
	SearchNoResults string `yaml:"search_no_results"`

	UserDisabledFollowers  string `yaml:"user_disabled_followers"`
	UserDisabledFollowing  string `yaml:"user_disabled_following"`
	UserDisabledProfile    string `yaml:"user_disabled_profile"`
	UserErrorFollowItself  string `yaml:"user_error_follow_itself"`
	UserMissingPermissions string `yaml:"user_missing_permissions"`
	UserNotFound           string `yaml:"user_not_found"`
	UserFollowingAlready   string `yaml:"user_following_already"`
	UserFollowingNot       string `yaml:"user_following_not"`
	UserFollowingStarted   string `yaml:"user_following_started"`
	UserFollowingStopped   string `yaml:"user_following_stopped"`

	ValidatorIncorrectEmailFormat    string `yaml:"validator_incorrect_email_format"`
	ValidatorIncorrectPasswordFormat string `yaml:"validator_incorrect_password_format"`
	ValidatorIncorrectUsernameFormat string `yaml:"validator_incorrect_username_format"`
	ValidatorMaxChars                string `yaml:"validator_max_chars"`
	ValidatorMinChars                string `yaml:"validator_min_chars"`
	ValidatorMustBeSupportedLanguage string `yaml:"validator_must_be_supported_language"`
	ValidatorRequired                string `yaml:"validator_required"`

	ProjectNotFound           string `yaml:"project_not_found"`
	ProjectCreated            string `yaml:"project_created"`
	ProjectUpdated            string `yaml:"project_updated"`
	ProjectDeleted            string `yaml:"project_deleted"`
	ProjectPublished          string `yaml:"project_published"`
	ProjectUnpublished        string `yaml:"project_unpublished"`
	ProjectInvalid            string `yaml:"project_invalid"`
	ProjectMissingPermissions string `yaml:"project_missing_permissions"`
	ProjectTrashed            string `yaml:"project_trashed"`
	ProjectRestored           string `yaml:"project_restored"`
	ProjectAlreadyTrashed     string `yaml:"project_already_trashed"`
	ProjectNotTrashed         string `yaml:"project_not_trashed"`
	ProjectFavorited          string `yaml:"project_favorite"`
	ProjectUnfavorited        string `yaml:"project_unfavorite"`
	ProjectAlreadyFavorited   string `yaml:"project_already_favorite"`
	ProjectNotFavorited       string `yaml:"project_not_favorite"`
	ProjectForked             string `yaml:"project_forked"`
	ProjectIsNotAFork         string `yaml:"project_is_not_a_fork"`
	ProjectUnlinked           string `yaml:"project_unlinked"`
	ProjectAssignedUser       string `yaml:"project_assigned_user"`
	ProjectUnassignedUser     string `yaml:"project_unassigned_user"`
	ProjectEmptyAssign        string `yaml:"project_empty_assign"`
	ProjectUserNotAssigned    string `yaml:"project_user_not_assigned"`
	ProjectCannotAssignOwner  string `yaml:"project_cannot_assign_owner"`

	UpstreamNotPublic    string `yaml:"upstream_not_public"`
	TrashCleared         string `yaml:"trash_cleared"`
	InsufficientArkhoins string `yaml:"insufficient_arkhoins"`

	ComponentCreated         string `yaml:"component_created"`
	ComponentUpdated         string `yaml:"component_updated"`
	ComponentPublished       string `yaml:"component_published"`
	ComponentUnpublished     string `yaml:"component_unpublished"`
	ComponentBought          string `yaml:"component_bought"`
	ComponentSold            string `yaml:"component_sold"`
	ComponentTrashed         string `yaml:"component_trashed"`
	ComponentRestored        string `yaml:"component_restored"`
	ComponentDeleted         string `yaml:"component_deleted"`
	ComponentInvalid         string `yaml:"component_invalid"`
	ComponentNotFound        string `yaml:"component_not_found"`
	ComponentAlreadyTrashed  string `yaml:"component_already_trashed"`
	ComponentNotTrashed      string `yaml:"component_not_trashed"`
	ComponentAlreadyOwned    string `yaml:"component_already_owned"`
	ComponentNotOwned        string `yaml:"component_not_owned"`
	ComponentAlreadyPublic   string `yaml:"component_already_public"`
	ComponentNotPublic       string `yaml:"yaml:component_not_public"`
	ComponentOwnerCannotBuy  string `yaml:"component_owner_cannot_buy"`
	ComponentOwnerCannotSell string `yaml:"component_owner_cannot_sell"`

	PasswordResetRequest       string `yaml:"password_reset_request"`
	PasswordResetSuccess       string `yaml:"password_reset_success"`
	PasswordResetEmailSubject  string `yaml:"password_reset_email_subject"`
	PasswordResetEmailTemplate string `yaml:"password_reset_email_template"`
	InvalidPasswordResetKey    string `yaml:"invalid_password_reset_key"`

	UnsupportedFileType string `yaml:"unsupported_file_type"`
	UnableToDecodeFile  string `yaml:"unable_to_decode_file"`
	UnableToEncodeFile  string `yaml:"unable_to_encode_file"`
	FileTooLarge        string `yaml:"file_too_large"`
}

var Translations = make(map[string]Translation)

func readYAMLFile(filename string) (*Translation, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var translation Translation
	err = yaml.Unmarshal(data, &translation)
	if err != nil {
		return nil, err
	}

	return &translation, nil
}

func Init(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("error reading directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		lang := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		translation, err := readYAMLFile(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatalf("error reading %s: %v", file.Name(), err)
		}

		Translations[lang] = *translation
	}
}

func GetTranslation(ctx *gin.Context) Translation {
	return ctx.Keys["lang"].(Translation)
}
