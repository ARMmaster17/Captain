package providers

// ProviderDriver is an abstract class that defines the basic capabilities that Captain expects of all provider drivers.
type ProviderDriver interface {
	// Connect will connect to the underlying provider, performing any
	// authentication or configuration necessary to service requests.
	Connect() error
	// BuildPlane takes a generic plane plan, and handles applying it to a new
	// instance in the underlying machine provider. Should return a unique
	// identifying string WITHOUT the driver's CUID added as a prefix.
	BuildPlane(p *GenericPlane) (string, error)
	DestroyPlane(cuid string, p *GenericPlane) error
	// GetCUIDPrefix returns the driver's unique identifier present in a plane's
	// CUID value in the state database. This is used to identify which provider
	// a plane belongs to, and stores a UID that the driver can use to reference the
	// plane internally.
	GetCUIDPrefix() string
	// GetYAMLTag returns the driver's unique YAML tag used in driver configurations.
	// These configured values will be passed with each driver call, so there is no
	// need to store them in the driver.
	GetYAMLTag() string
}
