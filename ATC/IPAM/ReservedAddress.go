package IPAM

import (
	"gorm.io/gorm"
	"net"
)

// ReservedAddress A database object that represents a reservation of an IP address by a plane on any provisioning
// driver.
type ReservedAddress struct {
	gorm.Model
	ReservedBlockID uint
	Address net.IP
}
