package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"SERVER_HOST"`
		Port int    `mapstructure:"SERVER_PORT"`
	}
	DB struct {
		User    string `mapstructure:"DB_USER"`
		Pass    string `mapstructure:"DB_PASS"`
		Host    string `mapstructure:"DB_HOST"`
		Port    int    `mapstructure:"DB_PORT"`
		Name    string `mapstructure:"DB_NAME"`
		SSLMode string `mapstructure:"SSL_MODE"`
	}
	Logger struct {
		Level string `mapstructure:"LOGGER_LEVEL"`
	}
	Blizzard struct {
		ClientID     string `mapstructure:"CLIENT_ID"`
		ClientSecret string `mapstructure:"CLIENT_SECRET"`
		RedirectURL  string `mapstructure:"REDIRECT_URL"`
	}
	JWT struct {
		Secret string `mapstructure:"JWT_SECRET"`
	}
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.SetConfigType("env")
	v.SetConfigName(".env")

	v.SetDefault("SERVER_HOST", "localhost")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("LOGGER_LEVEL", "info")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed read config: %v", err)
		return nil, err
	}

	v.AutomaticEnv()

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed mapping env to config struct: %v", err)
	}

	return &cfg, nil
}
