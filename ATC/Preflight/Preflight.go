package Preflight

import (
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/go-ping/ping"
	"github.com/spf13/viper"
	"time"
)

func PreflightSingleInstance(connectionURI string, playbookPath string) error {
	err := waitForHostToComeOnline(connectionURI)
	if err != nil {
		return fmt.Errorf("plane %s is not online, cannot run preflight:\n%w", connectionURI, err)
	}
	err = ansibleProvisionHost(connectionURI, viper.GetString("config.preflight.ansible.privatekeypath"), []string{processPlaybookPath(playbookPath)})
	if err != nil {
		return fmt.Errorf("unable to run preflight playbook on %s:\n%w", connectionURI, err)
	}
	return nil
}

func ansibleProvisionHost(connectionURI string, privatekeyPath string, playbookPaths []string) error {
	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection:    viper.GetString("config.preflight.ansible.connectiontype"),
		PrivateKey:    privatekeyPath,
		Timeout:       5,
		User:          viper.GetString("config.preflight.ansible.connectionuser"),
	}
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory:         connectionURI,
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