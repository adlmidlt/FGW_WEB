package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	PSQL PostgresqlConfig `yaml:"PSQL"`
}

type PostgresqlConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	SSLMode         string        `yaml:"SSLMode"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
}

// LoadConfigDatabase - загружает конфигурацию настройки БД из YAML-файла.
func (c *Config) LoadConfigDatabase(pathToYamlFile string) error {
	return c.readFromYaml(pathToYamlFile)
}

// readFromYaml - чтение конфигурационных настроек из YAML-файла.
func (c *Config) readFromYaml(pathToYamlFile string) error {
	yamlFile, err := os.ReadFile(pathToYamlFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	return nil
}
