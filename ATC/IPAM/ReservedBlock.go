package IPAM

import (
	"gorm.io/gorm"
	"net"
)

type ReservedBlock struct {
	gorm.Model
	BlockName string
	Subnet    net.IPNet
	Addresses []ReservedAddress
}
