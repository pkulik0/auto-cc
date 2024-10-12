package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type config struct {
	Port                 int16  `mapstructure:"port"`
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
	viper.BindEnv("KEYCLOAK_URL")
	viper.BindEnv("KEYCLOAK_REALM")
	viper.BindEnv("KEYCLOAK_CLIENT_ID")
	viper.BindEnv("KEYCLOAK_CLIENT_SECRET")

	viper.SetDefault("PORT", 8080)
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
