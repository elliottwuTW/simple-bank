package config

import (
	"github.com/spf13/viper"
)

type DBConfig struct {
	URI  string `mapstructure:"uri"`
	Name string `mapstructure:"name"`
}

type Config struct {
	DB            DBConfig `mapstructure:"db"`
	ServerAddress string   `mapstructure:"serverAddress"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	// Read from file
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("json")

	// Overwrite file config with environment variable if exists.
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
