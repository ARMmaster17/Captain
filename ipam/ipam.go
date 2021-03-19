package ipam

import (
	"errors"
	http2 "github.com/ARMmaster17/Captain/http"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"time"
)

const AuthTimeout = 60 * 60

type IPAM struct {
	Client *http2.HttpRestClient
	LastAuthTime int64
}

func NewIPAM() (*IPAM, error) {
	url := os.Getenv("IPAM_API_URL") + "api/" + os.Getenv("IPAM_APP_NAME") + "/"
	client, err := http2.NewHttpRestClient(url, "application/json", true)
	if err != nil {
		log.Println(err)
		return &IPAM{}, errors.New("unable to build HttpRestClient")
	}

	return &IPAM{
		Client: client,
	}, nil
}

func (i *IPAM) Authenticate() error {
	currentTime := time.Now().Unix()
	if currentTime - i.LastAuthTime < AuthTimeout {
		return nil
	}
	username := os.Getenv("IPAM_USERNAME")
	password := os.Getenv("IPAM_PASSWORD")
	body, err := i.Client.PostWithBasicAuth("user/", username, password, nil)
	if err != nil {
		log.Println(err)
		return errors.New("unable to authenticate with IPAM API")
	}
	token := gjson.Get(body, "data.token").String()
	if token == "" {
		// TODO: Read error fields
		return errors.New("authentication with IPAM API did not return a token header")
	}
	i.Client.Headers["token"] = token
	i.LastAuthTime = currentTime
	return nil
}
