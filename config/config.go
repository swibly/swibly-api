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
		GinMode     string `yaml:"gin_mode"`
		Address     string `yaml:"addr"`
		Port        uint16 `yaml:"port"`
		Environment string `yaml:"environment"`
	}

	Postgres struct {
		ConnectionString string `yaml:"connection_string"`
	}

	S3 struct {
		Access string `yaml:"access"`
		Secret string `yaml:"secret"`
		URL    string `yaml:"url"`
		SURL   string `yaml:"surl"`
		Region string `yaml:"region"`
		Bucket string `yaml:"bucket"`
	}

	Security struct {
		BcryptCost int    `yaml:"bcrypt_cost"`
		JWTSecret  string `yaml:"jwt_secret"`
	}

	SMTP struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Email    string `yaml:"email"`
		Password string `yaml:"password"`
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
		"POSTGRES_CONNECTION_STRING",
		"JWT_SECRET",
		"SMTP_HOST",
		"SMTP_PORT",
		"SMTP_USERNAME",
		"SMTP_EMAIL",
		"SMTP_PASSWORD",
		"S3_ACCESS_KEY",
		"S3_SECRET_KEY",
		"ENVIRONMENT",
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

	if err := yaml.Unmarshal(read("s3.yaml"), &S3); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("security.yaml"), &Security); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("smtp.yaml"), &SMTP); err != nil {
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
