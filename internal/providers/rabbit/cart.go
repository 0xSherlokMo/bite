package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
)

// cartActionReq is the body for PUT /api/cart/products/{id}/action.
type cartActionReq struct {
	StoreID      int     `json:"storeId"`
	Action       string  `json:"action"` // "increment" | "decrement"
	PurchaseMode string  `json:"purchaseMode"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

// cartAction increments or decrements a product in the cart.
func (c *Client) cartAction(ctx context.Context, productID, storeID int, action string, lat, lon float64) (json.RawMessage, error) {
	req := cartActionReq{StoreID: storeID, Action: action, PurchaseMode: "normal", Latitude: lat, Longitude: lon}
	url := fmt.Sprintf("%s/api/cart/products/%d/action", apiHost, productID)
	var raw json.RawMessage
	if err := c.send(ctx, "PUT", url, req, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// AddToCart adds one unit of a product (increment).
func (c *Client) AddToCart(ctx context.Context, productID, storeID int, lat, lon float64) (json.RawMessage, error) {
	return c.cartAction(ctx, productID, storeID, "increment", lat, lon)
}

// RemoveFromCart removes one unit of a product (decrement).
func (c *Client) RemoveFromCart(ctx context.Context, productID, storeID int, lat, lon float64) (json.RawMessage, error) {
	return c.cartAction(ctx, productID, storeID, "decrement", lat, lon)
}

// Cart returns the current cart (raw JSON).
func (c *Client) Cart(ctx context.Context) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.get(ctx, apiHost+"/api/cart_items/v2?purchaseMode=normal", &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
