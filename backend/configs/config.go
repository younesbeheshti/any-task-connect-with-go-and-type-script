package configs

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Server   ServerConfig   `mapstructure:"server"`
	Platform PlatformConfig `mapstructure:"platform"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

type PlatformConfig struct {
	CommissionPercent int64 `mapstructure:"commission_percent"`
}

// StorageConfig configures file upload storage.
type StorageConfig struct {
	LocalDir string `mapstructure:"local_dir"` // directory where uploaded files are written
	MaxSize  int64  `mapstructure:"max_size"`  // max single-file size in bytes
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
	Port        string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	URL      string `mapstructure:"url"`
	Exchange string `mapstructure:"exchange"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	AccessTTL  int    `mapstructure:"access_ttl"`
	RefreshTTL int    `mapstructure:"refresh_ttl"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type ServerConfig struct {
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// Load reads configuration from config.yaml, .env, and environment variables.
func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/etc/taskbridge")

	setDefaults(v)
	bindEnv(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	v.SetConfigFile(".env")
	_ = v.MergeInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "taskbridge")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.port", "8080")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.user", "taskbridge")
	v.SetDefault("database.password", "taskbridge")
	v.SetDefault("database.name", "taskbridge")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("rabbitmq.url", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.exchange", "taskbridge.events")

	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.access_ttl", 900)
	v.SetDefault("jwt.refresh_ttl", 604800)

	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "15s")
	v.SetDefault("server.idle_timeout", "60s")

	v.SetDefault("platform.commission_percent", 8)

	v.SetDefault("storage.local_dir", "./uploads")
	v.SetDefault("storage.max_size", 52428800) // 50 MiB
}

func bindEnv(v *viper.Viper) {
	_ = v.BindEnv("app.name", "APP_NAME")
	_ = v.BindEnv("app.environment", "APP_ENVIRONMENT")
	_ = v.BindEnv("app.port", "APP_PORT")

	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")
	_ = v.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	_ = v.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	_ = v.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = v.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")

	_ = v.BindEnv("redis.host", "REDIS_HOST")
	_ = v.BindEnv("redis.port", "REDIS_PORT")
	_ = v.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = v.BindEnv("redis.db", "REDIS_DB")

	_ = v.BindEnv("rabbitmq.url", "RABBITMQ_URL")
	_ = v.BindEnv("rabbitmq.exchange", "RABBITMQ_EXCHANGE")

	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("jwt.access_ttl", "JWT_ACCESS_TTL")
	_ = v.BindEnv("jwt.refresh_ttl", "JWT_REFRESH_TTL")

	_ = v.BindEnv("platform.commission_percent", "PLATFORM_COMMISSION_PERCENT")

	_ = v.BindEnv("storage.local_dir", "STORAGE_LOCAL_DIR")
	_ = v.BindEnv("storage.max_size", "STORAGE_MAX_SIZE")
}

// Validate checks required configuration values.
func (c *Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("app.port is required")
	}
	if c.Database.Host == "" || c.Database.Name == "" {
		return fmt.Errorf("database host and name are required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.App.Environment == "production" && c.JWT.Secret == "change-me-in-production" {
		return fmt.Errorf("jwt.secret must be changed in production")
	}
	env := strings.ToLower(c.App.Environment)
	switch env {
	case "development", "staging", "production":
	default:
		return fmt.Errorf("invalid app.environment: %s", c.App.Environment)
	}
	return nil
}

// IsDevelopment reports whether the app runs in development mode.
func (c *Config) IsDevelopment() bool {
	return strings.EqualFold(c.App.Environment, "development")
}

// IsProduction reports whether the app runs in production mode.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.App.Environment, "production")
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Addr returns the Redis address host:port.
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// AccessTokenDuration returns the access token TTL.
func (c *JWTConfig) AccessTokenDuration() time.Duration {
	return time.Duration(c.AccessTTL) * time.Second
}

// RefreshTokenDuration returns the refresh token TTL.
func (c *JWTConfig) RefreshTokenDuration() time.Duration {
	return time.Duration(c.RefreshTTL) * time.Second
}
