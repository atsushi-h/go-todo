package config

import (
	"fmt"
	"net/url"

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

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}
	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config: %w", err)
	}
	if err := c.OAuth.Validate(); err != nil {
		return fmt.Errorf("oauth config: %w", err)
	}
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}
	if err := c.Frontend.Validate(); err != nil {
		return fmt.Errorf("frontend config: %w", err)
	}
	return nil
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

// String returns a safe string representation with masked password
func (d *DatabaseConfig) String() string {
	return fmt.Sprintf(
		"DatabaseConfig{Host:%s, Port:%d, Database:%s, User:%s, Password:***}",
		d.Host, d.Port, d.Database, d.User,
	)
}

// Validate checks if the database configuration is valid
func (d *DatabaseConfig) Validate() error {
	if d.Port < 1 || d.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", d.Port)
	}
	return nil
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

// Validate checks if the Redis configuration is valid
func (r *RedisConfig) Validate() error {
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", r.Port)
	}
	return nil
}

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	GoogleClientID     string `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	GoogleClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
	CallbackURL        string `envconfig:"OAUTH_CALLBACK_URL" required:"true"`
}

// String returns a safe string representation with masked secret
func (o *OAuthConfig) String() string {
	return fmt.Sprintf(
		"OAuthConfig{GoogleClientID:%s, GoogleClientSecret:***, CallbackURL:%s}",
		o.GoogleClientID, o.CallbackURL,
	)
}

// Validate checks if the OAuth configuration is valid
func (o *OAuthConfig) Validate() error {
	if _, err := url.Parse(o.CallbackURL); err != nil {
		return fmt.Errorf("invalid callback URL: %w", err)
	}
	return nil
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int `envconfig:"BACKEND_CONTAINER_PORT" default:"4000"`
}

// Address returns server listen address
func (s *ServerConfig) Address() string {
	return fmt.Sprintf(":%d", s.Port)
}

// Validate checks if the server configuration is valid
func (s *ServerConfig) Validate() error {
	if s.Port < 1 || s.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", s.Port)
	}
	return nil
}

// FrontendConfig holds frontend-related configuration
type FrontendConfig struct {
	URL string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
}

// Validate checks if the frontend configuration is valid
func (f *FrontendConfig) Validate() error {
	if _, err := url.Parse(f.URL); err != nil {
		return fmt.Errorf("invalid frontend URL: %w", err)
	}
	return nil
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

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}
