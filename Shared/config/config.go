package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var (
	applicationName string
)

func InitConfiguration(appName string) error {
	applicationName = strings.ToUpper(appName)
	viper.Reset()
	viper.SetEnvPrefix("CAPTAIN")
	viper.AutomaticEnv()
	return nil
}

func GetGlobalString(key string) string {
	return viper.GetString(key)
}

func GetAppString(key string) string {
	return viper.GetString(fmt.Sprintf("%s_%s", applicationName, key))
}

func SetAppString(key string, value string) {
	viper.Set(fmt.Sprintf("%s_%s", applicationName, strings.ToUpper(key)), value)
}
