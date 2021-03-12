package ipam

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

type AddressRequest struct {
	SubnetID string `schema:"subnetId"`
	IP		 string	`schema:"ip"`
	Hostname	string	`schema:"hostname"`
}

type IPAddress string

func (i *IPAM) IPCreateFromFirstFree(hostname Hostname) (IPAddress, error) {
	i.Authenticate()
	address, err := i.IPFirstFree()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to obtain free IP address in block to register for host " + string(hostname))
	}
	err = i.IPCreate(address, hostname)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to register IP " + string(address))
	}
	return address, nil
}

func (i *IPAM) IPCreate(address IPAddress, hostname Hostname) error {
	i.Authenticate()
	url := fmt.Sprintf("api/%s/addresses/", os.Getenv("IPAM_APP_NAME"))
	var addressRequest = AddressRequest{
		SubnetID: "12",
		IP:       string(address),
		Hostname: string(hostname),
	}
	_, err := i.Client.PostUrlEncode(url, addressRequest)
	if err != nil {
		log.Println(err)
		return errors.New(fmt.Sprintf("unable to register IP address %s to host %s", address, hostname))
	}
	return nil
}

func (i *IPAM) IPFirstFree() (IPAddress, error) {
	i.Authenticate()
	url := fmt.Sprintf("api/%s/subnets/%d/first_free/", os.Getenv("IPAM_APP_NAME"), 12)
	body, err := i.Client.Get(url)
	if err != nil {
		return "", errors.New("unable to query first free IP address in subnet")
	}
	ip := gjson.Get(body, "data").String()
	if ip == "" {
		return "", errors.New("the IPAM API did not return a free IP address in response")
	}
	return IPAddress(ip), nil
}

func (a *IPAddress) Delete(i *IPAM) error {
	i.Authenticate()
	url := fmt.Sprintf("api/%s/addresses/%s/%d/", os.Getenv("IPAM_APP_NAME"), a, 12)
	_, err := i.Client.Delete(url)
	if err != nil {
		return errors.New("unable to delete address " + string(*a))
	}
	return nil
}