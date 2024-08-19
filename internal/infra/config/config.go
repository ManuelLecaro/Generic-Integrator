package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName    string `mapstructure:"APP_NAME"`
	Port       string `mapstructure:"PORT"`
	Env        string `mapstructure:"ENV"`
	APIKey     string `mapstructure:"API_KEY"`
	DB         DBConfig
	EventStore EventStoreConfig
}

type DBConfig struct {
	Host             string `mapstructure:"DB_HOST"`
	Port             string `mapstructure:"DB_PORT"`
	User             string `mapstructure:"DB_USER"`
	Password         string `mapstructure:"DB_PASSWORD"`
	Name             string `mapstructure:"DB_NAME"`
	SSLMode          string `mapstructure:"DB_SSLMODE"`
	ConnectionString string `mapstructure:"DB_CONNECTSTRING"`
}

type EventStoreConfig struct {
	ConnectionString string `mapstructure:"EVENTSTORE_DB_CONNECTION_STRING"`
}

var config *Config

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")

	viper.SetDefault("APP_NAME", "PaymentPlatform")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("API_KEY", "test")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "user")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "agap")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("EVENTSTORE_DB_CONNECTION_STRING", "esdb://localhost:2113?tls=false")
	viper.SetDefault("DB_CONNECTSTRING", "mongodb://root:example@localhost:27017/agap?timeoutMS=5000")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error leyendo el archivo de configuración, utilizando valores por defecto: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config no ha sido cargada. Llama a LoadConfig() primero.")
	}
	return config
}
