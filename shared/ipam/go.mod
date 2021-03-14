module github.com/ARMmaster17/Captain/shared/ipam

go 1.16

require (
	github.com/gorilla/schema v1.2.0
	github.com/streadway/amqp v1.0.0
	github.com/tidwall/gjson v1.6.8
	github.com/tidwall/pretty v1.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	github.com/ARMmaster17/Captain/shared/http v0.0.0-00010101000000-000000000000
)

replace (
	github.com/ARMmaster17/Captain/shared/http => "../http"
)