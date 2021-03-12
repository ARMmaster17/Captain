package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Wrapper for http.Client that provides a stateful object to make HTTP requests against.
// Provides some boilerplate to reduce code duplication.
type HttpRestClient struct {
	Client *http.Client
	RootUrl string
	ContentType string
	Headers map[string]string
}

// Creates a new HttpRestClient that will make all reqests from RootUrl with the specified settings.
func NewHttpRestClient(rootUrl string, contentType string, ignoreTls bool) (*HttpRestClient, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ignoreTls},
	}
	client := &http.Client{Transport: tr}
	return &HttpRestClient{
		Client: client,
		RootUrl: rootUrl,
		ContentType: contentType,
	}, nil
}

func (h *HttpRestClient) doRequest(request *http.Request) (string, error) {
	for k, v := range h.Headers {
		request.Header.Set(k, v)
	}
	response, err := h.Client.Do(request)
	if err != nil {
		log.Println(err)
		return "", errors.New(fmt.Sprintf("unable to reach endpoint %s", request.URL.Path))
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("target resource %s returned error code %d", request.URL.Path, response.StatusCode))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return "", errors.New(fmt.Sprintf("unable to read data returned from %s", request.URL.Path))
	}
	return string(body), nil
}