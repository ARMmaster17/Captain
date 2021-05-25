package IPAM

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"net"
)

type ReservedBlock struct {
	gorm.Model
	BlockName string
	Subnet    net.IPNet
	Addresses []ReservedAddress
}

func (ipam *IPAM) getAllReservedBlocks() ([]ReservedBlock, error) {
	var allBlocks []ReservedBlock
	result := ipam.db.Find(&allBlocks)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get IP blocks from database: %w", result.Error)
	}
	return allBlocks, nil
}

func (r *ReservedBlock) hasAvailableAddress(db *gorm.DB) (bool, error) {
	var addressCountInBlock int64
	result := db.Model(&ReservedAddress{}).Where("reservedblock_id = ?", r.ID).Count(&addressCountInBlock)
	if result.Error != nil {
		return false, fmt.Errorf("unable to get list of addresses in block: %w", result.Error)
	}
	return !subnetIsFull(addressCountInBlock, r.Subnet.Mask), nil
}

func (r *ReservedBlock) reserveAddress(db *gorm.DB) (net.IP, error) {
	ip, err := r.getNextAddress(db)
}

func (r *ReservedBlock) getNextAddress(db *gorm.DB) (net.IP, error) {
	for i := 0; int64(i) < getSubnetAddressSize(r.Subnet.Mask); i++ {
		ip := nextIP(r.Subnet.IP, uint(i))
		inUse, err := r.addressIsInUse(ip)
		if err != nil {
			return net.IP{}, fmt.Errorf("unable to check if next address is in use: %w", err)
		} else if !inUse {
			return ip, nil
		}
	}
	// This shouldn't happen because we already validate that the block has at least one open address.
	return net.IP{}, fmt.Errorf("the IP block metadata and address table have become desyncronized")
}

func (r *ReservedBlock) addressIsInUse(ip net.IP) (bool, error) {

}

func subnetIsFull(existingAddresses int64, subnetCIDR net.IPMask) bool {
	return existingAddresses >= getSubnetAddressSize(subnetCIDR)
}

func getSubnetAddressSize(subnetCIDR net.IPMask) int64 {
	leadingOnes, bits := subnetCIDR.Size()
	return int64(math.Pow(2, float64(bits - leadingOnes)) - 2)
}
