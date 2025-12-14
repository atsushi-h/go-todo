package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConfig_DSN(t *testing.T) {
	t.Run("通常のパスワード", func(t *testing.T) {
		cfg := DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "user",
			Password: "pass",
		}
		expected := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"
		assert.Equal(t, expected, cfg.DSN())
	})

	t.Run("特殊文字を含むパスワード", func(t *testing.T) {
		cfg := DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "user@domain",
			Password: "p@ss:w/ord",
		}
		// 特殊文字がURLエンコードされていることを確認
		dsn := cfg.DSN()
		assert.Contains(t, dsn, "user%40domain")        // @ -> %40
		assert.Contains(t, dsn, "p%40ss%3Aw%2Ford")     // @:/ -> %40%3A%2F
		assert.NotContains(t, dsn, "user@domain")
		assert.NotContains(t, dsn, "p@ss:w/ord")
	})
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
		errContains string
	}{
		{name: "valid HTTP URL", callbackURL: "http://localhost:4000/callback", wantErr: false},
		{name: "valid HTTPS URL", callbackURL: "https://example.com/callback", wantErr: false},
		{name: "invalid URL format", callbackURL: "://invalid", wantErr: true, errContains: "invalid callback URL"},
		{name: "invalid scheme - ftp", callbackURL: "ftp://example.com/callback", wantErr: true, errContains: "must use http or https scheme"},
		{name: "no scheme", callbackURL: "example.com/callback", wantErr: true, errContains: "must use http or https scheme"},
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
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFrontendConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantErr     bool
		errContains string
	}{
		{name: "valid HTTP URL", url: "http://localhost:3000", wantErr: false},
		{name: "valid HTTPS URL", url: "https://example.com", wantErr: false},
		{name: "invalid URL format", url: "://invalid", wantErr: true, errContains: "invalid frontend URL"},
		{name: "invalid scheme - ftp", url: "ftp://example.com", wantErr: true, errContains: "must use http or https scheme"},
		{name: "no scheme", url: "localhost:3000", wantErr: true, errContains: "must use http or https scheme"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := FrontendConfig{
				URL: tt.url,
			}

			err := cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoad_Success(t *testing.T) {
	// 環境変数を設定（t.Setenvを使用して自動クリーンアップ）
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DB", "testdb")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "pass")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("GOOGLE_CLIENT_ID", "client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "secret")
	t.Setenv("OAUTH_CALLBACK_URL", "http://localhost:4000/callback")

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
	// 無効なポート番号を設定（t.Setenvを使用して自動クリーンアップ）
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "99999")
	t.Setenv("POSTGRES_DB", "testdb")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "pass")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("GOOGLE_CLIENT_ID", "client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "secret")
	t.Setenv("OAUTH_CALLBACK_URL", "http://localhost:4000/callback")

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
