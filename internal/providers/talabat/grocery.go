package talabat

import (
	"context"
	"encoding/json"
	"fmt"
)

// VendorRef identifies a grocery store (talabat mart / darkstore) plus the
// customer's coordinates. All fields are required by the cart API.
type VendorRef struct {
	DHVendorID string  // dh_vendor_id (uuid)
	BranchID   string  // branch_id (numeric string, e.g. "760349")
	ChainID    string  // chain_id  (e.g. "607445")
	Lat        float64 // customer latitude
	Lon        float64 // customer longitude
}

// CartItem is a line in a grocery cart request.
type CartItem struct {
	ProductID       string   `json:"product_id"`
	Quantity        int      `json:"quantity"`
	IncrementWeight *float64 `json:"increment_weight"`
}

// cartRequest is the exact body Talabat's app posts to /grocery/v2/cart.
type cartRequest struct {
	DHVendorID             string     `json:"dh_vendor_id"`
	BranchID               string     `json:"branch_id"`
	ChainID                string     `json:"chain_id"`
	GlobalEntityID         string     `json:"global_entity_id"`
	Items                  []CartItem `json:"items"`
	Lat                    float64    `json:"lat"`
	Lon                    float64    `json:"lon"`
	CampaignProgressReq    bool       `json:"campaign_progress_required"`
	HasOutstandingPayments bool       `json:"has_outstanding_payments"`
}

// Cart is the server's view of the grocery cart after an operation.
type Cart struct {
	ID                    string          `json:"id"`
	SubtotalAfterDiscount float64         `json:"subtotal_after_discount"`
	Items                 json.RawMessage `json:"items"`
	Messages              json.RawMessage `json:"cart_update_messages"`
}

// SetCart replaces the grocery cart contents with items (POST /grocery/v2/cart).
// Talabat's mart cart is declarative: you post the full desired item set.
func (c *Client) SetCart(ctx context.Context, v VendorRef, items []CartItem) (*Cart, error) {
	if v.DHVendorID == "" || v.BranchID == "" || v.ChainID == "" {
		return nil, fmt.Errorf("VendorRef requires DHVendorID, BranchID and ChainID")
	}
	req := cartRequest{
		DHVendorID:             v.DHVendorID,
		BranchID:               v.BranchID,
		ChainID:                v.ChainID,
		GlobalEntityID:         globalEntity,
		Items:                  items,
		Lat:                    v.Lat,
		Lon:                    v.Lon,
		CampaignProgressReq:    true,
		HasOutstandingPayments: false,
	}
	var cart Cart
	if err := c.postJSON(ctx, apiHost+"/grocery/v2/cart", req, &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

// AddToCart is a convenience for adding a single product with a quantity.
func (c *Client) AddToCart(ctx context.Context, v VendorRef, productID string, qty int) (*Cart, error) {
	return c.SetCart(ctx, v, []CartItem{{ProductID: productID, Quantity: qty}})
}

// GetCart fetches the current grocery cart for a vendor (GET /grocery/v2/cart).
func (c *Client) GetCart(ctx context.Context, v VendorRef) (*Cart, error) {
	url := fmt.Sprintf(
		"%s/grocery/v2/cart?branch_id=%s&dh_vendor_id=%s&global_entity_id=%s&chain_id=%s&lat=%g&lon=%g&campaign_progress_required=true&has_outstanding_payments=false",
		apiHost, v.BranchID, v.DHVendorID, globalEntity, v.ChainID, v.Lat, v.Lon,
	)
	var cart Cart
	if err := c.get(ctx, url, &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

// Catalog returns a mart store's landing catalog (raw JSON: components/tabs).
// Typed modeling is deferred; callers can render or index the raw payload.
func (c *Client) Catalog(ctx context.Context, branchID string, lat, lon float64) (json.RawMessage, error) {
	url := fmt.Sprintf(
		"%s/grocery/v1/bff/eg/screens/vendor-home/%s?areaId=0&pageNumber=0&lat=%g&lon=%g&aisle_tab_wrapped_indicator_enabled=false",
		apiHost, branchID, lat, lon,
	)
	var raw json.RawMessage
	if err := c.get(ctx, url, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
