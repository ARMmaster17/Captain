package proxmox

type Ticket struct {
	Node string
	UPID string
}

func (t *Ticket) IsRunning(p *Proxmox) (bool, error) {

}
