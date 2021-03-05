package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Machine struct {
	Name	string	`json:"Name"`
}

func health(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Fprintf(w, "ID: " + key)
	fmt.Println("Endpoint hit: homePage")
}

func createMachine(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var machine Machine
	json.Unmarshal(reqBody, &machine)

	json.NewEncoder(w).Encode(machine)

	fmt.Println(string(reqBody))
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", health)
	router.HandleFunc("/machines/create", createMachine).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	handleRequests()
}
