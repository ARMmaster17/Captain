package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Performs a POST request against the specified URI and []byte body payload. Utilizes state of referenced HttpRestClient
// object and returns the contents of the response body in string form.
func (h *HttpRestClient) Post(url string, reqBody []byte) (string, error) {
	fullUrl := h.RootUrl + url
	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build request for POST endpoint " + fullUrl)
	}
	return h.doRequest(request)
}

func (h *HttpRestClient) PostJson(url string, reqBody interface{}) (string, error) {
	reqBodyParsed, err := json.Marshal(reqBody)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to parse JSON payload to be delivered")
	}
	return h.Post(url, reqBodyParsed)
}

func (h *HttpRestClient) PostUrlEncode(urlPath string, reqBody interface{}) (string, error) {
	fullUrl := h.RootUrl + urlPath
	encoder := schema.NewEncoder()
	form := url.Values{}
	err := encoder.Encode(reqBody, form)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to parse url encoded payload")
	}
	values := form.Encode()
	request, err := http.NewRequest("POST", urlPath, strings.NewReader(values))
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build request for POST endpoint " + fullUrl)
	}
	request.Form = form
	return h.doRequest(request)
}

func (h *HttpRestClient) PostWithBasicAuth(url string, username string, password string, reqBody []byte) (string, error) {
	fullUrl := h.RootUrl + url
	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build request for POST endpoint " + fullUrl)
	}
	request.SetBasicAuth(username, password)
	return h.doRequest(request)
}