package octopusenergy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// MeterPointService handles communication with the meter point related Octopus API.
type MeterPointService service

// MeterPointGetOptions is the options for GetMeterPoint.
type MeterPointGetOptions struct {
	// The electricity meter-point’s MPAN.
	MPAN string `url:"-"`
}

// MeterPointGetOutput is the returned struct from GetMeterPoint.
type MeterPointGetOutput struct {
	// Grid Supply Point
	GSP string `json:"gsp"`

	// The electricity meter-point’s MPAN.
	MPAN string `json:"mpan"`

	// ProfileClass
	ProfileClass int `json:"profile_class"`
}

// Get the GSP and profile of a given MPAN.
func (s *MeterPointService) Get(options *MeterPointGetOptions) (*MeterPointGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context
func (s *MeterPointService) GetWithContext(ctx context.Context, options *MeterPointGetOptions) (*MeterPointGetOutput, error) {
	path := fmt.Sprintf("v1/electricity-meter-points/%s", options.MPAN)
	rel := &url.URL{Path: path}
	u := s.client.BaseURL.ResolveReference(rel)
	url, err := addParameters(u, options)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	res := MeterPointGetOutput{}
	if err := s.client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
