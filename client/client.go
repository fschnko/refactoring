package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	c *http.Client

	config Config
}

type Config struct {
	BaseURL           string
	StatusGetAttempts int
	StatusGetDelay    time.Duration
}

func New(cli *http.Client, config Config) *Client {
	setConfigDefaults(&config)

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

func setConfigDefaults(config *Config) {
	if config.StatusGetAttempts < 1 {
		config.StatusGetAttempts = 1
	}

	if config.StatusGetDelay == 0 {
		config.StatusGetDelay = time.Millisecond
	}
}
