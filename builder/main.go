package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ampq"
	"github.com/streadway/amqp"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"time"
)

func connectToLockDB() (*concurrency.Session, error) {
	endpoint := os.Getenv("ETCD_URI")
	config := clientv3.Config{
		Endpoints: []string{endpoint},
	}
	cli, err := clientv3.New(config)
	if err != nil {
		log.Println(err)
		return &concurrency.Session{}, errors.New("unable to create etcd client")
	}
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Println(err)
		return &concurrency.Session{}, errors.New("unable to create etcd session")
	}
	return s, nil
}

func getLockableMutex(resource string) (*concurrency.Mutex, *context.Context, error) {
	session, err := connectToLockDB()
	if err != nil {
		log.Println(err)
		return &concurrency.Mutex{}, nil, errors.New("unable to connect to etcd")
	}
	lock := concurrency.NewMutex(session, fmt.Sprintf("/%s/", resource))
	ctx := context.Background()
	return lock, &ctx, nil
}

func lockResource(mutex concurrency.Mutex, context context.Context) error {
	log.Println("Obtaining lock...")
	err := mutex.Lock(context)
	if err != nil {
		log.Println(err)
		return errors.New("unable to lock resource")
	}
	log.Println("Lock obtained")
	return nil
}
func unlockResource(mutex concurrency.Mutex, context context.Context) error {
	log.Println("Unlocking resource...")
	err := mutex.Unlock(context)
	if err != nil {
		log.Println(err)
		return errors.New("unable to unlock resource")
	}
	log.Println("Resource unlocked")
	return nil
}

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
	plane := Plane{}
	err := json.Unmarshal(d.Body, &plane)
	if err != nil {
		log.Println(err)
		log.Println("unable to parse message JSON")
		return true
	}
	log.Printf("Recieved request to build plane %s in %d:%d:%d configuration\n", plane.Name, plane.CPU, plane.RAM, plane.Storage)
	mutex, context, err := getLockableMutex("pve-lxc-create")
	if err != nil {
		log.Println(err)
		log.Println("unable to parse message JSON")
		// TODO: might be smart to return false and let another container try
		return true
	}
	// TODO: remove debugging statements
	log.Println("Acquiring lock on pve-lxc-create...")
	err = lockResource(*mutex, *context)
	if err != nil {
		log.Println(err)
		log.Println("unable to lock resource")
		// TODO: might be smart to return false and let another container try
		return true
	}
	log.Println("Acquired lock")
	vmid, err := makePlane(plane)
	if err != nil {
		log.Println(err)
		log.Println("unable to execute plane build request")
		_ = unlockResource(*mutex, *context)
		return true
	}
	log.Printf("Plane %s built and deployed with VMID of %s", plane.Name, vmid)
	time.Sleep(45 * time.Second)
	log.Println("Destroying plane")
	err = destroyPlane(plane)
	if err != nil {
		log.Println(err)
		log.Println("unable to destroy plane")
		_ = unlockResource(*mutex, *context)
		return true
	}
	_ = unlockResource(*mutex, *context)
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
