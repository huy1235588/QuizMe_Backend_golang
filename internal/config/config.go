package config

import (
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Cloudinary CloudinaryConfig `mapstructure:"cloudinary"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug or release
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret              string `mapstructure:"secret"`
	ExpirationMs        int64  `mapstructure:"expiration_ms"`
	RefreshExpirationMs int64  `mapstructure:"refresh_expiration_ms"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	Path string `mapstructure:"path"`
}

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string            `mapstructure:"cloud_name"`
	APIKey    string            `mapstructure:"api_key"`
	APISecret string            `mapstructure:"api_secret"`
	BaseURL   string            `mapstructure:"base_url"`
	Folder    map[string]string `mapstructure:"folder"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.expiration_ms", 86400000)        // 1 day
	viper.SetDefault("jwt.refresh_expiration_ms", 2592000000) // 30 days
	viper.SetDefault("upload.path", "./uploads")

	// Read environment variables
	viper.AutomaticEnv()

	// Environment variable mappings
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("cloudinary.cloud_name", "CLOUDINARY_CLOUD_NAME")
	viper.BindEnv("cloudinary.api_key", "CLOUDINARY_API_KEY")
	viper.BindEnv("cloudinary.api_secret", "CLOUDINARY_API_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is not an error if env vars are set
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
