package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	EntraClientID     string
	EntraTenantID     string
	EntraClientSecret string
	DataDir           string
	Port              string
	MaxUploadSize     int64
}

var AppConfig Config

func LoadConfig() {
	maxUploadMBStr := getEnv("MAX_UPLOAD_SIZE_MB", "50")
	maxUploadMB, err := strconv.ParseInt(maxUploadMBStr, 10, 64)
	if err != nil || maxUploadMB <= 0 {
		maxUploadMB = 50
	}

	AppConfig = Config{
		EntraClientID:     getEnv("ENTRA_CLIENT_ID", ""),
		EntraTenantID:     getEnv("ENTRA_TENANT_ID", ""),
		EntraClientSecret: getEnv("ENTRA_CLIENT_SECRET", ""),
		DataDir:           getEnv("DATA_DIR", "./data"),
		Port:              getEnv("PORT", "8080"),
		MaxUploadSize:     maxUploadMB << 20,
	}

	if AppConfig.EntraClientID == "" || AppConfig.EntraTenantID == "" || AppConfig.EntraClientSecret == "" {
		log.Println("[WARNING] Missing Entra ID environment variables")
	} else {
		log.Println("[INFO] Config loaded")
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
