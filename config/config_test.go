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
