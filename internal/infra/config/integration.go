package config

import (
	"github.com/spf13/viper"
)

// PaymentProvider represents a single payment provider configuration
type PaymentProvider struct {
	Name       string           `mapstructure:"name"`
	Type       string           `mapstructure:"type"`
	BaseURL    string           `mapstructure:"base_url"`
	AuthHeader string           `mapstructure:"auth_header"`
	AuthToken  string           `mapstructure:"auth_token"`
	Currency   string           `mapstructure:"currency"`
	Endpoints  []EndpointConfig `mapstructure:"endpoints"`
}

// EndpointConfig represents the configuration for an endpoint of a payment provider
type EndpointConfig struct {
	Action string            `mapstructure:"action"`
	Method string            `mapstructure:"method"`
	Path   string            `mapstructure:"path"`
	Params map[string]string `mapstructure:"params"`
}

// IntegrationConfig represents the entire configuration structure
type IntegrationConfig struct {
	PaymentProviders []PaymentProvider `mapstructure:"payment_providers"`
}

func LoadIntegrationConfig(path string) (*IntegrationConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config IntegrationConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
