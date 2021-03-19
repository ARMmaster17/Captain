module github.com/ARMmaster17/Captain/shared/captain

go 1.16

require (
	github.com/ARMmaster17/Captain/shared/ipam v0.0.0-00010101000000-000000000000
	github.com/ARMmaster17/Captain/shared/proxmox v0.0.0-00010101000000-000000000000
)

replace (
	github.com/ARMmaster17/Captain/shared/ipam => "../ipam"
	github.com/ARMmaster17/Captain/shared/proxmox => "../proxmox"
)