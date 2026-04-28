package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres", cfg.DBHost)
	assert.Equal(t, "redis:6379", cfg.RedisAddr)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, 0, cfg.RedisDB)
}

func TestDBPath(t *testing.T) {
	cfg := Load()
	assert.Contains(t, cfg.DBPath(), "postgres://")
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", result)

	result2 := getEnv("NON_EXISTENT", "default")
	assert.Equal(t, "default", result2)
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT_KEY", "2")
	defer os.Unsetenv("TEST_INT_KEY")

	result := getEnvInt("TEST_INT_KEY", 0)
	assert.Equal(t, 2, result)

	result2 := getEnvInt("NON_EXISTENT_INT", 7)
	assert.Equal(t, 7, result2)
}
