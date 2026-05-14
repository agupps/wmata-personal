package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Api Api `yaml:"api"`
}

type Api struct {
	Key string `yaml:"key"`
}

func (c *Config) Parse() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal config file, please refer to config-example.yaml: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
