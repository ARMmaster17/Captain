package IPAM

import (
	"fmt"
	"gorm.io/gorm"
	"net"
)

// ReservedBlock A database object that represents a block of IP addresses represented by a base address and a CIDR
// block identifier. ReservedBlocks act as address pools, and are filled in the order specified in config.yaml.
type ReservedBlock struct {
	gorm.Model
	BlockName string
	Subnet    net.IPNet
	Addresses []ReservedAddress
}

// getAllReservedBlocks Returns all active ReservedBlocks in the database.
func (ipam *IPAM) getAllReservedBlocks() ([]ReservedBlock, error) {
	var allBlocks []ReservedBlock
	result := ipam.db.Find(&allBlocks)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get IP blocks from database: %w", result.Error)
	}
	return allBlocks, nil
}

// hasAvailableAddress Checks if the given ReservedBlock has at least one free address in the pool.
func (r *ReservedBlock) hasAvailableAddress(db *gorm.DB) (bool, error) {
	var addressCountInBlock int64
	result := db.Model(&ReservedAddress{}).Where("reservedblock_id = ?", r.ID).Count(&addressCountInBlock)
	if result.Error != nil {
		return false, fmt.Errorf("unable to get list of addresses in block: %w", result.Error)
	}
	return !subnetIsFull(addressCountInBlock, r.Subnet.Mask), nil
}

// reserveAddress Finds the first available address, and then reserves it into the database. This function is not
// thread-safe, and assumes that the database is locked by the calling function.
func (r *ReservedBlock) reserveAddress(db *gorm.DB) (net.IP, error) {
	ip, err := r.getNextAddress(db)
	if err != nil {
		return nil, fmt.Errorf("unable get next available address: %w", err)
	}
	newAddress := ReservedAddress{
		Address: ip,
		ReservedBlockID: r.ID,
	}
	result := db.Save(&newAddress)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to save new address: %w", err)
	} else {
		return ip, nil
	}
}

// getNextAddress Finds the first available address in a ReservedBlock pool. This function only retrieves a free address,
// and does not reserve it.
func (r *ReservedBlock) getNextAddress(db *gorm.DB) (net.IP, error) {
	var usedAddresses []ReservedAddress
	result := db.Model(&ReservedAddress{}).Where("reservedblock_id = ?", r.ID).First(&usedAddresses)
	if result.Error != nil {
		return nil, fmt.Errorf("unable to get addresses in block: %w", result.Error)
	}
	for i := 0; int64(i) < getSubnetAddressSize(r.Subnet.Mask); i++ {
		ip := nextIP(r.Subnet.IP, uint(i))
		if !r.addressIsInUse(ip, usedAddresses) {
			return ip, nil
		}
	}
	// This shouldn't happen because we already validate that the block has at least one open address.
	return net.IP{}, fmt.Errorf("the IP block metadata and address table have become desyncronized")
}

// addressIsInUse Checks if the given address is currently reserved the ReservedBlock.
func (r *ReservedBlock) addressIsInUse(ip net.IP, usedAddresses []ReservedAddress) bool {
	for _, address := range usedAddresses {
		if ip.String() == address.Address.String() {
			return true
		}
	}
	return false
}
