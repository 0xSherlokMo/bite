package talabat

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xSherlokMo/bite/internal/secrets"
)

// requiredHeaders must be present for a session to authenticate.
var requiredHeaders = []string{"authorization", "x-device-id"}

// ImportHeaders reads a JSON file of {headerName: value} (e.g. exported from a
// capture of a real app request) and stores it as the Talabat session. It
// validates that the essential auth headers are present.
func ImportHeaders(file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	var headers map[string]string
	if err := json.Unmarshal(b, &headers); err != nil {
		return "", fmt.Errorf("expected a JSON object of header:value pairs: %w", err)
	}
	// Normalize keys to lower-case; drop transport-managed headers.
	clean := make(map[string]string, len(headers))
	for k, v := range headers {
		lk := lower(k)
		if headersManaged[lk] {
			continue
		}
		clean[lk] = v
	}
	for _, h := range requiredHeaders {
		if clean[h] == "" {
			return "", fmt.Errorf("missing required header %q in %s", h, file)
		}
	}
	if err := secrets.Save(&secrets.Session{Provider: "talabat", Headers: clean}); err != nil {
		return "", err
	}
	return secrets.Location("talabat")
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
