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
	Clone                   string `yaml:"clone"`

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

	SearchIncorrect string `yaml:"search_incorrect"`
	SearchNoResults string `yaml:"search_no_results"`

	UserDisabledFollowers string `yaml:"user_disabled_followers"`
	UserDisabledFollowing string `yaml:"user_disabled_following"`
	UserDisabledProfile   string `yaml:"user_disabled_profile"`
	UserErrorFollowItself string `yaml:"user_error_follow_itself"`
	UserNotFound          string `yaml:"user_not_found"`
	UserFollowingAlready  string `yaml:"user_following_already"`
	UserFollowingNot      string `yaml:"user_following_not"`
	UserFollowingStarted  string `yaml:"user_following_started"`
	UserFollowingStopped  string `yaml:"user_following_stopped"`

	ValidatorIncorrectEmailFormat    string `yaml:"validator_incorrect_email_format"`
	ValidatorIncorrectPasswordFormat string `yaml:"validator_incorrect_password_format"`
	ValidatorIncorrectUsernameFormat string `yaml:"validator_incorrect_username_format"`
	ValidatorMaxChars                string `yaml:"validator_max_chars"`
	ValidatorMinChars                string `yaml:"validator_min_chars"`
	ValidatorMustBeNumericalBoolean  string `yaml:"validator_must_be_numerical_boolean"`
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

	TrashCleared string `yaml:"trash_cleared"`
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
