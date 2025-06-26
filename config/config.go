package config

import (
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	URI  string `mapstructure:"uri"`
	Name string `mapstructure:"name"`
}

type TokenConfig struct {
	SymmetricKey string        `mapstructure:"symmetricKey"`
	Duration     time.Duration `mapstructure:"duration"`
}

type Config struct {
	DB            DBConfig    `mapstructure:"db"`
	Token         TokenConfig `mapstructure:"token"`
	ServerAddress string      `mapstructure:"serverAddress"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	// Read from file
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("json")

	// Overwrite file config with environment variable if exists.
	viper.AutomaticEnv()
	// 因為使用 JSON 格式指定環境變數，因應會從外部修改值，讓 docker 指令的 -e
	// 也可以用 DB_URI 的方式給進來
	viper.BindEnv("db.uri", "DB_URI")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
