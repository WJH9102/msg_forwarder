package config

import (
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

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort:   envOr("SERVER_PORT", "8080"),
		AuthToken:    mustEnv("AUTH_TOKEN"),
		SMTPHost:     envOr("SMTP_HOST", "smtp.163.com"),
		SMTPPort:     envInt("SMTP_PORT", 465),
		SMTPUser:     mustEnv("SMTP_USER"),
		SMTPPassword: mustEnv("SMTP_PASSWORD"),
		SenderName:   envOr("SENDER_NAME", "Msg Forwarder"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing required env: " + key)
	}
	return v
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