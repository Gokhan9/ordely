package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource            string `mapstructure:"DB_SOURCE"`
	RedisAddress        string `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress   string `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey   string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration string `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
