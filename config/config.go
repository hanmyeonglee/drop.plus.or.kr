package config

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type Config struct {
	EntraClientID     string
	EntraTenantID     string
	EntraClientSecret string
	Port              string
	BaseURL           string
	DataDir           string
	MaxUploadSize     int64
	AutoDeleteSeconds int
}

var (
	AppConfig   Config
	OAuthConfig *oauth2.Config
	Store       *sessions.CookieStore
)

func LoadConfig() {
	maxUploadMBStr := getEnv("MAX_UPLOAD_SIZE_MB", "50")
	maxUploadMB, err := strconv.ParseInt(maxUploadMBStr, 10, 64)
	if err != nil || maxUploadMB <= 0 {
		maxUploadMB = 50
	}

	autoDeleteStr := getEnv("AUTO_DELETE_SECONDS", "2592000")
	autoDelete, err := strconv.Atoi(autoDeleteStr)
	if err != nil {
		autoDelete = 2592000
	}

	AppConfig = Config{
		EntraClientID:     getEnv("ENTRA_CLIENT_ID", ""),
		EntraTenantID:     getEnv("ENTRA_TENANT_ID", ""),
		EntraClientSecret: getEnv("ENTRA_CLIENT_SECRET", ""),
		Port:              getEnv("PORT", "8080"),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
		DataDir:           getEnv("DATA_DIR", "./data"),
		MaxUploadSize:     maxUploadMB << 20,
		AutoDeleteSeconds: autoDelete,
	}

	if AppConfig.EntraClientID == "" || AppConfig.EntraTenantID == "" || AppConfig.EntraClientSecret == "" {
		log.Println("[WARNING] Missing Entra ID environment variables")
	} else {
		log.Println("[INFO] Config loaded")
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     AppConfig.EntraClientID,
		ClientSecret: AppConfig.EntraClientSecret,
		RedirectURL:  AppConfig.BaseURL + "/auth/callback",
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     microsoft.AzureADEndpoint(AppConfig.EntraTenantID),
	}

	sessionKey := getEnv("SESSION_KEY", "")
	if sessionKey == "" {
		key := make([]byte, 32)
		rand.Read(key)
		sessionKey = string(key)
	}
	Store = sessions.NewCookieStore([]byte(sessionKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   strings.HasPrefix(AppConfig.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
