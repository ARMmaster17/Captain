package main

import (
	"errors"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func createLxcContainer(config MachineConfig) (string, error) {
	vmid, err := getNextVmid()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to obtain next available VMID")
	}
	config.VMID = vmid
	encoder := schema.NewEncoder()
	form := url.Values{}
	parseErr := encoder.Encode(config, form)
	if parseErr != nil {
		log.Println(err)
		return "", errors.New("unable to parse new LXC machine data")
	}
	values := form.Encode()
	url := os.Getenv("PROXMOX_API_URL")
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	client, request, err := prepareRequest("POST", url + "nodes/" + node + "/lxc", strings.NewReader(values))
	if err != nil {
		log.Println(err)
		return "", errors.New("failed prepare HTTP request for new LXC container")
	}
	request.Form = form
	response, err := client.Do(&request)
	defer response.Body.Close()
	return vmid, nil
}

func destroyLxcContainer(vmid string) error {
	url := os.Getenv("PROXMOX_API_URL")
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	client, request, err := prepareRequest("DELETE", url + "nodes/" + node + "/lxc/" + vmid, nil)
	if err != nil {
		log.Println(err)
		return errors.New("failed to prepare request to delete node")
	}
	response, err := client.Do(&request)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(err)
		return errors.New("request to delete node failed")
	}
	return nil
}

func stopLxcContainer(vmid string) error {
	url := os.Getenv("PROXMOX_API_URL")
	node := os.Getenv("PROXMOX_DEFAULT_NODE")
	client, request, err := prepareRequest("POST", url + "nodes/" + node + "/lxc/" + vmid + "/status/shutdown", nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to preprare LXC stop request")
	}
	response, err := client.Do(&request)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Println(err)
		return errors.New("request to stop node failed")
	}
	return nil
}