package providers

import "net"

type GenericPlane struct {
	FQDN string
	CUID string
	Cores int
	RAM int
	Disk int
	NetworkInterfaces []GenericNetworkInterface
}

type GenericNetworkInterface struct {
	Name string
	IP net.IP
	Subnet net.IPNet
	MTU int
	Gateway net.IP
}
