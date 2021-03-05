package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getAuthData(url string, username string, password string) (csrf string, token string) {
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		print(err)
		return
	}
	response, err := http.Post(url + "access/ticket", "application/json", bytes.NewBuffer(reqBody))
	defer response.Body.Close()
	var data map[string]map[string]string
	json.NewDecoder(response.Body).Decode(&data)
	csrf = data["data"]["CSRFPreventionToken"]
	token = data["data"]["ticket"]
	return
}
