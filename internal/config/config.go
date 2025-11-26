package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	Logging  LoggingConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type SecurityConfig struct {
	APIKeyEnabled  bool
	APIKey         string
	RateLimitRPS   int
	MaxRequestSize int64
	AllowedOrigins []string
	EnableCORS     bool
}

type LoggingConfig struct {
	Level  string
	Format string
}

type DatabaseConfig struct {
	Path string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", ":3000"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 15*time.Second),
		},
		Security: SecurityConfig{
			APIKeyEnabled:  getBoolEnv("SECURITY_API_KEY_ENABLED", false),
			APIKey:         getEnv("SECURITY_API_KEY", ""),
			RateLimitRPS:   getIntEnv("SECURITY_RATE_LIMIT_RPS", 100),
			MaxRequestSize: getInt64Env("SECURITY_MAX_REQUEST_SIZE", 1<<20), // 1MB
			AllowedOrigins: []string{getEnv("SECURITY_ALLOWED_ORIGINS", "*")},
			EnableCORS:     getBoolEnv("SECURITY_ENABLE_CORS", true),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "data/mcp.db"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Security.APIKeyEnabled && c.Security.APIKey == "" {
		return errors.New("API key required when API key auth is enabled")
	}
	if c.Security.RateLimitRPS < 1 {
		return errors.New("rate limit must be at least 1 RPS")
	}
	if c.Server.Port == "" {
		return errors.New("server port is required")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getInt64Env(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultVal
}

func getBoolEnv(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
