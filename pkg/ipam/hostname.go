package ipam

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"os"
)

type Hostname string

func (h *Hostname) Delete(i *IPAM) error {
	i.Authenticate()
	ip, err := h.GetIP(i)
	if err != nil {
		return errors.New("unable to lookup IP of hostname record to delete")
	}
	err = ip.Delete(i)
	if err != nil {
		return errors.New("unable to delete ip record " + string(ip))
	}
	return nil
}

func (h *Hostname) GetIP(i *IPAM) (IPAddress, error) {
	i.Authenticate()
	url := fmt.Sprintf("api/%s/addresses/search_hostname/%s/", os.Getenv("IPAM_APP_NAME"), h)
	body, err := i.Client.Get(url)
	if err != nil {
		return IPAddress(""), errors.New("unable to query for IP address of hostname")
	}
	address := gjson.Get(string(body), "data.0.ip").String()
	if address == "" {
		return IPAddress(""), errors.New("no IP address found for host " + string(*h))
	}
	return IPAddress(address), nil
}