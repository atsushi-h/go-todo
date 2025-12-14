package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	OAuth    OAuthConfig
	Server   ServerConfig
	Frontend FrontendConfig
	Cookie   CookieConfig
}

// DatabaseConfig holds PostgreSQL configuration
type DatabaseConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" required:"true"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	Database string `envconfig:"POSTGRES_DB" required:"true"`
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
}

// DSN returns PostgreSQL connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
	)
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host string `envconfig:"REDIS_HOST" required:"true"`
	Port int    `envconfig:"REDIS_PORT" default:"6379"`
}

// Address returns Redis connection address
func (r *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
	CallbackURL        string `envconfig:"OAUTH_CALLBACK_URL" required:"true"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int `envconfig:"BACKEND_CONTAINER_PORT" default:"4000"`
}

// Address returns server listen address
func (s *ServerConfig) Address() string {
	return fmt.Sprintf(":%d", s.Port)
}

// FrontendConfig holds frontend-related configuration
type FrontendConfig struct {
	URL string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
}

// CookieConfig holds cookie security configuration
type CookieConfig struct {
	Secure bool `envconfig:"COOKIE_SECURE" default:"false"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}
