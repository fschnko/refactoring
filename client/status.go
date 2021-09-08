package client

import "fmt"

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

	for i := 0; i < c.config.StatusGetAttempts; i++ {
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
