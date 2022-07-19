// Package eircodesaferplace allows to get the eircode location after querying
// the eircode service.
package eircodesaferplace

import (
	"encoding/json"
	"fmt"
	"net/http"

	"safer.place/internal/address"
)

type Resolver struct {
	c     *http.Client
	addr  string
	token string
}

func New(addr, token string) *Resolver {
	return &Resolver{
		addr:  addr,
		token: token,
		c:     http.DefaultClient,
	}
}

func (r *Resolver) Resolve(addr string) (string, float64, float64, error) {
	req, err := http.NewRequest(http.MethodGet, r.addr+"/"+addr, nil)
	if err != nil {
		return "", 0, 0, fmt.Errorf("unable to create request: %w", err)
	}
	req.Header.Set("Token", r.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.c.Do(req)
	if err != nil {
		return "", 0, 0, fmt.Errorf("unable to request eircode: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", 0, 0, address.ErrUnresolved
	}
	if resp.StatusCode != http.StatusOK {
		return "", 0, 0, fmt.Errorf("unexpected error: %w", err)
	}

	defer resp.Body.Close()

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", 0, 0, fmt.Errorf("unable to decode response")
	}

	return data.Address, data.Lattitude, data.Longitude, nil
}

type Data struct {
	Eircode   string  `json:"eircode"`
	Address   string  `json:"address"`
	Lattitude float64 `json:"lattitide"`
	Longitude float64 `josn:"longitude"`
}
