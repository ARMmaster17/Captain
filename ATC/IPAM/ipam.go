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
	ipam.syncDBBlocksWithConfig(viper.GetStringSlice("defaults.network.blocks"))
	return nil
}

func (ipam *IPAM) GetNextAddress() (net.IPAddr, error) {
	ipam.mutex.Lock()

	ipam.mutex.Unlock()
}

func (ipam *IPAM) ReleaseAddress(addr net.IPAddr) error {
	ipam.mutex.Lock()

	ipam.mutex.Unlock()
}

func (ipam *IPAM) syncDBBlocksWithConfig(blocks []string) error {
	netBlocks := parseSubnetBlocks(blocks)
	var existingReservedBlocks []ReservedBlock
	result := ipam.db.Find(&existingReservedBlocks)
	if result.Error != nil {
		return fmt.Errorf("unable to retrieve list of reserved IP blocks: %w", result.Error)
	}

}
