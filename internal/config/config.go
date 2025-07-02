package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ErrMissingGraphqlEndpoint = errors.New("missing GraphQL endpoint in configuration")
	ErrMissingCronMode        = errors.New("missing cron mode in configuration (should be 'fix' or 'window')")
	ErrInvalidCronMode        = errors.New("invalid cron mode, should be 'fix' or 'window'")
	ErrEULANotAccepted        = errors.New("EULA not accepted, please set EulaAccepted to true in the configuration (config.yaml)")
)

func GetConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	viper.AddConfigPath("/app/")
	viper.AddConfigPath("/config/")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	if !config.EulaAccepted {
		fmt.Print(EULA)
		return Config{}, ErrEULANotAccepted
	}

	if config.GraphqlEndpoint == "" {
		return Config{}, ErrMissingGraphqlEndpoint
	}

	if config.Cron.Mode == "" {
		return Config{}, ErrMissingCronMode
	}

	if config.Cron.Mode != "fix" && config.Cron.Mode != "window" {
		return Config{}, ErrInvalidCronMode
	}

	return config, nil
}
