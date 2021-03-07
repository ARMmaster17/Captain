package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getProxmoxAuthData(client *http.Client) (string, string, error) {
	url := os.Getenv("PROXMOX_API_URL")
	username := os.Getenv("PROXMOX_USERNAME")
	password := os.Getenv("PROXMOX_PASSWORD")
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		log.Println(err)
		return "", "", errors.New("invalid Proxmox username/password format")
	}
	response, err := client.Post(url + "access/ticket", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
		return "", "", errors.New("unable to access Proxmox API")
	}
	defer response.Body.Close()
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Println(err)
		return "", "", errors.New("the Proxmox API returned invalid data")
	}

	csrf := gjson.Get(string(body), "data.CSRFPreventionToken").String()
	token := gjson.Get(string(body), "data.ticket").String()
	return csrf, token, nil
}

func getProxmoxHttpClient() (http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return *client, nil
}

func prepareRequest(method string, url string, body io.Reader) (http.Client, http.Request, error) {
	request, _ := http.NewRequest(method, url, body)

	client, err := getProxmoxHttpClient()
	if err != nil {
		log.Println(err)
		return client, *request, errors.New("unable to prepare HTTP request object")
	}
	csrf, token, err := getProxmoxAuthData(&client)
	if err != nil {
		log.Println(err)
		return client, *request, errors.New("unable to retrieve authentication data")
	}
	request.Header.Set("CSRFPreventionToken", csrf)
	request.Header.Set("Cookie", "PVEAuthCookie=" + token)
	request.Header.Set("ContentType", "application/x-www-url-formencoding")
	return client, *request, nil
}

func getNextVmid() (string, error) {
	url := os.Getenv("PROXMOX_API_URL")
	client, request, err := prepareRequest("GET", url + "cluster/nextid", nil)
	if err != nil {
		log.Println(err)
		return "", errors.New("error preparing HTTP request")
	}
	response, err := client.Do(&request)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to obtain next available VMID")
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New("the Proxmox API returned invalid data")
	}

	vmid := gjson.Get(string(body), "data").String()
	return vmid, nil
}

type MachineConfig struct {
	VMID			string	`schema:"vmid"`
	OsTemplate		string	`schema:"ostemplate"`
	Net0			string	`schema:"net0"`
	Hostname		string	`schema:"hostname"`
	Cores			int		`schema:"cores"`
	Memory			int		`schema:"memory"`
	Swap			int		`schema:"swap"`
	Nameservers		string	`schema:"nameserver"`
	Storage			string	`schema:"storage"`
	RootFS			string	`schema:"rootfs"`
	OnBoot			int		`schema:"onboot"`
	Unprivileged	int		`schema:"unprivileged"`
	Start			int		`schema:"start"`
}

type NetworkConfig struct {
	Name		string	`schema:"name"`
	Bridge		string	`schema:"bridge"`
	IP			string	`schema:"ip"`
	Gateway		string	`schema:"gw"`
	Firewall	int		`schema:"firewall"`
	MTU			int		`schema:"mtu"`
}
