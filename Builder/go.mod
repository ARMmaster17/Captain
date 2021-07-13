module github.com/ARMmaster17/Captain/builder

go 1.16

require (
	github.com/ARMmaster17/Captain/Shared v0.0.0-20210709094226-b2e738dbb7ae
	github.com/Telmate/proxmox-api-go v0.0.0-20210713162220-06eec78e453b
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/ARMmaster17/Captain/Shared => ../Shared
