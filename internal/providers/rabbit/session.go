package rabbit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/0xSherlokMo/bite/internal/secrets"
)

var requiredHeaders = []string{"authorization", "deviceid"}

// ImportHeaders stores a captured Rabbit header set as the session.
func ImportHeaders(file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	var headers map[string]string
	if err := json.Unmarshal(b, &headers); err != nil {
		return "", fmt.Errorf("expected JSON object of header:value pairs: %w", err)
	}
	clean := make(map[string]string, len(headers))
	for k, v := range headers {
		lk := strings.ToLower(k)
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
	if err := secrets.Save(&secrets.Session{Provider: "rabbit", Headers: clean}); err != nil {
		return "", err
	}
	return secrets.Location("rabbit")
}
