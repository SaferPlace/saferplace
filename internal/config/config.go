package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config containing all configuration for the ingest engine
type Config struct {
	Port  int  `envconfig:"PORT" default:"8001"`
	Debug bool `envconfig:"DEBUG"`

	Cert CertConfig
	Auth AuthConfig

	Queue    string `default:"memory"`
	Database string `default:"sql"`
	Notifier string `default:"log"`
}

// CertConfig specifies how the certificates should be created
type CertConfig struct {
	Domains  []string `default:"localhost"`
	Provider string   `default:"insecure"`
}

// AuthConfig is used to configure OAuth
type AuthConfig struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Domain       string `default:"http://localhost:8001"`
}

// Parse the configuration
func Parse(prefix string) (*Config, error) {
	var cfg Config
	if err := envconfig.Process("saferplace", &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	return &cfg, nil
}
