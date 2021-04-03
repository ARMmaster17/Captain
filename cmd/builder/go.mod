module github.com/ARMmaster17/Captain/cmd/builder

go 1.16

require (
	github.com/gorilla/schema v1.2.0
	github.com/streadway/amqp v1.0.0
	github.com/tidwall/gjson v1.6.8
	github.com/tidwall/pretty v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

//replace (
//	github.com/ARMmaster17/Captain/shared/ampq => "../shared/ampq"
//	github.com/ARMmaster17/Captain/shared/ipam => "../shared/ipam"
//	github.com/ARMmaster17/Captain/shared/proxmox => "../shared/proxmox"
//	github.com/ARMmaster17/Captain/shared/locking => "../shared/locking"
//	github.com/ARMmaster17/Captain/shared/prep => "../shared/prep"
//)
