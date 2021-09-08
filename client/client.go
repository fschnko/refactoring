package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	c *http.Client

	config Config
}

type Config struct {
	BaseURL string
}

func New(cli *http.Client, config Config) *Client {
	return &Client{
		c:      cli,
		config: config,
	}
}

func (c *Client) request(path string, v interface{}) error {
	resp, err := c.c.Get(c.url(path))
	if err != nil {
		return fmt.Errorf("make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body data: %w", err)
	}

	err = json.Unmarshal(buf, v)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

func (c *Client) url(path string) string {
	return c.config.BaseURL + path
}
