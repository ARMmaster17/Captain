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
		return true
	}
	var message ampq.Message
	err := json.Unmarshal(d.Body, &message)
	if err != nil {
		log.Println(err)
		log.Println("unable to parse message JSON")
		return true
	}
	if message.Plane.Name == "" {
		log.Println("No plane name given in payload")
		return true
	}
	if message.Operation == "build" {
		return handleBuildMessage(message)
	} else if message.Operation == "destroy" {
		return handleDestroyMessage(message)
	} else {
		log.Println(fmt.Sprintf("Unknown operation: %s", message.Operation))
		return true
	}

}

func handleBuildMessage(message ampq.Message) bool {
	log.Printf("Recieved request to build plane %s in %d:%d:%d configuration\n", message.Plane.Name, message.Plane.CPU, message.Plane.RAM, message.Plane.Storage)
	mutex, context, err := locking.GetLockableMutex("pve-lxc-create")
	if err != nil {
		log.Println(err)
		log.Println("unable to create mutex")
		// TODO: might be smart to return false and let another container try
		return true
	}
	log.Println("Acquiring lock on pve-lxc-create...")
	err = locking.LockResource(*mutex, *context)
	if err != nil {
		log.Println(err)
		log.Println("unable to lock resource")
		// TODO: might be smart to return false and let another container try
		return true
	}
	log.Println("Acquired lock")
	vmid, err := makePlane(message.Plane)
	if err != nil {
		log.Println(err)
		log.Println("unable to execute plane build request")
		_ = locking.UnlockResource(*mutex, *context)
		return true
	}
	log.Printf("Plane %s built and deployed with VMID of %s", message.Plane.Name, vmid)
	_ = locking.UnlockResource(*mutex, *context)
	hostname, err := getFQDNHostname(message.Plane.Name)
	if err != nil {
		log.Println(err)
		log.Println("unable to build FQDN of container")
		// TODO: Possibly queue destruction of container?
		return true
	}
	if len(message.Prep) > 0 {
		//////////////////////////////////////////
		// TODO: Get rid of IP lookup
		ip, err := ipam.GetIPFromHostname(hostname)
		if err != nil {
			log.Println(err)
			log.Println("Failed to lookup IP address in IPAM")
			return true
		}
		//////////////////////////////////////////
		// TODO: Wait on PX status rather than arbitrary timeout
		time.Sleep(45 * time.Second)
		err = prep.DeployPlan(hostname, ip, message.Prep)
		if err != nil {
			log.Println(err)
			log.Println("pre-flight prep failed")
			// TODO: Possibly queue destruction of container?
			return true
		}
	}
	return true
}

func handleDestroyMessage(message ampq.Message) bool {
	log.Println("Destroying plane")
	err := destroyPlane(message.Plane)
	if err != nil {
		log.Println(err)
		log.Println("unable to destroy plane")
		return true
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
