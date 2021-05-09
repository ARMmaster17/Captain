package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

var (
	Version string
)

// Main entry point of the application. Handles the creation of the requested number of workers for each task, and sets
// them up to use pipe-based IPC or an external MQ service for communication.
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// TODO: Read verbosity level from command line args

	log.Info().Msg(fmt.Sprintf("Captain %s is starting up", getApplicationVersion()))

	bootstrapOnly := flag.Bool("boostrap", false, "Runs a stripped-down version of Captain to build the entire Captain stack from a single worker node.")

	if *bootstrapOnly {
		err := BootstrapCluster()
		if err != nil {
			log.Error().Msg("An error occurred: " + err.Error())
		}
		return
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = generateConfigFile()
			if err != nil {
				log.Fatal().Stack().Err(err).Msg("unable to create config file")
			} else {
				log.Warn().Msg("a config file was generated at /etc/captain/config.yaml that needs to be edited before starting Captain again")
			}
			return
		} else {
			log.Fatal().Stack().Err(err).Msg("unable to read from config file")
		}
	}

	apiServer := &APIServer{}
	err := apiServer.Start()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Captain API server has fatally crashed")
		return
	}
	apiServer.Serve(viper.GetInt("config.api.port"))

	err = StartMonitoring()
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Captain has fatally crashed")
		return
	}
}

// There is a bug in how 'go test' is implemented. This method does not
// have a unit test.
func getApplicationVersion() string {
	return Version
}

func generateConfigFile() error {
	viper.Set("defaults.publickey", "x")
	viper.Set("defaults.network.nameservers", "8.8.8.8 8.8.4.4")
	viper.Set("defaults.network.searchdomain", "")
	viper.Set("defaults.network.gateway", "10.1.0.1")
	viper.Set("defaults.network.mtu", 1450)
	// In the future when there is more than one driver, this section should not
	// be added automatically.
	viper.Set("config.drivers.provisioners", []string{"proxmoxlxc"})
	viper.Set("config.api.port", 5000)
	viper.Set("drivers.provisioners.proxmoxlxc.image", "pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz")
	viper.Set("drivers.provisioners.proxmoxlxc.publicnetwork", "internal")
	viper.Set("drivers.provisioners.proxmoxlxc.diskstorage", "pve-storage")
	viper.Set("drivers.provisioners.proxmoxlxc.defaultnode", "pxvh1")

	return viper.WriteConfig()
}
