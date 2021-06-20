package Shared

import "gorm.io/gorm"

// Plane Represents a running instance on any provider, such as an LXC container, a VM, or a Docker/Kubernetes container
// instance. Planes are usually not modified directly by the API, as they are automatically managed by the built-in
// health checks for each formation. In the event of a configuration change at the formation level, planes are usually
// destroyed and recreated in favor of modifying configuration through a complex provider adapter.
type Plane struct {
	gorm.Model
	// Used by provisioning drivers to identify this plane within the context of that hypervisor platform.
	DriverIdentifier string
	// ID of the formation that this plane belongs to. Contains most of the config fields for this plane.
	FormationID int64
	// The network identifier assigned by the IPAM module. Used for preflight provisioning.
	NetID string
}
