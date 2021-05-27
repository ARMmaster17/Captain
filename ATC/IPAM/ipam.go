package IPAM

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net"
	"sync"
)

// IPAM is a thread-safe simple implementation of an IP address management system that reserves and releases
// addresses for Planes based on allowed blocks in config.yaml.
type IPAM struct {
	db *gorm.DB
	mutex *sync.Mutex
}

// Initialize performs any needed migrations on the database and creates a mutex object for safe cross-thread
// transactions.
func (ipam *IPAM) Initialize(db *gorm.DB) error {
	ipam.db = db
	err := ipam.performMigrations()
	if err != nil {
		return fmt.Errorf("unable to perform IPAM database migrations: %w", err)
	}
	ipam.mutex = &sync.Mutex{}
	err = ipam.syncDBBlocksWithConfig(viper.GetStringSlice("defaults.network.blocks"))
	if err != nil {
		return fmt.Errorf("unable to sync database with new IP block config: %w", err)
	}
	return nil
}

// GetNewAddress Finds, reserves, and returns an IP address from the first available pool.
func (ipam *IPAM) GetNewAddress() (net.IP, error) {
	ipam.mutex.Lock()
	freeBlock, err := ipam.findFreeBlock()
	if err != nil {
		ipam.mutex.Unlock()
		return net.IP{}, fmt.Errorf("unable to find a free IP block: %w", err)
	}
	newAddress, err := freeBlock.reserveAddress(ipam.db)
	ipam.mutex.Unlock()
	return newAddress, err
}

// ReleaseAddress Finds what pool an address belongs to, and then deletes that address from the pool.
func (ipam *IPAM) ReleaseAddress(addr net.IP) error {
	ipam.mutex.Lock()
	// TODO: Not implemented.
	ipam.mutex.Unlock()
	return nil
}

// syncDBBlocksWithConfig Loads values from config.yaml and performs a reconciliation with the available ReservedBlocks
// in the database.
func (ipam *IPAM) syncDBBlocksWithConfig(blocks []string) error {
	netBlocks := parseSubnetBlocks(blocks)
	existingReservedBlocks, err := ipam.getAllReservedBlocks()
	if err != nil {
		return fmt.Errorf("unable to retrieve list of reserved IP blocks: %w", err)
	}
	for _, netBlock := range netBlocks {
		var err error
		existingReservedBlocks, err = ipam.addNetblockIfNotExists(netBlock, existingReservedBlocks)
		if err != nil {
			return fmt.Errorf("unable to merge IP block configuration with database: %w", err)
		}
	}
	// TODO: Remove blocks that are no longer in the config. This will require deleting planes, and not allowing
	// new planes to be built until the IP block table is fully synced.
	return nil
}

// addNetblockIfNotExists Searches a given list of ReservedBlocks. If the given CIDR block does not exist, it will
// be created and saved into the database.
func (ipam *IPAM) addNetblockIfNotExists(netBlock net.IPNet, existingReservedBlocks []ReservedBlock) ([]ReservedBlock, error) {
	for _, block := range existingReservedBlocks {
		if block.BlockName == netBlock.String() {
			return existingReservedBlocks, nil
		}
	}
	// Block does not exist, create it.
	newBlock := ReservedBlock{
		BlockName: netBlock.String(),
		Subnet:    netBlock,
	}
	result := ipam.db.Create(&newBlock)
	if result.Error != nil {
		return existingReservedBlocks, fmt.Errorf("unable to create new reserved IP block: %w", result.Error)
	}
	return append(existingReservedBlocks, newBlock), nil
}

// findFreeBlock Searches through all CIDR blocks given in config.yaml, and returns the first block with at least one
// free address.
func (ipam *IPAM) findFreeBlock() (ReservedBlock, error) {
	reservedBlocks, err := ipam.getAllReservedBlocks()
	if err != nil {
		return ReservedBlock{}, fmt.Errorf("unable to get IP block list: %w", err)
	}
	for _, block := range reservedBlocks {
		isAvailable, err := block.hasAvailableAddress(ipam.db)
		if err != nil {
			return ReservedBlock{}, fmt.Errorf("unable to check if block has free addresses: %w", err)
		} else if isAvailable {
			return block, nil
		}
	}
	return ReservedBlock{}, fmt.Errorf("all reserved IP blocks are full")
}