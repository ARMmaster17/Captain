package metadata

var (
	// Version is the semantic version of this application. Set at compile time with GCC flags.
	Version string
)

// GetCaptainVersion returns the current version of the application as set by the compiler.
func GetCaptainVersion() string {
	return Version
}
