package CaptainLib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// restGET performs a GET request against the given path at the base URL.
func (c *CaptainClient) restGET(path string) ([]byte, error) {
	//return c.restRequest("GET", path, map[string]interface{}{})
	response, err := http.Get(fmt.Sprintf("%s%s", c.BaseUrl, path))
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request:\n%w", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s response body:\n%w", "GET", err)
	}
	return body, nil
}

// restPOST performs a POST request against the given path at the base URL.
func (c *CaptainClient) restPOST(path string, payload map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to convert %s payload to JSON:\n%w", "POST", err)
	}
	response, err := http.Post(fmt.Sprintf("%s%s", c.BaseUrl, path), "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unable to make GET request:\n%w", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s response body:\n%w", "GET", err)
	}
	return body, nil
}

// restPUT performs a PUT request with the given payload at the specified path appended to the base URL.
func (c *CaptainClient) restPUT(path string, payload map[string]interface{}) ([]byte, error) {
	return c.restRequest("PUT", path, payload)
}

// restDELETE sends a DELETE request to the given path at the base URL.
func (c *CaptainClient) restDELETE(path string) ([]byte, error) {
	return c.restRequest(http.MethodDelete, path, nil)
}

// restRequest is a generic REST request handler that performs the specified action at the specified URL.
func (c *CaptainClient) restRequest(method string, uri string, payload map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseUrl, uri)
	fmt.Println(url)
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to convert %s payload to JSON:\n%w", method, err)
	}
	request, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unable to create %s request:\n%w", method, err)
	}
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to make %s request:\n%w", method, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s response body:\n%w", method, err)
	}
	return body, nil
}
