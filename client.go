package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HelloClient struct {
	serviceUrl string
	httpClient *http.Client
}

type HelloMessageResponse struct {
	Message string
}

func NewHelloClient(serviceUrl string) *HelloClient {
	return &HelloClient{
		serviceUrl: serviceUrl,
		httpClient: &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		},
	}
}

func (c *HelloClient) GetHello(name string) (message string, err error) {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/hello", c.serviceUrl), nil)
	req.Header.Add("Accept", "application/json")

	if name != "" {
		q := req.URL.Query()
		q.Add("name", name)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Check the status code is what we expect.
	if status := resp.StatusCode; status != http.StatusOK {
		return message, fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return message, err
	}

	// Unmarshal JSON response
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return message, fmt.Errorf("invalid JSON response received: got '%s'", body)
	}

	// Get the message
	message, ok := data["message"].(string)
	if !ok {
		return message, fmt.Errorf("response has no 'message' field: got '%s'", data)
	}

	return message, nil
}
