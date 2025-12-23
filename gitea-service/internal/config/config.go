package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the Gitea service
type Config struct {
	// Gitea configuration
	GiteaURL   string `envconfig:"GITEA_URL" required:"true"`
	GiteaToken string `envconfig:"GITEA_TOKEN" required:"true"`

	// LDAP Manager service configuration (for inter-service communication)
	LDAPManagerURL string `envconfig:"LDAP_MANAGER_URL" required:"true"`

	// Server configuration
	Port        int    `envconfig:"PORT" default:"8081"`
	MetricsPort int    `envconfig:"METRICS_PORT" default:"9091"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`

	// JWT configuration (must match LDAP Manager for token validation)
	JWTSecret     string        `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiration time.Duration `envconfig:"JWT_EXPIRATION" default:"24h"`

	// CORS configuration
	CORSOrigins []string `envconfig:"CORS_ORIGINS" default:"*"`

	// Graceful shutdown timeout
	ShutdownTimeout int `envconfig:"SHUTDOWN_TIMEOUT" default:"30"`

	// HTTP client timeouts
	HTTPClientTimeout time.Duration `envconfig:"HTTP_CLIENT_TIMEOUT" default:"30s"`
}

// Load reads configuration from environment variables
func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	return &cfg
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
