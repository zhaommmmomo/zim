package config

import (
	"fmt"
	"github.com/spf13/viper"
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

func Get(key string) string {
	return viper.GetString(key)
}
