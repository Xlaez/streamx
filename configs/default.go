package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBUri               string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri            string        `mapstructure:"REDIS_URL"`
	Port                string        `mapstructure:"PORT"`
	TokenKey            string        `mapstructure:"TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	CloudinaryEnv       string        `mapstructure:"CLOUDINARY_API_ENV"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
