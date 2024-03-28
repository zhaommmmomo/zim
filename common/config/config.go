package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	ENV      = "env"
	ENV_DEV  = "dev"
	ENV_PROD = "prod"
	LOG_PATH = "log.path"
)

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("load config {%s} fail: %w", path, err))
	}
	viper.WatchConfig()
}

func GetLogPath() string {
	return viper.GetString(LOG_PATH)
}

func Get(key string) string {
	return viper.GetString(key)
}

func SetDefault(key string, value any) {
	viper.SetDefault(key, value)
}

func IsDebug() bool {
	return viper.GetString(ENV) == ENV_DEV
}

func IsProd() bool {
	return viper.GetString(ENV) == ENV_PROD
}
