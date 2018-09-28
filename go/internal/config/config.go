package config

import (
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const (
	configFile = "config.yml"
	logFile    = "gitlab-shell.log"
)

type MigrationConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Features []string `yaml:"features"`
}

type Config struct {
	RootDir   string
	LogFile   string          `yaml:"log_file"`
	LogFormat string          `yaml:"log_format"`
	Migration MigrationConfig `yaml:"migration"`
}

func New() (*Config, error) {
	cfg := Config{}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg.RootDir = dir

	configBytes, err := ioutil.ReadFile(path.Join(cfg.RootDir, configFile))
	if err != nil {
		return nil, err
	}

	if err := parseConfig(configBytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// parseConfig expects YAML data in configBytes and a Config instance with RootDir set.
func parseConfig(configBytes []byte, cfg *Config) error {
	if err := yaml.Unmarshal(configBytes, cfg); err != nil {
		return err
	}

	if cfg.LogFile == "" {
		cfg.LogFile = logFile
	}

	if len(cfg.LogFile) > 0 && cfg.LogFile[0] != '/' {
		cfg.LogFile = path.Join(cfg.RootDir, cfg.LogFile)
	}

	if cfg.LogFormat == "" {
		cfg.LogFormat = "text"
	}

	return nil
}
