package config_test

import (
	"avito_shop/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_AllFields(t *testing.T) {
	t.Setenv("DATABASE_HOST", "localhost")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("DATABASE_USER", "user")
	t.Setenv("DATABASE_NAME", "dbname")
	t.Setenv("DATABASE_PASSWORD", "password")
	t.Setenv("JWT_SECRET_KEY", "secret")
	t.Setenv("SERVER_PORT", "8080")

	cfg := config.Load()

	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "user", cfg.DBUser)
	assert.Equal(t, "dbname", cfg.DBName)
	assert.Equal(t, "password", cfg.DBPassword)
	assert.Equal(t, "secret", cfg.JWTSecretKey)
	assert.Equal(t, "8080", cfg.ServerPort)
}
