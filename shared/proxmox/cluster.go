package proxmox

import (
	"errors"
	"github.com/tidwall/gjson"
	"log"
)

func (p *Proxmox) getNextVmid() (string, error) {
	body, err := p.doGet("cluster/nextid")
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to query API for next free VMID")
	}
	vmid := gjson.Get(body, "data").String()
	if vmid == "" {
		// TODO: read error fields
		return "", errors.New("Proxmox API did not return a free VMID")
	}
	return vmid, nil
}