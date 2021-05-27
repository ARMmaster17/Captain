package IPAM

import (
	"github.com/rs/zerolog/log"
	"math"
	"net"
)

// parseSubnetBlocks parses a slice of string representations of IP address blocks + CIDRs into a slice of IPNet
// objects. Invalid strings are ignored, and a debug statement is logged.
func parseSubnetBlocks(blocks []string) []net.IPNet {
	var outBlocks []net.IPNet
	for _, block := range blocks {
		_, net, err := net.ParseCIDR(block)
		if err != nil {
			log.Debug().Msgf("unrecognized IP block: %s", block)
		} else {
			outBlocks = append(outBlocks, *net)
		}
	}
	return outBlocks
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

// subnetIsFull Takes a CIDR block and the number of addresses in use. Returns true if the subnet is full. This
// function does not take into account broadcast addresses or other special address blocks.
func subnetIsFull(existingAddresses int64, subnetCIDR net.IPMask) bool {
	return existingAddresses >= getSubnetAddressSize(subnetCIDR)
}

// getSubnetAddressSize Returns the maximum number of addresses that can exist in the given CIDR block. Does not
// perform special operations for P2P blocks, or take into account broadcast addresses as it is assumed that addresses
// given to IPAM do not represent the actual operating subnet.
func getSubnetAddressSize(subnetCIDR net.IPMask) int64 {
	leadingOnes, bits := subnetCIDR.Size()
	return int64(math.Pow(2, float64(bits - leadingOnes))/* - 2*/)
}
