package Preflight

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/go-ping/ping"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"time"
)

func PreflightSingleInstance(connectionURI string, playbookPath string) error {
	err := waitForHostToComeOnline(connectionURI)
	if err != nil {
		return fmt.Errorf("plane %s is not online, cannot run preflight:\n%w", connectionURI, err)
	}
	hostFilePath, err := generateHostsFile(connectionURI)
	if err != nil {
		return fmt.Errorf("unable to generate temporary hosts file:\n%w", err)
	}
	defer cleanupHostsFile(hostFilePath)
	err = ansibleProvisionHost(hostFilePath, viper.GetString("config.preflight.ansible.privatekeypath"), []string{processPlaybookPath(playbookPath)})
	if err != nil {
		return fmt.Errorf("unable to run preflight playbook on %s:\n%w", connectionURI, err)
	}
	return nil
}

func generateHostsFile(connectionURI string) (string, error) {
	tmpHostsFile, err := ioutil.TempFile(os.TempDir(), "hosts-")
	if err != nil {
		return "", fmt.Errorf("unable to create temporary file:\n%w", err)
	}
	defer tmpHostsFile.Close()
	fileBody := []byte(connectionURI)
	if _, err = tmpHostsFile.Write(fileBody); err != nil {
		return "", fmt.Errorf("unable to write host data to temporary file:\n%w", err)
	}
	return tmpHostsFile.Name(), nil
}

func cleanupHostsFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Warn().Msgf("unable to delete temporary file at %s:\n%w", path, err)
	}
}

func ansibleProvisionHost(hostFilePath string, privatekeyPath string, playbookPaths []string) error {
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection:    viper.GetString("config.preflight.ansible.connectiontype"),
		PrivateKey:    privatekeyPath,
		Timeout:       5,
		User:          viper.GetString("config.preflight.ansible.connectionuser"),
	}
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:         hostFilePath,
	}
	playbookObj := &playbook.AnsiblePlaybookCmd{
		Playbooks:                  playbookPaths,
		Options:                    ansiblePlaybookOptions,
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
	}
	err := playbookObj.Run(context.Background())
	if err != nil {
		return fmt.Errorf("preflight prep failed:\n%w", err)
	}
	return nil
}

func waitForHostToComeOnline(connectionURI string) error {
	var err error
	for i := 0; i < viper.GetInt("config.preflight.ping.retries"); i++ {
		err = pingHost(connectionURI)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(viper.GetInt64("config.preflight.ping.timeout")) * time.Second)
	}
	return fmt.Errorf("unable to check if host is online:\n%w", err)
}

func pingHost(connectionURI string) error {
	pinger, err := ping.NewPinger(connectionURI)
	if err != nil {
		return fmt.Errorf("unable to build ping request:\n%w", err)
	}
	pinger.Count = 1
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		return fmt.Errorf("unable to ping host:\n%w", err)
	}
	return nil
}

func processPlaybookPath(path string) string {
	return fmt.Sprintf("/etc/captain/atc/playbooks/%s", path)
}