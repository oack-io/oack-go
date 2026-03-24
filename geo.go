package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Checker represents a network checker node.
type Checker struct {
	ID      string `json:"id"`
	Region  string `json:"region"`
	Country string `json:"country"`
	IP      string `json:"ip"`
	ASN     any    `json:"asn"`
	Mode    string `json:"mode"`
	Status  string `json:"status"`
}

// GeoRegion represents a geographic region with its countries.
type GeoRegion struct {
	Code      string       `json:"code"`
	Name      string       `json:"name"`
	Countries []GeoCountry `json:"countries"`
}

// GeoCountry represents a country within a region.
type GeoCountry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GeoRegionsResponse wraps the regions list response.
type GeoRegionsResponse struct {
	Regions []GeoRegion `json:"regions"`
}

// ListCheckers returns all available checker nodes.
func (c *Client) ListCheckers(ctx context.Context) ([]Checker, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/checkers", nil)
	if err != nil {
		return nil, err
	}
	var checkers []Checker
	if err := json.Unmarshal(body, &checkers); err != nil {
		return nil, fmt.Errorf("unmarshal checkers: %w", err)
	}
	return checkers, nil
}

// ListRegions returns all available geographic regions and their countries.
func (c *Client) ListRegions(ctx context.Context) (*GeoRegionsResponse, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/regions", nil)
	if err != nil {
		return nil, err
	}
	var resp GeoRegionsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal regions: %w", err)
	}
	return &resp, nil
}
