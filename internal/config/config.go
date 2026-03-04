package config

import (
	"github.com/spf13/viper"
)

type Service struct {
	Name   string `mapstructure:"name"`
	Prefix string `mapstructure:"prefix"`
	Target string `mapstructure:"target"`
	RequiresAuth bool `mapstructure:"requires_auth"` // New field
	RateLimit    float64 `mapstructure:"rate_limit"` // Per-service limit
}

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Services []Service `mapstructure:"services"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
