package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	_, currFile, _, _ := runtime.Caller(0)
	currDir := filepath.Dir(currFile)
	envPath := filepath.Join(currDir, "..", "..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}
