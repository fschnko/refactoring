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
	sleep := func() {
		time.Sleep(delay)

		delay *= time.Duration(c.config.StatusGetDelayFactor)
		if delay > c.config.StatusGetMaxDelay {
			delay = c.config.StatusGetMaxDelay
		}
	}

	for i := 0; i < c.config.StatusGetAttempts; i++ {
		err := c.request("/status/"+token, &resp)
		if err != nil {
			return StatusUnknown, fmt.Errorf("request status: %w", err)
		}

		status := status(resp.Message)

		if status != StatusUnknown {
			return status, nil
		}

		sleep()
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
