package config

import "github.com/spf13/viper"

func LoadConfig() error {
	viper.AddConfigPath("./pkg/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
