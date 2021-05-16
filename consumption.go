package octopusenergy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ConsumptionService handles communication with the consumption related Octopus API.
type ConsumptionService service

// ConsumptionGetOptions is the options for GetConsumption.
type ConsumptionGetOptions struct {
	// The Meter Point Number this is the electricity meter-point’s MPAN or gas meter-point’s MPRN
	MPN string `url:"-"`

	// The meter’s serial number.
	SerialNumber string `url:"-"`

	// Fueltype: electricity or gas
	FuelType FuelType `url:"-"`

	// Show consumption from the given datetime (inclusive). This parameter can be provided on its own.
	PeriodFrom *time.Time `url:"period_from,omitempty" layout:"2006-01-02T15:04:05Z" optional:"true"`

	// Show consumption to the given datetime (exclusive).
	// This parameter also requires providing the period_from parameter to create a range.
	PeriodTo *time.Time `url:"period_to,omitempty" layout:"2006-01-02T15:04:05Z" optional:"true"`

	// Page size of returned results.
	// Default is 100, maximum is 25,000 to give a full year of half-hourly consumption details.
	PageSize *int `url:"page_size,omitempty" optional:"true"`

	// Ordering of results returned.
	// Default is that results are returned in reverse order from latest available figure.
	// Valid values: * ‘period’, to give results ordered forward. * ‘-period’, (default), to give results ordered from most recent backwards.
	OrderBy *string `url:"order_by,omitempty" optional:"true"`

	// Aggregates consumption over a specified time period.
	// A day is considered to start and end at midnight in the server’s timezone.
	// The default is that consumption is returned in half-hour periods. Accepted values are: * ‘hour’ * ‘day’ * ‘week’ * ‘month’ * ‘quarter’
	GroupBy *string `url:"group_by,omitempty" optional:"true"`

	// Pagination page to be returned on this request
	Page *int `url:"page,omitempty" optional:"true"`
}

// ConsumptionGetOutput is the returned struct from GetConsumption.
type ConsumptionGetOutput struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Consumption   float64 `json:"consumption"`
		IntervalStart string  `json:"interval_start"`
		IntervalEnd   string  `json:"interval_end"`
	} `json:"results"`
}

// Get consumption data for give meter details. This endpoint is paginated, it will return
// next and previous links if returned data is larger than the set page size, you are responsible
// to request the next page if required.
func (s *ConsumptionService) Get(options *ConsumptionGetOptions) (*ConsumptionGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context.
func (s *ConsumptionService) GetWithContext(ctx context.Context, options *ConsumptionGetOptions) (*ConsumptionGetOutput, error) {
	path := fmt.Sprintf("v1/%s-meter-points/%s/meters/%s/consumption", options.FuelType.String(), options.MPN, options.SerialNumber)
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

	res := ConsumptionGetOutput{}
	if err := s.client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetPages same as Get except it returns all pages in one request.
func (s *ConsumptionService) GetPages(options *ConsumptionGetOptions) (*ConsumptionGetOutput, error) {
	return s.GetPagesWithContext(context.Background(), options)
}

// GetPagesWithContext same as GetPages except it takes a Context.
func (s *ConsumptionService) GetPagesWithContext(ctx context.Context, options *ConsumptionGetOptions) (*ConsumptionGetOutput, error) {
	options.PageSize = Int(1500)
	options.Page = nil

	fullResp := ConsumptionGetOutput{}
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
