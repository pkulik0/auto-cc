package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type config struct {
	Port int16 `mapstructure:"port"`
}

func parseConfig() (*config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("AUTOCC")
	viper.BindEnv("PORT")

	viper.SetDefault("PORT", 8080)

	err := viper.ReadInConfig()
	switch err.(type) {
	case nil:
		log.Info().Msg("Using config file")
	case viper.ConfigFileNotFoundError:
		log.Info().Msg("Config file not found")
	default:
		return nil, err
	}

	var c config
	err = viper.Unmarshal(&c)
	return &c, err
}
