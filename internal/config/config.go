package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
	IsDev  bool
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Name, c.Port)
}

type JWTConfig struct {
	Secret string
}

type ServerConfig struct {
	GRPCPort string
	HTTPPort string
}

func Load() (*Config, error) {
	cfg := &Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_KEY"),
		},
		Server: ServerConfig{
			GRPCPort: getEnvOrDefault("GRPC_PORT", "50051"),
			HTTPPort: getEnvOrDefault("HTTP_PORT", "8080"),
		},
		IsDev: os.Getenv("ENV") != "production",
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DB.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DB.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.DB.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_KEY is required")
	}
	return nil
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
