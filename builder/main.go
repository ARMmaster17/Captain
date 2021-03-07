package main

import (
	"errors"
	"fmt"
	"log"
)

func makePlane(plane Plane) (string, error) {
	machineConfig, err := buildPlaneConfig(plane)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane configuration")
	}
	vmid, err := createLxcContainer(machineConfig)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane")
	}
	return vmid, nil
}

func main() {
	plane := Plane{
		Name:    "test2",
		CPU:     1,
		RAM:     256,
		Storage: 5,
	}
	vmid, err := makePlane(plane)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(vmid)
}
