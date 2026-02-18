package types

import (
	"fmt"

	_ "github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment string `yaml:"environment" envconfig:"ENVIRONMENT"`
	Port        uint   `yaml:"port" envconfig:"PORT"`

	Database DatabaseConfig `yaml:"database"`
	NTFY     NTFYConfig     `yaml:"ntfy"`
}

type NTFYConfig struct {
	URL   string `yaml:"url" envconfig:"NTFY_URL"`
	Topic string `yaml:"topic" envconfig:"NTFY_TOPIC"`
}

type DatabaseConfig struct {
	URL      string `yaml:"url" envconfig:"DATABASE_URL"`
	Host     string `yaml:"host" envconfig:"DB_HOST"`
	Port     int    `yaml:"port" envconfig:"DB_PORT"`
	User     string `yaml:"user" envconfig:"DB_USER"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	Database string `yaml:"database" envconfig:"DB_NAME"`
	SSLMode  string `yaml:"sslmode" envconfig:"DB_SSLMODE"`
}

// BuildDatabaseURL constructs connection string from components if URL not provided
func (d *DatabaseConfig) BuildDatabaseURL() string {
	if d.URL != "" {
		return d.URL
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode)
}
