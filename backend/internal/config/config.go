package config

import (
	"os"
	"strings"
)

// Config holds all application configuration gathered from environment variables.
type Config struct {
	Port           string
	DBPath         string
	JWTSecret      string
	GinMode        string
	FrontendDir    string
	AllowedOrigins []string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "fastchem.db"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		GinMode:     getEnv("GIN_MODE", "debug"),
		FrontendDir: os.Getenv("FRONTEND_DIR"), // empty means auto-detect
	}

	originsStr := os.Getenv("ALLOWED_ORIGINS")
	if originsStr != "" {
		for _, o := range strings.Split(originsStr, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, o)
			}
		}
	}
	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:8080",
			"https://fastchem.takumihomelab.works",
		}
	}

	return cfg
}

// IsProduction returns true when GinMode is "release".
func (c *Config) IsProduction() bool {
	return c.GinMode == "release"
}

// AllowedOriginsSet returns a set for O(1) origin lookups.
func (c *Config) AllowedOriginsSet() map[string]bool {
	m := make(map[string]bool, len(c.AllowedOrigins))
	for _, o := range c.AllowedOrigins {
		m[o] = true
	}
	return m
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
