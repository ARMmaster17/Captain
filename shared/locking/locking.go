package locking

import (
	"context"
	"errors"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
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

func GetLockableMutex(resource string) (*concurrency.Mutex, *context.Context, error) {
	session, err := connectToLockDB()
	if err != nil {
		log.Println(err)
		return &concurrency.Mutex{}, nil, errors.New("unable to connect to etcd")
	}
	lock := concurrency.NewMutex(session, fmt.Sprintf("/%s/", resource))
	ctx := context.Background()
	return lock, &ctx, nil
}

func LockResource(mutex concurrency.Mutex, context context.Context) error {
	log.Println("Obtaining lock...")
	err := mutex.Lock(context)
	if err != nil {
		log.Println(err)
		return errors.New("unable to lock resource")
	}
	log.Println("Lock obtained")
	return nil
}
func UnlockResource(mutex concurrency.Mutex, context context.Context) error {
	log.Println("Unlocking resource...")
	err := mutex.Unlock(context)
	if err != nil {
		log.Println(err)
		return errors.New("unable to unlock resource")
	}
	log.Println("Resource unlocked")
	return nil
}
