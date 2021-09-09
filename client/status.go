package client

import (
	"fmt"
	"time"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusProcessing
	StatusFailed
	StatusSuccess
)

func (c *Client) Status(token string) (Status, error) {
	var resp struct {
		Message string `json:"message"`
	}

	delay := c.config.StatusGetMinDelay
	counter := 0
	next := func() bool {
		counter++

		if counter == 1 {
			return true
		}

		if counter > c.config.StatusGetAttempts {
			return false
		}

		time.Sleep(delay)
		delay *= time.Duration(c.config.StatusGetDelayFactor)
		if delay > c.config.StatusGetMaxDelay {
			delay = c.config.StatusGetMaxDelay
		}

		return true
	}

	for next() {
		err := c.request("/status/"+token, &resp)
		if err != nil {
			return StatusUnknown, fmt.Errorf("request status: %w", err)
		}

		status := status(resp.Message)

		if status != StatusUnknown {
			return status, nil
		}
	}

	return StatusUnknown, fmt.Errorf("failed %d attempts to get status", c.config.StatusGetAttempts)
}

func status(message string) Status {
	switch message {
	case "processing":
		return StatusProcessing
	case "failed":
		return StatusFailed
	case "success":
		return StatusSuccess
	default:
		return StatusUnknown
	}
}
