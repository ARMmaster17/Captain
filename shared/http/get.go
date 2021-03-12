package http

import (
	"bytes"
	"errors"
	"log"
	"net/http"
)

func (h *HttpRestClient) Get(url string) (string, error) {
	fullUrl := h.RootUrl + url
	request, err := http.NewRequest("GET", url, bytes.NewReader(nil))
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build request for GET endpoint " + fullUrl)
	}
	return h.doRequest(request)
}