package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultBaseURL               = "http://localhost:8080"
	defaultServiceName           = "backend-challenge"
	defaultMongoDBUrl            = "mongodb://localhost:27017"
	defaultMongoDBName           = "backend-challenge"
	defaultMongoDBRequestTimeout = 10 * time.Second
	defaultJWTSecret             = "secret"
	defaultJWTExpire             = 3600 * time.Second
)

type Config struct {
	ServiceName           string
	BaseURL               string
	DebugMode             bool
	MongoDBUrl            string
	MongoDBName           string
	MongoDBUser           string
	MongoDBPassword       string
	MongoDBRequestTimeout time.Duration
	JWTSecret             string
	JWTExpire             time.Duration
}

func NewConfig() (*Config, error) {

	return &Config{
		ServiceName:           getEnvString("SERVICE_NAME", defaultServiceName),
		BaseURL:               getEnvString("BASE_URL", defaultBaseURL),
		DebugMode:             getEnvBool("DEBUG_MODE", false),
		MongoDBUrl:            getEnvString("MONGODB_URL", defaultMongoDBUrl),
		MongoDBName:           getEnvString("MONGODB_NAME", defaultMongoDBName),
		MongoDBUser:           getEnvString("MONGODB_USER", ""),
		MongoDBPassword:       getEnvString("MONGODB_PASSWORD", ""),
		MongoDBRequestTimeout: getEnvSeconds("MONGODB_REQUEST_TIMEOUT", defaultMongoDBRequestTimeout),
		JWTSecret:             getEnvString("JWT_SECRET", defaultJWTSecret),
		JWTExpire:             getEnvSeconds("JWT_EXPIRE", defaultJWTExpire),
	}, nil
}

// getEnvString returns the value of the environment variable key, or
// fallback when it is unset or empty.
func getEnvString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvBool parses the environment variable key as a bool, returning
// fallback when it is unset or cannot be parsed.
func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// getEnvSeconds parses the environment variable key as a whole number of
// seconds, returning fallback when it is unset or cannot be parsed.
func getEnvSeconds(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
