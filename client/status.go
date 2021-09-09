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

	next := backoffRepeater(c.config.StatusGetAttempts, c.config.StatusGetMinDelay, c.config.StatusGetMaxDelay, c.config.StatusGetDelayFactor)

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

func backoffRepeater(attempts int, min, max time.Duration, factor int) func() bool {
	delay := min
	counter := 0
	return func() bool {
		counter++

		if counter == 1 {
			return true
		}

		if counter > attempts {
			return false
		}

		time.Sleep(delay)
		delay *= time.Duration(factor)
		if delay > max {
			delay = max
		}

		return true
	}
}
