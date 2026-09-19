// Package rabbit is a pure-API client for Rabbit (rabbitmart.com), Egypt's
// ~15-minute grocery + food delivery app. Rabbit ships no anti-tamper and pins
// nothing, so requests replay with the app's header set: an `authorization` JWT
// plus a device id and a few metadata headers (no anti-fraud tokens).
package rabbit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xSherlokMo/bite/internal/secrets"
)

const (
	apiHost  = "https://api-prod.rabbitmart.com"
	foodHost = "https://food-portal-api-prod.rabbitmart.com"
)

var headersManaged = map[string]bool{
	"host": true, "http2-enabled": true, "content-length": true, "accept-encoding": true,
}

// Client is a Rabbit API client bound to one stored session.
type Client struct {
	http    *http.Client
	session *secrets.Session
}

// New loads the stored Rabbit session and returns a client.
func New() (*Client, error) {
	s, err := secrets.Load("rabbit")
	if err != nil {
		return nil, err
	}
	return &Client{http: &http.Client{Timeout: 30 * time.Second}, session: s}, nil
}

// Envelope is Rabbit's standard response wrapper.
type Envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

// APIError carries a non-2xx response.
type APIError struct {
	Status      int
	Method, URL string
	Body        string
}

func (e *APIError) Error() string {
	b := e.Body
	if len(b) > 200 {
		b = b[:200] + "…"
	}
	return fmt.Sprintf("rabbit %s %s -> %d: %s", e.Method, e.URL, e.Status, b)
}

// Expired reports an auth failure (token needs refresh).
func (e *APIError) Expired() bool { return e.Status == 401 || e.Status == 403 }

func (c *Client) do(ctx context.Context, method, url string, body []byte, out any) error {
	var rdr io.Reader
	if body != nil {
		rdr = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil {
		return err
	}
	for k, v := range c.session.Headers {
		if headersManaged[strings.ToLower(k)] {
			continue
		}
		req.Header.Set(k, v)
	}
	if body != nil {
		req.Header.Set("content-type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Status: resp.StatusCode, Method: method, URL: url, Body: string(raw)}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func (c *Client) get(ctx context.Context, url string, out any) error {
	return c.do(ctx, http.MethodGet, url, nil, out)
}

func (c *Client) send(ctx context.Context, method, url string, in, out any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return c.do(ctx, method, url, b, out)
}
