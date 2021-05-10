package CaptainLib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func (c *CaptainClient) restPUT(path string, payload map[string]interface{}) ([]byte, error) {
	return c.restRequest("PUT", path, payload)
}

func (c *CaptainClient) restDELETE(path string) ([]byte, error) {
	return c.restRequest("DELETE", path, nil)
}

func (c *CaptainClient) restRequest(method string, uri string, payload map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseUrl, uri)
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