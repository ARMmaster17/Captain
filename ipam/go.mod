module github.com/ARMmaster17/Captain/shared/ipam

go 1.16

require (
	github.com/ARMmaster17/Captain/shared/http v0.0.0-00010101000000-000000000000
	github.com/tidwall/gjson v1.6.8
	github.com/tidwall/pretty v1.1.0 // indirect
)

replace (
	github.com/ARMmaster17/Captain/shared/http => ./../http
)
