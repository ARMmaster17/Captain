package ipam

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

func getIpamHttpClient() (http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return *client, nil
}

func getIpamAuthData(client *http.Client) (string, error) {
	url := os.Getenv("IPAM_API_URL")
	appName := os.Getenv("IPAM_APP_NAME")
	username := os.Getenv("IPAM_USERNAME")
	password := os.Getenv("IPAM_PASSWORD")

	request, err := http.NewRequest("POST", fmt.Sprintf("%sapi/%s/user/", url, appName), nil)
	if err != nil {
		log.Println(err)
		return "", errors.New("failed to build IPAM authentication request")
	}
	request.SetBasicAuth(username, password)
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to contact IPAM API")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New("the IPAM API returned invalid data")
	}

	token := gjson.Get(string(body), "data.token").String()

	return token, nil
}

func authneticateIpam() (*http.Client, string, error) {
	client, err := getIpamHttpClient()
	if err != nil {
		log.Println(err)
		return nil, "", errors.New("unable to create HTTP client")
	}
	token, err := getIpamAuthData(&client)
	if err != nil {
		log.Println(err)
		return nil, "", errors.New("unable to retrieve auth data from IPAM")
	}
	return &client, token, nil
}

func prepareIpamRequest(method string, url string, body io.Reader) (http.Client, http.Request, error) {
	ipamUrl := os.Getenv("IPAM_API_URL")
	client, token, err := authneticateIpam()
	if err != nil {
		log.Println(err)
		return http.Client{}, http.Request{}, errors.New("unable to authenticate with IPAM")
	}
	request, _ := http.NewRequest(method, fmt.Sprintf("%s%s", ipamUrl, url), body)
	request.Header.Set("token", token)
	return *client, *request, nil
}

func registerIPAddress(address string, hostname string) error {
	url := fmt.Sprintf("api/%s/addresses/", os.Getenv("IPAM_APP_NAME"))
	values := url2.Values{}
	values.Set("subnetId", "12")
	values.Set("ip", address)
	values.Set("hostname", hostname)
	client, request, err := prepareIpamRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		log.Println(err)
		return errors.New("unable to prepare HTTP request")
	}

	response, err := client.Do(&request)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return errors.New("the IPAM API returned invalid data")
	}
	return nil
}

func getFirstFreeIPAddress() (string, error) {
	url := fmt.Sprintf("api/%s/subnets/%d/first_free/", os.Getenv("IPAM_APP_NAME"), 12)
	client, request, err := prepareIpamRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to prepare HTTP request")
	}
	response, err := client.Do(&request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New("the Proxmox API returned invalid data")
	}
	ip := gjson.Get(string(body), "data").String()
	return ip, nil
}

func GetIPAddress(hostname string) (string, error) {
	ip, err := getFirstFreeIPAddress()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to obtain free IP address")
	}

	err = registerIPAddress(ip, hostname)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to register new IP address")
	}
	return ip, nil
}

func GetIPFromHostname(hostname string) (string, error) {
	url := fmt.Sprintf("api/%s/addresses/search_hostname/%s/", os.Getenv("IPAM_APP_NAME"), hostname)
	client, request, err := prepareIpamRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to prepare HTTP request")
	}

	response, err := client.Do(&request)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to send request to IPAM API")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New("the IPAM API returned invalid data")
	}
	address := gjson.Get(string(body), "data.0.ip").String()
	return address, nil
}

func deleteAddress(ip string) error {
	url := fmt.Sprintf("api/%s/addresses/%s/%d/", os.Getenv("IPAM_APP_NAME"), ip, 12)
	client, request, err := prepareIpamRequest("DELETE", url, nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to prepare HTTP request")
	}

	response, err := client.Do(&request)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return errors.New("the IPAM API returned invalid data")
	}
	return nil
}

func ReleaseIPAddress(hostname string) error {
	ip, err := GetIPFromHostname(hostname)
	if err != nil {
		log.Println(err)
		return errors.New("hostname lookup in IPAM failed")
	}
	err = deleteAddress(ip)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete address in IPAM")
	}
	return nil
}
