module github.com/ARMmaster17/Captain/Builder

go 1.16

require (
	github.com/ARMmaster17/Captain/Shared v0.0.0-00010101000000-000000000000
	github.com/Selvatico/go-mocket v1.0.7 // indirect
	github.com/Telmate/proxmox-api-go 518c081d3063
	github.com/rs/zerolog v1.25.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/ARMmaster17/Captain/Shared => ../Shared
