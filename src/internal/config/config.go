package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Environment string        `mapstructure:"environment"`
	Server      ServerConfig  `mapstructure:"server"`
	Database    DatabaseConfig `mapstructure:"database"`
	Auth        AuthConfig    `mapstructure:"auth"`
	CORS        CORSConfig    `mapstructure:"cors"`
	Logging     LoggingConfig `mapstructure:"logging"`
	Services    ServicesConfig `mapstructure:"services"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	JWTSecret           string        `mapstructure:"jwt_secret"`
	AccessTokenExpiry   time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry  time.Duration `mapstructure:"refresh_token_expiry"`
	PasswordHashCost    int           `mapstructure:"password_hash_cost"`
}

// CORSConfig represents the CORS configuration
type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	TimeFormat string `mapstructure:"time_format"`
}

// ServicesConfig represents the microservices configuration
type ServicesConfig struct {
	DataService      ServiceConfig `mapstructure:"data_service"`
	AuthService      ServiceConfig `mapstructure:"auth_service"`
	AnalyticsService ServiceConfig `mapstructure:"analytics_service"`
}

// ServiceConfig represents a microservice configuration
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	TLS  bool   `mapstructure:"tls"`
}

// LoadConfig loads the application configuration
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/go-data-api")
	
	// Set default values
	setDefaults()
	
	// Read environment variables
	viper.AutomaticEnv()
	
	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found, use environment variables and defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "5s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.shutdown_timeout", "5s")
	
	// Database defaults
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.database", "data_api")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	
	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "your-secret-key")
	viper.SetDefault("auth.access_token_expiry", "15m")
	viper.SetDefault("auth.refresh_token_expiry", "7d")
	viper.SetDefault("auth.password_hash_cost", 10)
	
	// CORS defaults
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.expose_headers", []string{"Content-Length"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", "12h")
	
	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.time_format", time.RFC3339)
	
	// Services defaults
	viper.SetDefault("services.data_service.host", "localhost")
	viper.SetDefault("services.data_service.port", 50051)
	viper.SetDefault("services.data_service.tls", false)
	
	viper.SetDefault("services.auth_service.host", "localhost")
	viper.SetDefault("services.auth_service.port", 50052)
	viper.SetDefault("services.auth_service.tls", false)
	
	viper.SetDefault("services.analytics_service.host", "localhost")
	viper.SetDefault("services.analytics_service.port", 50053)
	viper.SetDefault("services.analytics_service.tls", false)
}

