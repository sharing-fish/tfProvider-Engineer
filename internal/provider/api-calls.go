package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
}

func (c *Client) GetEngineers() ([]EngineerModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineers := []EngineerModel{}
	err = json.Unmarshal(body, &engineers)
	if err != nil {
		return nil, err
	}

	return engineers, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	// Implementation of the request execution
	// For example, using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// NewClient -
func NewClient(endpoint *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
		Endpoint: *endpoint,
	}

	return &c, nil
}
