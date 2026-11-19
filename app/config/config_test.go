package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_UsesDefaultsWhenEnvIsUnset(t *testing.T) {
	for _, key := range []string{
		"DEBUG_MODE",
		"MONGODB_URL",
		"MONGODB_REQUEST_TIMEOUT",
		"JWT_SECRET",
		"JWT_EXPIRE",
	} {
		require.NoError(t, os.Unsetenv(key))
	}

	cfg, err := config.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.False(t, cfg.DebugMode)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoDBUrl)
	assert.Equal(t, 10*time.Second, cfg.MongoDBRequestTimeout)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, 3600*time.Second, cfg.JWTExpire)
}

func TestNewConfig_UsesEnvValuesWhenValid(t *testing.T) {
	t.Setenv("DEBUG_MODE", "true")
	t.Setenv("MONGODB_URL", "mongodb://mongo:27017")
	t.Setenv("MONGODB_REQUEST_TIMEOUT", "45")
	t.Setenv("JWT_SECRET", "my-secret")
	t.Setenv("JWT_EXPIRE", "7200")

	cfg, err := config.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.True(t, cfg.DebugMode)
	assert.Equal(t, "mongodb://mongo:27017", cfg.MongoDBUrl)
	assert.Equal(t, 45*time.Second, cfg.MongoDBRequestTimeout)
	assert.Equal(t, "my-secret", cfg.JWTSecret)
	assert.Equal(t, 7200*time.Second, cfg.JWTExpire)
}

func TestNewConfig_FallsBackForInvalidBoolAndDuration(t *testing.T) {
	t.Setenv("DEBUG_MODE", "not-a-bool")
	t.Setenv("MONGODB_REQUEST_TIMEOUT", "not-a-number")
	t.Setenv("JWT_EXPIRE", "invalid")

	cfg, err := config.NewConfig()

	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.False(t, cfg.DebugMode)
	assert.Equal(t, 10*time.Second, cfg.MongoDBRequestTimeout)
	assert.Equal(t, 3600*time.Second, cfg.JWTExpire)
}
