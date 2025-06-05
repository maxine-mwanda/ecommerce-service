// internal/config/config.go
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Logging       LoggingConfig
	Auth          AuthConfig
	Notifications NotificationConfig
}

type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"8080"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"3306"`
	User     string `envconfig:"DB_USER" default:"root"`
	Password string `envconfig:"DB_PASSWORD" default:"root"`
	Name     string `envconfig:"DB_NAME" default:"ecommerce"`
}

type LoggingConfig struct {
	Directory string `envconfig:"LOG_DIR" default:"./logs"`
	Filename  string `envconfig:"LOG_FILE" default:"ecommerce.log"`
}

type AuthConfig struct {
	IssuerURL    string `envconfig:"OIDC_ISSUER_URL" required:"true"`
	ClientID     string `envconfig:"OIDC_CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"OIDC_CLIENT_SECRET" required:"true"`
	RedirectURL  string `envconfig:"OIDC_REDIRECT_URL" required:"true"`
}

type NotificationConfig struct {
	SMS   SMSConfig
	Email EmailConfig
}

type SMSConfig struct {
	APIKey   string `envconfig:"SMS_API_KEY"`
	Username string `envconfig:"SMS_USERNAME"`
	SenderID string `envconfig:"SMS_SENDER_ID" default:"ECOMM"`
	Sandbox  bool   `envconfig:"SMS_SANDBOX" default:"true"`
}

type EmailConfig struct {
	SMTPHost         string `envconfig:"SMTP_HOST"`
	SMTPPort         int    `envconfig:"SMTP_PORT" default:"587"`
	SMTPUser         string `envconfig:"SMTP_USER"`
	SMTPPassword     string `envconfig:"SMTP_PASSWORD"`
	FromEmail        string `envconfig:"FROM_EMAIL" default:"no-reply@ecommerce.com"`
	SMTPEmail        string `envconfig:"SMTP_EMAIL" default:"`
	SMTPClientID     string `envconfig:"SMTP_CLIENT_ID"`
	SMTPClientSecret string `envconfig:"SMTP_CLIENT_SECRET"`
	SMTPRefreshToken string `envconfig:"SMTP_REFRESH_TOKEN"`
	ClientID         string // from Google Cloud Console OAuth2
	ClientSecret     string
	RefreshToken     string // Google OAuth2 refresh token with mail.google.com scope
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Printf("envconfig.Process failed: %v", err)
		return nil, err
	}
	return &cfg, nil
}
