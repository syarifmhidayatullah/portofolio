package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	SessionSecret string

	// Email
	EmailDriver  string // "smtp" or "resend"
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	ResendAPIKey string
	NotifyEmail  string

	// Admin seed
	AdminEmail    string
	AdminPassword string

	// App info
	AppName string
	AppURL  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Port:          getEnv("PORT", ":8080"),
		SessionSecret: getEnv("SESSION_SECRET", "super-secret-change-this"),

		EmailDriver:  getEnv("EMAIL_DRIVER", "smtp"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		NotifyEmail:  getEnv("NOTIFY_EMAIL", ""),

		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),

		AppName: getEnv("APP_NAME", "Syarif Hidayatullah"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),
	}

	cfg.DatabaseURL = buildDSN()
	return cfg
}

func buildDSN() string {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	name := getEnv("DB_NAME", "portfolio")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
