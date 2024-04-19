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

	Postgre struct {
		Host     string `yaml:"host"`
		DB       string `yaml:"db"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}

	JWT struct {
		Secret string `yaml:"secret"`
	}
)

func Parse() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("router.yaml"), &Router); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("postgre.yaml"), &Postgre); err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(read("jwt.yaml"), &JWT); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func read(file string) []byte {
	data, err := os.ReadFile(filepath.Join("config", file))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return []byte(os.ExpandEnv(string(data)))
}
