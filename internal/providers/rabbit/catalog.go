package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// StoreMode identifies a Rabbit store + purchase mode (normal / bulk / food).
type StoreMode struct {
	StoreID      int    `json:"storeId"`
	StoreName    string `json:"storeName"`
	PurchaseMode string `json:"purchaseMode"`
}

type categoryTileReq struct {
	StoresModes []StoreMode `json:"storesModes"`
	Type        string      `json:"type"`
}

// Categories returns the category tiles for a store (raw JSON). storeName is
// Rabbit's store code (e.g. "EGY010SOD"); get it from ServingMode.
func (c *Client) Categories(ctx context.Context, storeID int, storeName string) (json.RawMessage, error) {
	req := categoryTileReq{
		StoresModes: []StoreMode{{storeID, storeName, "normal"}},
		Type:        "EntryPointBranded",
	}
	var env Envelope
	err := c.send(ctx, "POST", apiHost+"/catalog-management-service/api/v1/categoryTile/byTypeStoreMode", req, &env)
	if err != nil {
		return nil, err
	}
	if len(env.Data) > 0 {
		return env.Data, nil
	}
	return json.Marshal(env)
}

// Swimlanes returns the merchandised product rows for a store (raw JSON).
func (c *Client) Swimlanes(ctx context.Context, storeName string, productSize int) (json.RawMessage, error) {
	if productSize <= 0 {
		productSize = 10
	}
	u := fmt.Sprintf("%s/catalog-management-service/api/v1/swimlane?purchaseMode=normal&storeName=%s&productSize=%d",
		apiHost, url.QueryEscape(storeName), productSize)
	var raw json.RawMessage
	if err := c.get(ctx, u, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
