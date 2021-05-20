package octopusenergy

import (
	"context"
	"net/http"
	"net/url"
)

// GridSupplyPointService handles communication with the grid supply point related Octopus API.
type GridSupplyPointService service

// GridSupplyPointGetOptions is the options for GetGridSupplyPoint.
type GridSupplyPointGetOptions struct {
	// A postcode to filter on.
	// If Octopus are unable to map the passed postcode to a GSP, an empty list will be returned.
	Postcode *string `url:"postcode,omitempty" optional:"true"`
}

// GridSupplyPointGetOutput is the returned struct from GetGridSupplyPoint.
type GridSupplyPointGetOutput struct {
	Count   int `json:"count"`
	Results []struct {
		GroupID string `json:"group_id"`
	} `json:"results"`
}

// Get gets the GSP and group ID, filtered by postcode if one is given.
func (s *GridSupplyPointService) Get(options *GridSupplyPointGetOptions) (*GridSupplyPointGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context.
func (s *GridSupplyPointService) GetWithContext(ctx context.Context, options *GridSupplyPointGetOptions) (*GridSupplyPointGetOutput, error) {
	rel := &url.URL{Path: "v1/industry/grid-supply-points"}
	u := s.client.BaseURL.ResolveReference(rel)
	url, err := addParameters(u, options)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	res := GridSupplyPointGetOutput{}
	if err := s.client.sendRequest(req, false, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
