package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"safer.place/internal/database/sqldatabase"
	"safer.place/internal/storage/minio"
)

// Config containing all configuration for saferplace.
type Config struct {
	// File from which the config was parsed. Empty if the file was not used.
	File  string
	Debug bool `yaml:"debug" envconfig:"DEBUG"`

	Webserver WebserverConfig `yaml:"webserver"`
	Queue     QueueConfig     `yaml:"queue"`
	Database  DatabaseConfig  `yaml:"database"`
	Storage   StorageConfig   `yaml:"storage"`
	Notifier  NotifierConfig  `yaml:"notifier"`
}

// WebserverConfig contains all configuration used to setup the webserver and middleware
type WebserverConfig struct {
	Port        int        `yaml:"port" envconfig:"PORT" default:"8001"`
	CORSDomains []string   `yaml:"cors_domains" default:""`
	Cert        CertConfig `yaml:"cert"`
	Auth        AuthConfig `yaml:"auth"`
}

// QueueConfig provides the configuration to consume and produce from the queue
type QueueConfig struct {
	Provider string `yaml:"provider" default:"memory"`
}

// DatabaseConfig configures the database used as a backend for all incident data.
type DatabaseConfig struct {
	Provider string `yaml:"provider" default:"sql"`

	SQL sqldatabase.Config `yaml:"sql"`
}

// StorageConfig configures the storage for user uploads.
type StorageConfig struct {
	Provider string `yaml:"provider" default:"minio"`

	Minio minio.Config `yaml:"minio"`
}

// Notifier can be configured to notify a third party of a incident.
type NotifierConfig struct {
	Provider string `yaml:"provider" default:"log"`
}

// CertConfig specifies how the certificates should be created
type CertConfig struct {
	Provider string   `default:"insecure"`
	Domains  []string `default:"localhost"`
}

// AuthConfig is used to configure OAuth
type AuthConfig struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	Domain       string `default:"http://localhost:8001"`
}

// Parse the configuration from a specific file. We first load the configuration from the
// environment variables, and then override with values in a file.
func Parse(file string) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("saferplace", cfg); err != nil {
		return nil, fmt.Errorf("unable to read from environment: %w", err)
	}

	f, err := os.Open(file)
	if err != nil {
		// Skip file decoding if it doesn't exist.
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, fmt.Errorf("unable to open config file: %w", err)
	}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode: %w", err)
	}
	cfg.File = file

	return cfg, nil
}
