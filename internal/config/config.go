package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	AuthToken    string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SenderName   string
}

func Load() (*Config, error) {
	godotenv.Load()

	authToken := os.Getenv("AUTH_TOKEN")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	for _, key := range []string{"AUTH_TOKEN", "SMTP_USER", "SMTP_PASSWORD"} {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("missing required env: %s", key)
		}
	}

	return &Config{
		ServerPort:   envOr("SERVER_PORT", "8080"),
		AuthToken:    authToken,
		SMTPHost:     envOr("SMTP_HOST", "smtp.163.com"),
		SMTPPort:     envInt("SMTP_PORT", 465),
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		SenderName:   envOr("SENDER_NAME", "Msg Forwarder"),
	}, nil
}

func envOr(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}