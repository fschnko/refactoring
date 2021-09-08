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

	err := c.request("/status/"+token, &resp)
	if err != nil {
		return StatusUnknown, fmt.Errorf("request status: %w", err)
	}

	return status(resp.Message), nil
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
