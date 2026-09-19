// Package secrets stores per-provider session tokens outside the repository.
//
// Sessions are written to $XDG_CONFIG_HOME/gofer (default ~/.config/gofer) with
// 0600 permissions so credentials never live in source control or world-readable
// locations. Only the raw HTTP header set needed to authenticate is persisted.
package secrets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Session is the minimal auth material for a provider: the exact HTTP headers the
// mobile app sends (authorization JWT, device/anti-fraud tokens, etc.).
type Session struct {
	Provider string            `json:"provider"`
	Headers  map[string]string `json:"headers"`
}

// dir returns the gofer config directory, creating it with 0700 if needed.
func dir() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	d := filepath.Join(base, "gofer")
	if err := os.MkdirAll(d, 0o700); err != nil {
		return "", err
	}
	return d, nil
}

func path(provider string) (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, provider+".json"), nil
}

// Load reads a stored session for the given provider.
func Load(provider string) (*Session, error) {
	p, err := path(provider)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no session for %q: run `gofer %s session import <file>` first", provider, provider)
		}
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("corrupt session file %s: %w", p, err)
	}
	if len(s.Headers) == 0 {
		return nil, fmt.Errorf("session file %s has no headers", p)
	}
	return &s, nil
}

// Save writes a session for the provider with 0600 perms.
func Save(s *Session) error {
	if s.Provider == "" {
		return fmt.Errorf("session has no provider name")
	}
	p, err := path(s.Provider)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	// Write to a temp file then rename, keeping 0600 throughout.
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// Location returns the on-disk path where a provider's session is stored.
func Location(provider string) (string, error) { return path(provider) }
