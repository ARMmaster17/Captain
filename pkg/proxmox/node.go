package proxmox

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
)

type Node struct {
	Name string `json:"node"`
	Status string	`json:"status"`
	// TODO: find out if CPU is an int or float
	RAMUsed int	`json:"mem"`
	RAMTotal	int `json:"maxmem"`
}

// Gets a list of nodes in the cluster with their current stats.
func (p *Proxmox) getNodes() ([]Node, error) {
	body, err := p.doGet("nodes")
	if err != nil {
		log.Println(err)
		return nil, errors.New("unable to get list of nodes")
	}
	nodeListRaw := gjson.Get(body, "data").Array()
	nodeList := make([]Node, len(nodeListRaw))
	for i := 0; i < len(nodeListRaw); i++ {
		err = json.Unmarshal([]byte(nodeListRaw[i].String()), &nodeList[i])
		if err != nil {
			log.Println(err)
			return nodeList, errors.New("unable to parse list of Proxmox nodes")
		}
	}
	return nodeList, nil
}

// Returns the node with the least % of RAM used.
func (p *Proxmox) GetLowestRAMUtilizationNode() (Node, error) {
	nodeList, err := p.getNodes()
	if err != nil {
		log.Println(err)
		return Node{}, errors.New("unable to get list of nodes")
	}
	var lowestNodeIndex int = -1
	var lowestRamPct float64 = 1
	for i := 1; i < len(nodeList); i++ {
		if nodeList[i].Status != "online" {
			continue
		}
		nodeRamPct := float64(nodeList[i].RAMUsed) / float64(nodeList[i].RAMTotal)
		if nodeRamPct < lowestRamPct {
			lowestRamPct = nodeRamPct
			lowestNodeIndex = i
		}
	}
	if lowestNodeIndex == -1 {
		return Node{}, errors.New("no eligible nodes found")
	}
	return nodeList[lowestNodeIndex], nil
}

func (p *Proxmox) FindNodeWithVMID(vmid string) (Node, error) {
	nodeList, err := p.getNodes()
	if err != nil {
		log.Println(err)
		return Node{}, errors.New("unable to obtain list of nodes")
	}
	for _, node := range nodeList {
		containerList, err := node.GetLXCContainers(p)
		if err != nil {
			log.Println(err)
			return Node{}, errors.New("unable to obtain list of containers on node " + node.Name)
		}
		for _, container := range containerList {
			if container.VMID == vmid {
				return node, nil
			}
		}
	}
	return Node{}, errors.New("VMID " + vmid + " not found")
}

func (n *Node) GetLXCContainers(p *Proxmox) ([]LXC, error) {
	url := fmt.Sprintf("nodes/%s/lxc", n.Name)
	body, err := p.doGet(url)
	if err != nil {
		log.Println(err)
		return []LXC{}, errors.New("unable to query Proxmox API for list of containers")
	}
	containers := gjson.Get(body, "data").Array()
	containerList := make([]LXC, len(containers))
	for i := 0; i < len(containers); i++ {
		err = json.Unmarshal([]byte(containers[i].String()), &containerList[i])
	}
	return containerList, nil
}
