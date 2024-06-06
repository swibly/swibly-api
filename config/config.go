package config

import (
	"log"
	"os"
	"path/filepath"

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
		ManageUser        string `yaml:"manage_user"`
		ManagePermissions string `yaml:"manager_permissions"`
		ManageProjects    string `yaml:"manage_projects"`
		ManageStore       string `yaml:"manage_store"`
	}
)

func Parse() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error: %v", err)
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
