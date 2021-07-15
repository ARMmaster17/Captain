module github.com/ARMmaster17/Captain/Builder

go 1.16

require (
	github.com/ARMmaster17/Captain/Shared v0.0.0-00010101000000-000000000000
	github.com/Telmate/proxmox-api-go v0.0.0-20210713162220-06eec78e453b
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/ARMmaster17/Captain/Shared => ../Shared
