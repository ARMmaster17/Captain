package proxmox

import (
	"errors"
	http2 "github.com/ARMmaster17/Captain/pkg/http"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"time"
)

const AuthTimeout = 60 * 60 * 2

type Proxmox struct {
	Client *http2.HttpRestClient
	LastAuthTime int64
}

func NewProxmox() (*Proxmox, error) {
	client, err := http2.NewHttpRestClient(os.Getenv("PROXMOX_API_URL"), "application/json", true)
	if err != nil {
		log.Println(err)
		return &Proxmox{}, errors.New("unable to build HttpRestClient")
	}

	return &Proxmox{
		Client: client,
	}, nil
}

func (p *Proxmox) Authenticate() error {
	currentTime := time.Now().Unix()
	if currentTime - p.LastAuthTime < AuthTimeout {
		return nil
	}
	body, err := p.Client.PostJson("access/ticket", map[string]string{
		"username": os.Getenv("PROXMOX_USERNAME"),
		"password": os.Getenv("PROXMOX_PASSWORD"),
	})
	if err != nil {
		log.Println(err)
		return errors.New("unable to authenticate with Proxmox API")
	}
	csrf := gjson.Get(body, "data.CSRFPreventionToken").String()
	token := gjson.Get(body, "data.ticket").String()
	if csrf == "" || token == "" {
		// TODO: Parse data array for error message
		return errors.New("authentication with Proxmox API did not return a CSRF or token header")
	}
	p.LastAuthTime = currentTime
	p.Client.Headers["CSRFPreventionToken"] = csrf
	p.Client.Headers["Cookie"] = "PVEAuthCookie=" + token
	p.Client.ContentType = "application/x-www-url-formencoded"
	return nil
}
