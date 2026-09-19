package talabat

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// Restaurant is a food vendor from the composite listing. Talabat uses terse,
// occasionally inconsistent JSON keys, so we extract tolerantly.
type Restaurant struct {
	VendorID     int     `json:"vendorId"`
	BranchID     int     `json:"branchId"`
	Name         string  `json:"name"`
	Rating       float64 `json:"rating"`
	RatingsCount int     `json:"ratingsCount"`
	Cuisine      string  `json:"cuisine"`
	DeliveryText string  `json:"deliveryText"`
	TimeEstimate string  `json:"timeEstimate"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	IsTalabatPro bool    `json:"isTalabatPro"`
}

type restaurantListEnvelope struct {
	Vendors      []map[string]json.RawMessage `json:"vendors"`
	TotalVendors int                          `json:"total_vendors"`
}

// Restaurants lists food vendors that deliver to the given coordinates. areaID is
// Talabat's area identifier for the location (see Address.AreaID); page starts at 1.
func (c *Client) Restaurants(ctx context.Context, lat, lon float64, areaID, page int) ([]Restaurant, error) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf("%s/vendor-list/v1/composite-list/%g/%g?countrycode=%s&areaid=%d&vertical_id=0&isCustomerPro=true&page=%d&size=40",
		apiHost, lat, lon, countryID, areaID, page)
	var env restaurantListEnvelope
	if err := c.get(ctx, url, &env); err != nil {
		return nil, err
	}
	out := make([]Restaurant, 0, len(env.Vendors))
	for _, v := range env.Vendors {
		out = append(out, Restaurant{
			VendorID:     rawInt(v["id"]),
			BranchID:     rawInt(v["bid"]),
			Name:         rawStr(v["na"]),
			Rating:       rawFloat(v["rat"]),
			RatingsCount: rawInt(v["ratings_count"]),
			Cuisine:      rawStr(v["cus"]),
			DeliveryText: rawStr(v["delivery_text"]),
			TimeEstimate: rawStr(v["time_estimation"]),
			Latitude:     rawFloat(v["latitude"]),
			Longitude:    rawFloat(v["longitude"]),
			IsTalabatPro: rawBool(v["isTalabatPro"]),
		})
	}
	return out, nil
}

// Menu returns a restaurant branch's menu. The payload is deep and varies by
// vendor, so the sections/items are returned as raw JSON under Result for the
// caller to render or index. (Typed menu modeling is deferred: see docs.)
type Menu struct {
	BaseURL string          `json:"baseUrl"`
	Result  json.RawMessage `json:"result"`
}

type menuEnvelope struct {
	Menu Menu `json:"menu"`
}

// Menu fetches the menu for a restaurant branch (use Restaurant.BranchID).
func (c *Client) Menu(ctx context.Context, branchID int) (*Menu, error) {
	url := fmt.Sprintf("%s/menubff/v4/branches/%d/menu", apiHost, branchID)
	var env menuEnvelope
	if err := c.get(ctx, url, &env); err != nil {
		return nil, err
	}
	return &env.Menu, nil
}

// --- tolerant JSON extractors (Talabat mixes string/number types) ---

func rawStr(r json.RawMessage) string {
	if len(r) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(r, &s) == nil {
		return s
	}
	return string(r)
}

func rawInt(r json.RawMessage) int {
	if len(r) == 0 {
		return 0
	}
	var n int
	if json.Unmarshal(r, &n) == nil {
		return n
	}
	var s string
	if json.Unmarshal(r, &s) == nil {
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
	}
	return 0
}

func rawFloat(r json.RawMessage) float64 {
	if len(r) == 0 {
		return 0
	}
	var f float64
	if json.Unmarshal(r, &f) == nil {
		return f
	}
	var s string
	if json.Unmarshal(r, &s) == nil {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return 0
}

func rawBool(r json.RawMessage) bool {
	if len(r) == 0 {
		return false
	}
	var b bool
	_ = json.Unmarshal(r, &b)
	return b
}
