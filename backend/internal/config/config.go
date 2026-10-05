package config

import (
	"time"
)

// Config is the application-wide configuration loaded from environment.
type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	Security SecurityConfig
	Log      LogConfig
}

type AppConfig struct {
	Env          string
	Name         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BodyLimit    string
}

type DBConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

type SecurityConfig struct {
	FrontendURL     string
	AllowedOrigins  []string
	RateLimitMax    int
	RateLimitWindow time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables. It first loads a
// .env file from the working directory (if present) so local development
// works without exporting variables manually.
func Load() *Config {
	LoadDotEnv(".env")
	return &Config{
		App: AppConfig{
			Env:          Get("APP_ENV", "development"),
			Name:         Get("APP_NAME", "eLabTrack V2"),
			Port:         Get("APP_PORT", "8080"),
			ReadTimeout:  Duration("APP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: Duration("APP_WRITE_TIMEOUT", 15*time.Second),
			BodyLimit:    Get("APP_BODY_LIMIT", "1MB"),
		},
		DB: DBConfig{
			Host:            Get("DB_HOST", "127.0.0.1"),
			Port:            Get("DB_PORT", "5432"),
			Name:            Get("DB_NAME", "elabtrack_v2"),
			User:            Get("DB_USER", "postgres"),
			Password:        Get("DB_PASSWORD", ""),
			MaxConns:        int32(Int("DB_MAX_CONNS", 20)),
			MinConns:        int32(Int("DB_MIN_CONNS", 2)),
			MaxConnLifetime: Duration("DB_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: Duration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute),
		},
		JWT: JWTConfig{
			Secret:     Get("JWT_SECRET", "change-me-in-production"),
			AccessTTL:  Duration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL: Duration("JWT_REFRESH_TTL", 168*time.Hour),
			Issuer:     Get("JWT_ISSUER", "elabtrack-v2"),
		},
		Security: SecurityConfig{
			FrontendURL:     Get("FRONTEND_URL", "http://localhost:5173"),
			AllowedOrigins:  List("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://localhost:4173"}),
			RateLimitMax:    Int("RATE_LIMIT_MAX", 120),
			RateLimitWindow: Duration("RATE_LIMIT_WINDOW", time.Minute),
		},
		Log: LogConfig{
			Level:  Get("LOG_LEVEL", "info"),
			Format: Get("LOG_FORMAT", "console"),
		},
	}
}

// IsProduction reports whether the app runs in production.
func (c *Config) IsProduction() bool { return c.App.Env == "production" }

// IsDevelopment reports whether the app runs in development.
func (c *Config) IsDevelopment() bool { return c.App.Env == "development" }
