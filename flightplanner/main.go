package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ampq"
	"github.com/ARMmaster17/Captain/shared/ipam"
	"github.com/ARMmaster17/Captain/shared/locking"
	"github.com/ARMmaster17/Captain/shared/prep"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func main() {
	err := startListening()
	if err != nil {
		log.Println(err)
	}
	forever := make(chan bool)
	<-forever
}
