package http

import (
	"errors"
	"log"
	"net/http"
)

func (h *HttpRestClient) Delete(url string) (string, error) {
	fullUrl := h.RootUrl + url
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build request for DELETE endpoint " + fullUrl)
	}
	return h.doRequest(request)
}