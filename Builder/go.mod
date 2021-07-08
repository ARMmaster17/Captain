module github.com/ARMmaster17/Captain/builder

go 1.16

require (
	github.com/ARMmaster17/Captain/Shared v0.0.0-20210708232013-b74a4aa813cd
	github.com/Telmate/proxmox-api-go v0.0.0-20210708200918-d27e0fa5a4a4
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/ARMmaster17/Captain/Shared => ../Shared
