package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	env     = "env"
	envDev  = "dev"
	envProd = "prod"
	logPath = "log.path"
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
	return viper.GetString(logPath)
}

func Get(key string) string {
	return viper.GetString(key)
}

func SetDefault(key string, value any) {
	viper.SetDefault(key, value)
}

func IsDebug() bool {
	return viper.GetString(env) == envDev
}

func IsProd() bool {
	return viper.GetString(env) == envProd
}
