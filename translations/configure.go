package translations

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Translation struct {
	Hello                   string `yaml:"hello"`
	InvalidAPIKey           string `yaml:"invalid_api_key"`
	MaximumAPIKey           string `yaml:"maximum_api_key"`
	RequirePermissionAPIKey string `yaml:"require_permission_api_key"`

	InternalServerError string `yaml:"internal_server_error"`
	InvalidBody         string `yaml:"invalid_body"`

	NoAPIKeyFound   string `yaml:"no_api_key_found"` // Used in queries for getting the permissions of keys
	APIKeyDestroyed string `yaml:"api_key_destroyed"`
	APIKeyUpdated   string `yaml:"api_key_updated"`
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
