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
	IP              string
}

// GetIP Converts the underlying DB representation to a net.IP object.
func (r *ReservedAddress) GetIP() net.IP {
	return net.ParseIP(r.IP)
}

// SetIP Converts a net.IP object into the underlying DB representation of the field.
func (r *ReservedAddress) SetIP(ip net.IP) {
	r.IP = ip.String()
}
