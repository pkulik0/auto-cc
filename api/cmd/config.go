package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type config struct {
	Port uint16 `mapstructure:"port"`

	RedisURL string `mapstructure:"redis_url"`

	PostgresHost string `mapstructure:"postgres_host"`
	PostgresPort uint16 `mapstructure:"postgres_port"`
	PostgresUser string `mapstructure:"postgres_user"`
	PostgresPass string `mapstructure:"postgres_pass"`
	PostgresDB   string `mapstructure:"postgres_db"`

	KeycloakURL          string `mapstructure:"keycloak_url"`
	KeycloakRealm        string `mapstructure:"keycloak_realm"`
	KeycloakClientId     string `mapstructure:"keycloak_client_id"`
	KeycloakClientSecret string `mapstructure:"keycloak_client_secret"`
}

func parseConfig() (*config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("AUTOCC")
	viper.BindEnv("PORT")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("POSTGRES_HOST")
	viper.BindEnv("POSTGRES_PORT")
	viper.BindEnv("POSTGRES_USER")
	viper.BindEnv("POSTGRES_PASS")
	viper.BindEnv("POSTGRES_DB")
	viper.BindEnv("KEYCLOAK_URL")
	viper.BindEnv("KEYCLOAK_REALM")
	viper.BindEnv("KEYCLOAK_CLIENT_ID")
	viper.BindEnv("KEYCLOAK_CLIENT_SECRET")

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "autocc")
	viper.SetDefault("POSTGRES_PASS", "autocc")
	viper.SetDefault("POSTGRES_DB", "autocc")
	viper.SetDefault("KEYCLOAK_URL", "https://sso.ony.sh")
	viper.SetDefault("KEYCLOAK_REALM", "onysh")

	err := viper.ReadInConfig()
	switch err.(type) {
	case nil:
		log.Debug().Msg("Using config file")
	case viper.ConfigFileNotFoundError:
		log.Debug().Msg("Config file not found")
	default:
		return nil, err
	}

	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
