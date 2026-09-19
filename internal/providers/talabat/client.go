// Package talabat is a pure-API client for Talabat (Delivery Hero) covering the
// account, food (restaurants) and grocery (talabat mart) verticals.
//
// It authenticates by replaying the exact header set the mobile app sends — an
// `authorization` JWT plus device/anti-fraud tokens (incognia-token, x-device-id,
// x-segmenttoken, ...). Those live in a secrets.Session, never in code.
package talabat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/salah/gofer/internal/secrets"
)

const (
	apiHost      = "https://api.talabat.com"
	locationHost = "https://userlocation.talabat.com"

	// EG = Egypt. countryID 9 / global entity HF_EG are Talabat's identifiers.
	countryCode  = "EG"
	countryID    = "9"
	globalEntity = "HF_EG"
)

// headersManaged are set by the transport and must not be replayed verbatim.
var headersManaged = map[string]bool{
	"host":            true,
	"http2-enabled":   true,
	"content-length":  true,
	"accept-encoding": true,
}

// Client is a Talabat API client bound to one stored session.
type Client struct {
	http    *http.Client
	session *secrets.Session
}

// New loads the stored Talabat session and returns a client.
func New() (*Client, error) {
	s, err := secrets.Load("talabat")
	if err != nil {
		return nil, err
	}
	return &Client{
		http:    &http.Client{Timeout: 30 * time.Second},
		session: s,
	}, nil
}

// APIError carries a non-2xx response for callers that want to branch on status
// (e.g. 401 → token expired → re-import session).
type APIError struct {
	Status int
	Method string
	URL    string
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("talabat %s %s -> %d: %s", e.Method, e.URL, e.Status, truncate(e.Body, 200))
}

// Expired reports whether the error looks like an auth failure.
func (e *APIError) Expired() bool {
	return e.Status == http.StatusUnauthorized || e.Status == http.StatusForbidden
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// do issues a request with the session headers attached and decodes JSON into out.
// out may be nil to discard the body. body may be nil for GET.
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
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode %s %s: %w", method, url, err)
	}
	return nil
}

func (c *Client) get(ctx context.Context, url string, out any) error {
	return c.do(ctx, http.MethodGet, url, nil, out)
}

func (c *Client) postJSON(ctx context.Context, url string, in, out any) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPost, url, b, out)
}
