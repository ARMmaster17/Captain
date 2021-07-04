package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var (
	// ApplicationName is the globally-recognized name of the currently running application.
	ApplicationName string
)

// InitConfiguration Initializes the static environment to read configuration files for the currently running service.
func InitConfiguration(appName string) {
	ApplicationName = strings.ToUpper(appName)
	viper.Reset()
	viper.SetEnvPrefix("CAPTAIN")
	viper.AutomaticEnv()
}

// GetGlobalString Gets a string with the CAPTAIN_ prefix from any configured datasource.
func GetGlobalString(key string) string {
	return viper.GetString(key)
}

// GetAppString Gets a string with the CAPTAIN_<AppName> prefix where AppName is the configured application name.
func GetAppString(key string) string {
	return viper.GetString(fmt.Sprintf("%s_%s", ApplicationName, key))
}

// SetAppString Sets an application-specific configuration setting in memory.
func SetAppString(key string, value string) {
	viper.Set(fmt.Sprintf("%s_%s", ApplicationName, strings.ToUpper(key)), value)
}

// GetAppInt Gets an integer with the CAPTAIN_<AppName> prefix from any configured datasource.
func GetAppInt(key string) int {
	return viper.GetInt(fmt.Sprintf("%s_%s", ApplicationName, key))
}

// GetGlobalInt gets an integer with the CAPTAIN_ prefix from any configured datasource.
func GetGlobalInt(key string) int {
	return viper.GetInt(fmt.Sprintf("%s_%s", ApplicationName, key))
}
