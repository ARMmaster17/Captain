package proxmox

import (
	"errors"
	"github.com/tidwall/gjson"
	"log"
)

// Gets the next available VMID from the cluster. If this VMID will be used
// for something, note that this operation is not thread-safe, and requires some
// kind of external locking mechanism.
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