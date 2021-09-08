package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientStatus(t *testing.T) {
	t.Run("status success", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"message":"success"}`)
		}))
		defer svr.Close()

		c := New(http.DefaultClient, Config{BaseURL: svr.URL})
		result, err := c.Status("dummy")
		assert.NoError(t, err)
		assert.Equal(t, result, StatusSuccess)
	})

	t.Run("status processing", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"message":"processing"}`)
		}))
		defer svr.Close()

		c := New(http.DefaultClient, Config{BaseURL: svr.URL})
		result, err := c.Status("dummy")
		assert.NoError(t, err)
		assert.Equal(t, result, StatusProcessing)
	})

	t.Run("status failed", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"message":"failed"}`)
		}))
		defer svr.Close()

		c := New(http.DefaultClient, Config{BaseURL: svr.URL})
		result, err := c.Status("dummy")
		assert.NoError(t, err)
		assert.Equal(t, result, StatusFailed)
	})

	t.Run("empty json", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{}`)
		}))
		defer svr.Close()

		c := New(http.DefaultClient, Config{BaseURL: svr.URL})
		result, err := c.Status("dummy")
		assert.Error(t, err)
		assert.Equal(t, result, StatusUnknown)
	})
}
