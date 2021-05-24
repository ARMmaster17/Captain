package IPAM

import "fmt"

func (ipam *IPAM) performMigrations() error {
	err := ipam.db.AutoMigrate(&ReservedAddress{})
	if err != nil {
		return fmt.Errorf("unable to initialize IPAM database schema for ReservedAddress: %w", err)
	}
	err = ipam.db.AutoMigrate(&ReservedBlock{})
	if err != nil {
		return fmt.Errorf("unable to initialize IPAM database schema for ReservedBlock: %w", err)
	}
	return nil
}
