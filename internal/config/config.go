package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

const (
	devJWTSecretKey   = "this-is-obviously-very-insecure#$"
	minJWTSecretBytes = 32

	defaultJWTExpiresIn           = 15 * time.Minute
	defaultJWTRefreshTokenExpires = 7 * 24 * time.Hour
)

type Config struct {
	Env      Environment
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
}

type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

func Load() (Config, error) {
	var errs []error

	env := Environment(getEnv("APP_ENV", string(Development)))
	switch env {
	case Development, Staging, Production:
	default:
		return Config{}, fmt.Errorf("APP_ENV must be one of development, staging, production; got %q", env)
	}

	jwtSecret, err := resolveJWTSecret(env)
	if err != nil {
		errs = append(errs, err)
	}

	cfg := Config{
		Env: env,
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", ":3030"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:              jwtSecret,
			ExpiresIn:           getDuration("JWT_EXPIRES_IN", defaultJWTExpiresIn, &errs),
			RefreshTokenExpires: getDuration("JWT_REFRESH_EXPIRES_IN", defaultJWTRefreshTokenExpires, &errs),
		},
	}

	if len(errs) > 0 {
		return Config{}, fmt.Errorf("config: %w", errors.Join(errs...))
	}
	return cfg, nil
}

func (c Config) IsDev() bool { return c.Env == Development }

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration, errs *[]error) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: %w", key, err))
		return fallback
	}
	return d
}

func resolveJWTSecret(env Environment) (string, error) {
	const key = "JWT_SECRET"
	secret := getEnv(key, devJWTSecretKey)
	if env == Development {
		return secret, nil
	}
	if secret == devJWTSecretKey {
		return "", fmt.Errorf("%s must be set for %s environment", key, env)
	}
	if len(secret) < minJWTSecretBytes {
		return "", fmt.Errorf("%s does not meet the minimum length of %d", key, minJWTSecretBytes)
	}
	return secret, nil
}
