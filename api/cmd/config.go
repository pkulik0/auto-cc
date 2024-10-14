package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type config struct {
	IsProduction bool   `mapstructure:"production"`
	Port         uint16 `mapstructure:"port"`

	PostgresHost string `mapstructure:"postgres_host"`
	PostgresPort uint16 `mapstructure:"postgres_port"`
	PostgresUser string `mapstructure:"postgres_user"`
	PostgresPass string `mapstructure:"postgres_pass"`
	PostgresDB   string `mapstructure:"postgres_db"`

	KeycloakURL          string `mapstructure:"keycloak_url"`
	KeycloakRealm        string `mapstructure:"keycloak_realm"`
	KeycloakClientId     string `mapstructure:"keycloak_client_id"`
	KeycloakClientSecret string `mapstructure:"keycloak_client_secret"`

	GoogleCallbackURL string `mapstructure:"google_callback_url"`
}

const (
	envPrefix = "AUTOCC"

	envProd = "PROD"
	envPort = "PORT"

	envPostgresHost = "POSTGRES_HOST"
	envPostgresPort = "POSTGRES_PORT"
	envPostgresUser = "POSTGRES_USER"
	envPostgresPass = "POSTGRES_PASS"
	envPostgresDB   = "POSTGRES_DB"

	envKeycloakURL          = "KEYCLOAK_URL"
	envKeycloakRealm        = "KEYCLOAK_REALM"
	envKeycloakClientId     = "KEYCLOAK_CLIENT_ID"
	envKeycloakClientSecret = "KEYCLOAK_CLIENT_SECRET"

	envGoogleCallbackURL = "GOOGLE_CALLBACK_URL"
)

func bindEnvs(key ...string) error {
	for _, k := range key {
		err := viper.BindEnv(k)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseConfig() (*config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err = bindEnvs(
		envPort, envPort,
		envPostgresHost, envPostgresPort, envPostgresUser, envPostgresPass, envPostgresDB,
		envKeycloakURL, envKeycloakRealm, envKeycloakClientId, envKeycloakClientSecret,
		envGoogleCallbackURL,
	)
	if err != nil {
		return nil, err
	}

	viper.SetDefault(envProd, false)
	viper.SetDefault(envPort, 8080)

	viper.SetDefault(envPostgresHost, "postgres")
	viper.SetDefault(envPostgresPort, 5432)
	viper.SetDefault(envPostgresUser, "autocc")
	viper.SetDefault(envPostgresDB, "autocc")

	if os.Getenv(envPrefix+"_"+envProd) == "" {
		viper.SetDefault(envPostgresPass, "autocc")
		viper.SetDefault(envKeycloakURL, "http://localhost:8081")
		viper.SetDefault(envKeycloakRealm, "autocc")
		viper.SetDefault(envGoogleCallbackURL, "http://localhost:8080/sessions/google/callback")
	} else {
		viper.SetDefault(envKeycloakURL, "https://sso.ony.sh")
		viper.SetDefault(envKeycloakRealm, "onysh")
		viper.SetDefault(envGoogleCallbackURL, "https://autocc.pkulik.com/sessions/google/callback")
	}

	err = viper.ReadInConfig()
	switch err.(type) {
	case nil:
	case viper.ConfigFileNotFoundError:
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
