package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "user",
		Password: "pass",
	}
	expected := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
	assert.Equal(t, expected, cfg.DSN())
}

func TestDatabaseConfig_String(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "user",
		Password: "secret123",
	}

	str := cfg.String()

	// パスワードがマスクされていることを確認
	assert.Contains(t, str, "Password:***")
	assert.NotContains(t, str, "secret123")
	assert.Contains(t, str, "localhost")
	assert.Contains(t, str, "testdb")
}

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{name: "valid port", port: 5432, wantErr: false},
		{name: "min valid port", port: 1, wantErr: false},
		{name: "max valid port", port: 65535, wantErr: false},
		{name: "port too low", port: 0, wantErr: true},
		{name: "port too high", port: 65536, wantErr: true},
		{name: "negative port", port: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DatabaseConfig{
				Host:     "localhost",
				Port:     tt.port,
				Database: "testdb",
				User:     "user",
				Password: "pass",
			}

			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRedisConfig_Address(t *testing.T) {
	cfg := RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	expected := "localhost:6379"
	assert.Equal(t, expected, cfg.Address())
}

func TestRedisConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{name: "valid port", port: 6379, wantErr: false},
		{name: "invalid port", port: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RedisConfig{
				Host: "localhost",
				Port: tt.port,
			}

			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServerConfig_Address(t *testing.T) {
	cfg := ServerConfig{
		Port: 4000,
	}
	expected := ":4000"
	assert.Equal(t, expected, cfg.Address())
}

func TestOAuthConfig_String(t *testing.T) {
	cfg := OAuthConfig{
		GoogleClientID:     "client123",
		GoogleClientSecret: "secret456",
		CallbackURL:        "http://localhost:4000/auth/callback",
	}

	str := cfg.String()

	// シークレットがマスクされていることを確認
	assert.Contains(t, str, "GoogleClientSecret:***")
	assert.NotContains(t, str, "secret456")
	assert.Contains(t, str, "client123")
}

func TestOAuthConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		callbackURL string
		wantErr     bool
	}{
		{name: "valid URL", callbackURL: "http://localhost:4000/callback", wantErr: false},
		{name: "valid HTTPS URL", callbackURL: "https://example.com/callback", wantErr: false},
		{name: "invalid URL", callbackURL: "://invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := OAuthConfig{
				GoogleClientID:     "client",
				GoogleClientSecret: "secret",
				CallbackURL:        tt.callbackURL,
			}

			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFrontendConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{name: "valid URL", url: "http://localhost:3000", wantErr: false},
		{name: "valid HTTPS URL", url: "https://example.com", wantErr: false},
		{name: "invalid URL", url: "://invalid", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := FrontendConfig{
				URL: tt.url,
			}

			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoad_Success(t *testing.T) {
	// 環境変数を設定
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DB", "testdb")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "pass")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("GOOGLE_CLIENT_ID", "client")
	os.Setenv("GOOGLE_CLIENT_SECRET", "secret")
	os.Setenv("OAUTH_CALLBACK_URL", "http://localhost:4000/callback")

	defer func() {
		// クリーンアップ
		os.Unsetenv("POSTGRES_HOST")
		os.Unsetenv("POSTGRES_PORT")
		os.Unsetenv("POSTGRES_DB")
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("REDIS_PORT")
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")
		os.Unsetenv("OAUTH_CALLBACK_URL")
	}()

	cfg, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "testdb", cfg.Database.Database)
}

func TestLoad_MissingRequired(t *testing.T) {
	// 必須の環境変数を設定しない
	os.Clearenv()

	_, err := Load()

	assert.Error(t, err)
}

func TestLoad_InvalidPort(t *testing.T) {
	// 無効なポート番号を設定
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "99999")
	os.Setenv("POSTGRES_DB", "testdb")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "pass")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("GOOGLE_CLIENT_ID", "client")
	os.Setenv("GOOGLE_CLIENT_SECRET", "secret")
	os.Setenv("OAUTH_CALLBACK_URL", "http://localhost:4000/callback")

	defer os.Clearenv()

	_, err := Load()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port")
}

func TestConfig_Validate(t *testing.T) {
	t.Run("すべて有効", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				User:     "user",
				Password: "pass",
			},
			Redis: RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
			OAuth: OAuthConfig{
				GoogleClientID:     "client",
				GoogleClientSecret: "secret",
				CallbackURL:        "http://localhost:4000/callback",
			},
			Server: ServerConfig{
				Port: 4000,
			},
			Frontend: FrontendConfig{
				URL: "http://localhost:3000",
			},
		}

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("無効なDatabaseポート", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Port: 99999,
			},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database config")
	})
}
