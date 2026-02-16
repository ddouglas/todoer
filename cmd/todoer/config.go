package main

import (
	"fmt"
	"os"

	"github.com/ddouglas/todoer/pkg/types"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var c *types.Config = new(types.Config)

func LoadConfig(configPath string) error {
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}

		err = yaml.Unmarshal(data, c)
		if err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}

	}
	err := envconfig.Process("TODOER", c)
	if err != nil {
		return fmt.Errorf("failed to process env vars: %w", err)
	}

	err = validateConfig()
	if err != nil {
		return err

	}

	if c.Database.URL == "" {
		c.Database.URL = c.Database.BuildDatabaseURL()
	}

	return nil
}

func validateConfig() error {
	if c.Port == 0 {
		return fmt.Errorf("port is required")
	}

	if c.Database.URL == "" && c.Database.Host == "" {
		return fmt.Errorf("database configuration is required")
	}

	return nil

}
