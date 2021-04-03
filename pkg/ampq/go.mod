module github.com/ARMmaster17/Captain/shared/ampq

go 1.16

require (
	github.com/streadway/amqp v1.0.0
	github.com/ARMmaster17/Captain/shared/captain v0.0.0-00010101000000-000000000000
)

replace (
	github.com/ARMmaster17/Captain/shared/captain => "../captain"
)