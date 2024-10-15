package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Endpoint   string
	HTTPClient *http.Client
}

// EngineerModel
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

func (c *Client) GetEngineerById(id string) (*EngineerModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.Endpoint, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	engineer := EngineerModel{}
	err = json.Unmarshal(body, &engineer)
	if err != nil {
		return nil, err
	}

	return &engineer, nil
}

func (c *Client) CreateEngineer(name, email string) (*EngineerModel, error) {
	engineer := EngineerModel{
		Name:  name,
		Email: email,
	}

	engineerBytes, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.Endpoint), bytes.NewBuffer(engineerBytes))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newEngineer := EngineerModel{}
	err = json.Unmarshal(body, &newEngineer)
	if err != nil {
		return nil, err
	}

	return &newEngineer, nil
}

func (c *Client) DeleteEngineer(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/engineers/%s", c.Endpoint, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateEngineer(id, name, email string) (*EngineerModel, error) {
	engineer := EngineerModel{
		Name:  name,
		Email: email,
	}

	engineerBytes, err := json.Marshal(engineer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.Endpoint, id), bytes.NewBuffer(engineerBytes))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedEngineer := EngineerModel{}
	err = json.Unmarshal(body, &updatedEngineer)
	if err != nil {
		return nil, err
	}

	return &updatedEngineer, nil
}

// DevModel
func (c *Client) GetDevs() ([]DevModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dev := []DevModel{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (c *Client) GetDevById(id string) (*DevModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/id/%s", c.Endpoint, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	dev := DevModel{}
	err = json.Unmarshal(body, &dev)
	if err != nil {
		return nil, err
	}

	return &dev, nil
}

func (c *Client) CreateDev(name string, engineers []EngineerModel) (*DevModel, error) {
	dev := DevModel{
		Name:      name,
		Engineers: engineers,
	}

	devBytes, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating Dev: %s", devBytes) // Log the request payload

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/dev", c.Endpoint), bytes.NewBuffer(devBytes))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	log.Printf("Response Body: %s", body) // Log the response body

	newDev := DevModel{}
	err = json.Unmarshal(body, &newDev)
	if err != nil {
		return nil, err
	}

	return &newDev, nil
}

func (c *Client) UpdateDev(id, name string, engineers []EngineerModel) (*DevModel, error) {
	dev := DevModel{
		Name:      name,
		Engineers: engineers,
	}

	devBytes, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/dev/%s", c.Endpoint, id), bytes.NewBuffer(devBytes))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedDev := DevModel{}
	err = json.Unmarshal(body, &updatedDev)
	if err != nil {
		return nil, err
	}

	return &updatedDev, nil
}

func (c *Client) DeleteDev(id string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/dev/%s", c.Endpoint, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// general purpose client/request functions
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
