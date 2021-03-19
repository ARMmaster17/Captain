package proxmox

import (
	"errors"
	"github.com/tidwall/gjson"
	"log"
)

func (p *Proxmox) doGet(urlPath string) (string, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to authenticate with Proxmox API")
	}
	body, err := p.Client.Get(urlPath)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to GET against Proxmox API")
	}
	return parseErrors(body)
}

func (p *Proxmox) doPost(urlPath string, reqBody []byte) (string, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to authenticate with Proxmox API")
	}
	body, err := p.Client.Post(urlPath, reqBody)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to POST against Proxmox API")
	}
	return parseErrors(body)
}

func (p *Proxmox) doPostUrlEncode(urlPath string, reqBody interface{}) (string, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to authenticate with Proxmox API")
	}
	body, err := p.Client.PostUrlEncode(urlPath, reqBody)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to POST against Proxmox API")
	}
	return parseErrors(body)
}

func (p *Proxmox) doDelete(urlPath string) (string, error) {
	err := p.Authenticate()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to authenticate with Proxmox API")
	}
	body, err := p.Client.Delete(urlPath)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to DELETE against Proxmox API")
	}
	return parseErrors(body)
}

func parseErrors(body string) (string, error) {
	errorList := gjson.Get(body, "errors").String()
	if errorList == "" {
		return body, nil
	} else {
		log.Println(errorList)
		return "", errors.New("the Proxmox API returned errors with the previous request")
	}
}