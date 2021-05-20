//go:generate stringer -linecomment -type=Rate

package octopusenergy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TariffChargeService handles communication with the tariff related Octopus API.
type TariffChargeService service

// Rate is the type of charge, example standing-charges, day-unit-rates, standard-unit-rates
type Rate int

const (
	RateStandingCharge Rate = iota // standing-charges
	RateStandardUnit               // standard-unit-rates
	RateDayUnit                    // day-unit-rates
	RateNightUnit                  // night-unit-rates
)

// TariffChargesGetOptions is the options for GetTariffCharges.
type TariffChargesGetOptions struct {
	// The code of the product to be retrieved.
	ProductCode string `url:"-"`

	// The code of the tariff to be retrieved.
	TariffCode string `url:"-"`

	// Fueltype: electricity or gas
	FuelType FuelType `url:"-"`

	// The type of charge
	Rate Rate `url:"-"`

	// Show charges active from the given datetime (inclusive). This parameter can be provided on its own.
	PeriodFrom *time.Time `url:"period_from,omitempty" layout:"2006-01-02T15:04:05Z" optional:"true"`

	// Show charges active up to the given datetime (exclusive).
	// You must also provide the period_from parameter in order to create a range.
	PeriodTo *time.Time `url:"period_to,omitempty" layout:"2006-01-02T15:04:05Z" optional:"true"`

	// Page size of returned results.
	// Default is 100, maximum is 1,500 to give up to a month of half-hourly prices.
	PageSize *int `url:"page_size,omitempty" optional:"true"`

	// Pagination page to be returned on this request
	Page *int `url:"page,omitempty" optional:"true"`
}

// TariffChargesGetOutput is the returned struct from GetTariffCharges.
type TariffChargesGetOutput struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		ValueExcVat float64   `json:"value_exc_vat"`
		ValueIncVat float64   `json:"value_inc_vat"`
		ValidFrom   time.Time `json:"valid_from"`
		ValidTo     time.Time `json:"valid_to"`
	} `json:"results"`
}

// Get retrieves the details of a tariffs changes. This endpoint is paginated, it will return
// next and previous links if returned data is larger than the set page size, you are responsible
// to request the next page if required.
func (s *TariffChargeService) Get(options *TariffChargesGetOptions) (*TariffChargesGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context
func (s *TariffChargeService) GetWithContext(ctx context.Context, options *TariffChargesGetOptions) (*TariffChargesGetOutput, error) {
	path := fmt.Sprintf("/v1/products/%s/%s-tariffs/%s/%s/", options.ProductCode, options.FuelType.String(), options.TariffCode, options.Rate.String())
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

	res := TariffChargesGetOutput{}
	if err := s.client.sendRequest(req, false, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetPages same as Get except it returns all pages in one request
func (s *TariffChargeService) GetPages(options *TariffChargesGetOptions) (*TariffChargesGetOutput, error) {
	return s.GetPagesWithContext(context.Background(), options)
}

// GetPagesWithContext same as GetPages except it takes a Context
func (s *TariffChargeService) GetPagesWithContext(ctx context.Context, options *TariffChargesGetOptions) (*TariffChargesGetOutput, error) {
	options.PageSize = Int(1500)
	options.Page = nil

	fullResp := TariffChargesGetOutput{}
	var lastPage bool
	var i int

	for !lastPage {
		page, err := s.GetWithContext(ctx, options)
		if err != nil {
			return nil, err
		}
		fullResp.Count = page.Count
		fullResp.Results = append(fullResp.Results, page.Results...)
		if page.Next == "" {
			lastPage = true
		}
		i++
		options.Page = Int(i)
	}

	return &fullResp, nil
}
