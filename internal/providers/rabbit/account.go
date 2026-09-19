package rabbit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// UserID extracts the numeric user id from the session's JWT (authorization).
// Rabbit's profile endpoint is /api/users/{id}; the id lives in the token.
func (c *Client) UserID() (string, error) {
	var tok string
	for k, v := range c.session.Headers {
		if strings.ToLower(k) == "authorization" {
			tok = v
			break
		}
	}
	tok = strings.TrimPrefix(strings.TrimPrefix(tok, "Bearer "), "bearer ")
	parts := strings.Split(tok, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("authorization header is not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode JWT: %w", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", err
	}
	for _, k := range []string{"userId", "user_id", "id", "sub", "uid"} {
		if v, ok := claims[k]; ok {
			switch t := v.(type) {
			case float64:
				return fmt.Sprintf("%.0f", t), nil
			case string:
				return t, nil
			}
		}
	}
	return "", fmt.Errorf("no user id claim in JWT")
}

// Profile returns the authenticated customer's profile (raw JSON under data).
func (c *Client) Profile(ctx context.Context) (json.RawMessage, error) {
	id, err := c.UserID()
	if err != nil {
		return nil, err
	}
	var env Envelope
	if err := c.get(ctx, apiHost+"/api/users/"+id, &env); err != nil {
		return nil, err
	}
	return env.Data, nil
}

// ServingMode returns the store/serving info Rabbit picks for the account's
// current location (which store fulfils, working hours, purchase modes).
func (c *Client) ServingMode(ctx context.Context) (json.RawMessage, error) {
	var env Envelope
	if err := c.get(ctx, apiHost+"/api/servingMode", &env); err != nil {
		// some deployments return a bare object, not the envelope
		var raw json.RawMessage
		if err2 := c.get(ctx, apiHost+"/api/servingMode", &raw); err2 == nil {
			return raw, nil
		}
		return nil, err
	}
	if len(env.Data) > 0 {
		return env.Data, nil
	}
	return json.Marshal(env)
}
