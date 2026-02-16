package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Auth        AuthConfig        `mapstructure:"auth"`
	Storage     StorageConfig     `mapstructure:"storage"`
	Email       EmailConfig       `mapstructure:"email"`
	EventBridge EventBridgeConfig `mapstructure:"eventbridge"`
	SQS         SQSConfig         `mapstructure:"sqs"`
	AWS         AWSConfig         `mapstructure:"aws"`
	SpringDoc   SpringDocConfig   `mapstructure:"springdoc"`
}

type ServerConfig struct {
	Port                   int    `mapstructure:"port"`
	Mode                   string `mapstructure:"mode"`
	CompressionEnabled     bool   `mapstructure:"compression_enabled"`
	ForwardHeadersStrategy string `mapstructure:"forward_headers_strategy"`
}

type DatabaseConfig struct {
	Host                     string `mapstructure:"host"`
	Port                     int    `mapstructure:"port"`
	User                     string `mapstructure:"user"`
	Password                 string `mapstructure:"password"`
	Name                     string `mapstructure:"name"`
	Charset                  string `mapstructure:"charset"`
	ParseTime                bool   `mapstructure:"parseTime"`
	Loc                      string `mapstructure:"loc"`
	CreateDatabaseIfNotExist bool   `mapstructure:"create_database_if_not_exist"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.ParseTime, d.Loc)
}

type JWTConfig struct {
	JWKSURL   string         `mapstructure:"jwks_url"`
	Issuer    string         `mapstructure:"issuer"`
	SecretKey string         `mapstructure:"secret_key"`
	Keystore  KeystoreConfig `mapstructure:"keystore"`
}

type KeystoreConfig struct {
	JKSLocation  string `mapstructure:"jks_location"`
	Password     string `mapstructure:"password"`
	KeypairAlias string `mapstructure:"keypair_alias"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AuthConfig struct {
	ProviderURL string            `mapstructure:"provider_url"`
	OpaqueToken OpaqueTokenConfig `mapstructure:"opaque_token"`
}

type OpaqueTokenConfig struct {
	IntrospectionURI string `mapstructure:"introspection_uri"`
	ClientID         string `mapstructure:"client_id"`
	ClientSecret     string `mapstructure:"client_secret"`
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
	SES     SESEmailConfig     `mapstructure:"ses"`
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

type SESEmailConfig struct {
	Region string `mapstructure:"region"`
}

type EventBridgeConfig struct {
	Type         string `mapstructure:"type"`
	Region       string `mapstructure:"region"`
	EventBusName string `mapstructure:"event_bus_name"`
	Source       string `mapstructure:"source"`
}

type SQSConfig struct {
	Type              string `mapstructure:"type"`
	Region            string `mapstructure:"region"`
	QueueURL          string `mapstructure:"queue_url"`
	MaxMessages       int    `mapstructure:"max_messages"`
	WaitTimeSeconds   int    `mapstructure:"wait_time_seconds"`
	VisibilityTimeout int    `mapstructure:"visibility_timeout"`
}

type AWSConfig struct {
	EndpointURL string               `mapstructure:"endpoint_url"`
	Region      string               `mapstructure:"region"`
	Credentials AWSCredentialsConfig `mapstructure:"credentials"`
}

type AWSCredentialsConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
}

type SpringDocConfig struct {
	OAuthFlow      OAuthFlowConfig `mapstructure:"oauth_flow"`
	SwaggerUI      SwaggerUIConfig `mapstructure:"swagger_ui"`
	PackagesToScan string          `mapstructure:"packages_to_scan"`
	PathsToMatch   string          `mapstructure:"paths_to_match"`
	EnableHateoas  bool            `mapstructure:"enable_hateoas"`
}

type OAuthFlowConfig struct {
	AuthorizationURL string `mapstructure:"authorization_url"`
	TokenURL         string `mapstructure:"token_url"`
}

type SwaggerUIConfig struct {
	OAuth SwaggerUIOAuthConfig `mapstructure:"oauth"`
}

type SwaggerUIOAuthConfig struct {
	ClientID                                  string `mapstructure:"client_id"`
	ClientSecret                              string `mapstructure:"client_secret"`
	UsePKCEWithAuthorizationCodeGrant         bool   `mapstructure:"use_pkce_with_authorization_code_grant"`
	UseBasicAuthenticationWithAccessCodeGrant bool   `mapstructure:"use_basic_authentication_with_access_code_grant"`
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
