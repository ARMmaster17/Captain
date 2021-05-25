package IPAM

import (
	"gorm.io/gorm"
	"net"
)

type ReservedAddress struct {
	gorm.Model

}

// nextIP Shamelessly borrowed from https://stackoverflow.com/a/49057611. Increments an IP address by the
// specified amount.
func nextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}
