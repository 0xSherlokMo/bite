package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Restaurants returns the food (Rabbit Kitchen / restaurants) swimlane for a
// location + store (raw JSON). storeName is the store code from ServingMode.
func (c *Client) Restaurants(ctx context.Context, lat, lon float64, storeName string, limit int) (json.RawMessage, error) {
	if limit <= 0 {
		limit = 20
	}
	u := fmt.Sprintf("%s/restaurant-management-service/api/v1/restaurants_swimlane?lat=%g&long=%g&storeName=%s&status=active&type=servingMode&limit=%d",
		foodHost, lat, lon, url.QueryEscape(storeName), limit)
	var raw json.RawMessage
	if err := c.get(ctx, u, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
