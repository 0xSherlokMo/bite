package talabat

import (
	"bytes"
	"context"
	"encoding/json"
)

// flexStr unmarshals from either a JSON string or number (Talabat is inconsistent
// about whether ids are quoted).
type flexStr string

func (f *flexStr) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) >= 2 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = flexStr(s)
		return nil
	}
	*f = flexStr(string(b)) // number or null → keep literal
	return nil
}

// Profile is the authenticated customer's core profile.
type Profile struct {
	ID              flexStr `json:"id"`
	UserID          flexStr `json:"user_id"`
	Email           string  `json:"email"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	MobileCountry   string  `json:"mobile_country_code"`
	MobileNumber    string  `json:"mobile_number"`
	RegistrationTyp flexStr `json:"registration_type"`
}

// Profile returns the logged-in customer's profile. Doubles as a cheap auth check.
func (c *Client) Profile(ctx context.Context) (*Profile, error) {
	var p Profile
	if err := c.get(ctx, apiHost+"/customers/v1/profile", &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Address is a saved delivery address (Talabat returns these under "result").
type Address struct {
	ID          flexStr `json:"decrypted_id"`
	ProfileName string  `json:"profile_name"`
	Type        flexStr `json:"type"`
	AreaID      int     `json:"area_id"`
	AreaName    string  `json:"area_name"`
	Street      string  `json:"street"`
	Building    flexStr `json:"building_no"`
	Floor       flexStr `json:"floor"`
	Block       flexStr `json:"block"`
	Directions  string  `json:"extra_directions"`
}

type addressesEnvelope struct {
	Addresses []Address `json:"result"`
}

// Addresses returns the customer's saved delivery addresses.
func (c *Client) Addresses(ctx context.Context) ([]Address, error) {
	var env addressesEnvelope
	url := locationHost + "/api/v2/user/addresses?countryId=" + countryID
	if err := c.get(ctx, url, &env); err != nil {
		return nil, err
	}
	return env.Addresses, nil
}
