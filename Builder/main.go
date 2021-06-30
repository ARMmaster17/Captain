package Builder

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

func main() {
	initLogging()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain/atc")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = generateConfigFile()
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("unable to create config file")
			} else {
				log.Warn().Msg("a config file was generated at /etc/captain/atc/config.yaml that needs to be edited before starting Captain again")
			}
		} else {
			log.Fatal().Stack().Err(err).Msg("unable to read from config file")
		}
		return
	}

	zerolog.SetGlobalLevel(getLogLevel(viper.GetInt("config.loglevel")))
}

func initLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	viper.SetDefault("config.loglevel", -1)
}

// generateConfigFile creates a config file when one doesn't exist. Actual implementation is handled by Viper.
func generateConfigFile() error {
	viper.Set("defaults.publickey", "x")
	viper.Set("defaults.network.nameservers", "8.8.8.8 8.8.4.4")
	viper.Set("defaults.network.searchdomain", "")
	viper.Set("defaults.network.gateway", "10.1.0.1")
	viper.Set("defaults.network.mtu", 1450)
	viper.Set("defaults.network.cidr", 16)
	viper.Set("defaults.network.blocks", []string{"10.1.5.0/24"})
	viper.Set("defaults.image", "debian-10")
	// In the future when there is more than one driver, this section should not
	// be added automatically.
	viper.Set("config.drivers.provisioners", []string{"proxmoxlxc"})
	viper.Set("config.api.port", 5000)
	viper.Set("images.debian-10.proxmoxlxc", "pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz")
	viper.Set("drivers.provisioners.proxmoxlxc.publicnetwork", "internal")
	viper.Set("drivers.provisioners.proxmoxlxc.diskstorage", "pve-storage")
	viper.Set("drivers.provisioners.proxmoxlxc.defaultnode", "pxvh1")
	viper.Set("drivers.provisioners.proxmoxlxc.forcessl", false)
	viper.Set("drivers.provisioners.proxmoxlxc.url", "https://192.168.1.241:8006/api2/json")

	return viper.WriteConfig()
}

func getLogLevel(level int) zerolog.Level {
	switch l := level; l {
	case 5:
		return zerolog.PanicLevel
	case 4:
		return zerolog.FatalLevel
	case 3:
		return zerolog.ErrorLevel
	case 2:
		return zerolog.WarnLevel
	case 1:
		return zerolog.InfoLevel
	case 0:
		return zerolog.DebugLevel
	default:
		return zerolog.TraceLevel
	}
}