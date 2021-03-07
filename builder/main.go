package main

import (
	"encoding/json"
	"errors"
	"github.com/ARMmaster17/Captain/shared/ampq"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func startListening() error {
	rabbitURI := os.Getenv("RABBITMQ_URI")
	conn, err := ampq.GetQueueConnection(rabbitURI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to connect to RabbitMQ")
	}

	err = conn.StartConsumer("builder-queue", "builder", handleMessage, 1)
	if err != nil {
		log.Println(err)
		return errors.New("unable to connect to builder queue")
	}
	return nil
}

func handleMessage(d amqp.Delivery) bool {
	if d.Body == nil {
		log.Println("empty message received")
		return false
	}
	plane := Plane{}
	err := json.Unmarshal(d.Body, &plane)
	if err != nil {
		log.Println("unable to parse message JSON")
		return false
	}
	_, err = makePlane(plane)
	if err != nil {
		log.Println("unable to execute plane build request")
		return false
	}
	return true
}

func main() {
	err := startListening()
	if err != nil {
		log.Println(err)
	}
	forever := make(chan bool)
	<-forever
}
