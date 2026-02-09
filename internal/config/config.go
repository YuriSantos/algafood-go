package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Email    EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Name      string `mapstructure:"name"`
	Charset   string `mapstructure:"charset"`
	ParseTime bool   `mapstructure:"parseTime"`
	Loc       string `mapstructure:"loc"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.ParseTime, d.Loc)
}

type JWTConfig struct {
	JWKSURL   string `mapstructure:"jwks_url"`
	Issuer    string `mapstructure:"issuer"`
	SecretKey string `mapstructure:"secret_key"`
}

type StorageConfig struct {
	Type  string             `mapstructure:"type"`
	Local LocalStorageConfig `mapstructure:"local"`
	S3    S3StorageConfig    `mapstructure:"s3"`
}

type LocalStorageConfig struct {
	Directory string `mapstructure:"directory"`
}

type S3StorageConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Directory string `mapstructure:"directory"`
}

type EmailConfig struct {
	Type    string             `mapstructure:"type"`
	From    string             `mapstructure:"from"`
	Sandbox SandboxEmailConfig `mapstructure:"sandbox"`
	SMTP    SMTPEmailConfig    `mapstructure:"smtp"`
}

type SandboxEmailConfig struct {
	Recipient string `mapstructure:"recipient"`
}

type SMTPEmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variable override
	viper.SetEnvPrefix("ALGAFOOD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
