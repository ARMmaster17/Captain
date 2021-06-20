package providers

import (
	"fmt"
	"github.com/ARMmaster17/Captain/Shared"
)

// GenericPlane is a generic description of a Plane with information that is relevant for provider drivers.
type GenericPlane struct {
	// The fully-qualified domain name of the plane on the public network.
	FQDN string
	// A string that is used by Captain to determine what driver manages what Plane objects. Also contains a unique
	// key that is used by provisioning drivers to map a Plane to the internal structure used by the hypervisor.
	CUID string
	// Number of CPU cores that should be assigned to the Plane.
	Cores int
	// Amount of RAM in megabytes that should be assigned to the Plane.
	RAM int
	// Size of disk in gigabytes that should be assigned to the Plane.
	Disk int
	// Identifier used for direct connections.
	NetID string

	//NetworkInterfaces []GenericNetworkInterface
}

//type GenericNetworkInterface struct {
//	Name string
//	IP string
//	Subnet string
//	MTU int
//	Gateway string
//}

func ConvertToGenericPlane(plane Shared.Plane, formation Shared.Formation) GenericPlane {
	return GenericPlane{
		FQDN: fmt.Sprintf("%s%d.%s", formation.BaseName, plane.ID, formation.Domain),
		CUID: plane.DriverIdentifier,
		Cores: formation.CPU,
		RAM: formation.RAM,
		Disk: formation.Disk,
		NetID: plane.NetID,
	}
}
