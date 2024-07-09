package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	Router struct {
		GinMode string `yaml:"gin_mode"`
		Address string `yaml:"addr"`
		Port    uint16 `yaml:"port"`
	}

	Postgres struct {
		Host     string `yaml:"host"`
		DB       string `yaml:"db"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	}

	Security struct {
		BcryptCost int    `yaml:"bcrypt_cost"`
		JWTSecret  string `yaml:"jwt_secret"`
	}

	Permissions struct {
		Admin             string `yaml:"admin"`
		ManageAPIKey      string `yaml:"manage_api_key"`
		ManageUser        string `yaml:"manage_user"`
		ManagePermissions string `yaml:"manage_permissions"`
		ManageProjects    string `yaml:"manage_projects"`
		ManageStore       string `yaml:"manage_store"`
	}
)

func Parse() {
	if missingVars := checkEnvVars(
		"POSTGRES_HOST",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_SSLMODE",
		"JWT_SECRET",
	); len(missingVars) > 0 {
		log.Println("You can override the following env variables to get rid of this error:")
		log.Println(strings.Join(missingVars, ", "))
		log.Println("Loading .env file...")
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}

	if err := yaml.Unmarshal(read("router.yaml"), &Router); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("postgres.yaml"), &Postgres); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("security.yaml"), &Security); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("permissions.yaml"), &Permissions); err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Print("Loaded config files")
}

func read(file string) []byte {
	data, err := os.ReadFile(filepath.Join("config", file))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return []byte(os.ExpandEnv(string(data)))
}

func checkEnvVars(vars ...string) []string {
	missingVars := []string{}
	for _, v := range vars {
		if _, exists := os.LookupEnv(v); !exists {
			missingVars = append(missingVars, v)
		}
	}
	if len(missingVars) > 0 {
		log.Printf("Missing environment variables: %s", strings.Join(missingVars, ", "))
		return missingVars
	}
	return nil
}
