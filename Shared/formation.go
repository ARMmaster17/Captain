package Shared

import "gorm.io/gorm"

// Formation is the lowest-level object that is directly addressable by the user (within the context of Captain). A
// formation manages a group of planes that are scaled up and down automatically. A formation is a logical representation
// of an internal service for an application. For example, all web servers that serve the same web app would be
// considered part of the same formation. All planes in a formation will be exactly the same except for the FQDN.
type Formation struct {
	gorm.Model
	// Name of the service. Used only in user-facing queries, and is not used internally. Should be as unique as
	// possible for easy identification.
	Name string `validate:"required,min=1"`
	// Number of CPU cores to assign to each plane. The actual implementation of this varies depending on which
	// provider adapter is used.
	CPU int `validate:"required,gte=1,lte=8192"`
	// Amount of RAM in megabytes to assign to each plane.
	RAM int `validate:"required,gte=1,lte=307200"`
	// Size of disk in gigabytes to assign to each plane. It is important that this disk is big enough to store the
	// container OS in addition to application data.
	Disk int `validate:"required,gte=1"`
	// URL-safe name for each plane in formation. Should be unique within the flight. Will be used in the FQDN of each
	// plane that is provisioned within the formation. For example: formation1.example.com.
	BaseName	string `validate:"alphanum,min=1,max=256"`
	// Domain name that forms the ending of the FQDN for each plane. In the future this will be moved to be the same
	// airspace-wide or stack-wide.
	Domain		string `validate:"required,fqdn,min=1"`
	// Desired number of planes that should be operational at any given moment. At each health check interval,
	// remediations will be made to adjust the number of healthy planes in service until it equals this number.
	TargetCount int     `validate:"gte=0"`
	FlightID    int
}
