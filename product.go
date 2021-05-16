package octopusenergy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ProductService handles communication with the product related Octopus API.
type ProductService service

// ProductsListOptions is the options for ListProduct.
type ProductsListOptions struct {
	// Show only variable products.
	IsVariable *bool `url:"is_variable,omitempty" optional:"true"`

	// Show only green products.
	IsGreen *bool `url:"is_green,omitempty" optional:"true"`

	// Show only tracker products.
	IsTracker *bool `url:"is_tracker,omitempty" optional:"true"`

	// Show only pre-pay products.
	IsPrepay *bool `url:"is_prepay,omitempty" optional:"true"`

	// Show only business products.
	IsBusiness *bool `url:"is_business,omitempty" optional:"true"`

	// Show products available for new agreements on the given datetime.
	// Defaults to current datetime, effectively showing products that are currently available.
	AvailableAt *time.Time `url:"available_at,omitempty" layout:"2006-01-02T15:04:05Z" optional:"true"`

	// Pagination page to be returned on this request
	Page *int `url:"page,omitempty" optional:"true"`
}

// ProductsListOutput is the returned struct from ListProduct.
type ProductsListOutput struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Code          string      `json:"code"`
		FullName      string      `json:"full_name"`
		DisplayName   string      `json:"display_name"`
		Description   string      `json:"description"`
		IsVariable    bool        `json:"is_variable"`
		IsGreen       bool        `json:"is_green"`
		IsTracker     bool        `json:"is_tracker"`
		IsPrepay      bool        `json:"is_prepay"`
		IsBusiness    bool        `json:"is_business"`
		IsRestricted  bool        `json:"is_restricted"`
		Term          int         `json:"term"`
		Brand         string      `json:"brand"`
		AvailableFrom time.Time   `json:"available_from"`
		AvailableTo   interface{} `json:"available_to"`
		Links         []struct {
			Href   string `json:"href"`
			Method string `json:"method"`
			Rel    string `json:"rel"`
		} `json:"links"`
	} `json:"results"`
}

// List return a list of energy products. This endpoint is paginated, it will return
// next and previous links if returned data is larger than the set page size, you are responsible
// to request the next page if required.
func (s *ProductService) List(options *ProductsListOptions) (*ProductsListOutput, error) {
	return s.ListWithContext(context.Background(), options)
}

// ListWithContext same as ListProducts except it takes a Context.
func (s *ProductService) ListWithContext(ctx context.Context, options *ProductsListOptions) (*ProductsListOutput, error) {
	rel := &url.URL{Path: "v1/products/"}
	u := s.client.BaseURL.ResolveReference(rel)
	url, err := addParameters(u, options)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	res := ProductsListOutput{}
	if err := s.client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListPages same as List except it returns all pages in one request.
func (s *ProductService) ListPages(options *ProductsListOptions) (*ProductsListOutput, error) {
	return s.ListPagesWithContext(context.Background(), options)
}

// ListPagesWithContext same as ListPages except it takes a Context.
func (s *ProductService) ListPagesWithContext(ctx context.Context, options *ProductsListOptions) (*ProductsListOutput, error) {
	options.Page = nil

	fullResp := ProductsListOutput{}
	var lastPage bool
	var i int

	for !lastPage {
		page, err := s.ListWithContext(ctx, options)
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

// ProductsGetOptions is the options for GetProduct.
type ProductsGetOptions struct {
	// The code of the product to be retrieved.
	ProductCode string `url:"-"`

	// The point in time in which to show the active charges. Defaults to current datetime.
	TariffsActiveAt *time.Time `url:"tariffs_active_at" layout:"2006-01-02T15:04:05Z" optional:"true"`
}

// ProductsGetOutput is the returned struct from GetProduct.
type ProductsGetOutput struct {
	Code                             string                 `json:"code"`
	FullName                         string                 `json:"full_name"`
	DisplayName                      string                 `json:"display_name"`
	Description                      string                 `json:"description"`
	IsVariable                       bool                   `json:"is_variable"`
	IsGreen                          bool                   `json:"is_green"`
	IsTracker                        bool                   `json:"is_tracker"`
	IsPrepay                         bool                   `json:"is_prepay"`
	IsBusiness                       bool                   `json:"is_business"`
	IsRestricted                     bool                   `json:"is_restricted"`
	Brand                            string                 `json:"brand"`
	Term                             int                    `json:"term"`
	AvailableFrom                    time.Time              `json:"available_from"`
	AvailableTo                      time.Time              `json:"available_to"`
	TariffsActiveAt                  time.Time              `json:"tariffs_active_at"`
	SingleRegisterElectricityTariffs map[string]Tariff      `json:"single_register_electricity_tariffs"`
	DualRegisterElectricityTariffs   map[string]Tariff      `json:"dual_register_electricity_tariffs"`
	SingleRegisterGasTariffs         map[string]Tariff      `json:"single_register_gas_tariffs"`
	SampleQuotes                     map[string]SampleQuote `json:"sample_quotes"`
	SampleConsumption                SampleConsumption      `json:"sample_consumption"`
	Links                            []struct {
		Href   string `json:"href"`
		Method string `json:"method"`
		Rel    string `json:"rel"`
	} `json:"links"`
}

type Tariff struct {
	DirectDebitMonthly   TariffDirectDebit `json:"direct_debit_monthly"`
	DirectDebitQuarterly TariffDirectDebit `json:"direct_debit_quarterly"`
}

type TariffDirectDebit struct {
	Code                   string  `json:"code"`
	StandardUnitRateExcVat float64 `json:"standard_unit_rate_exc_vat"`
	StandardUnitRateIncVat float64 `json:"standard_unit_rate_inc_vat"`
	StandingChargeExcVat   float64 `json:"standing_charge_exc_vat"`
	StandingChargeIncVat   float64 `json:"standing_charge_inc_vat"`
	OnlineDiscountExcVat   int     `json:"online_discount_exc_vat"`
	OnlineDiscountIncVat   int     `json:"online_discount_inc_vat"`
	DualFuelDiscountExcVat int     `json:"dual_fuel_discount_exc_vat"`
	DualFuelDiscountIncVat int     `json:"dual_fuel_discount_inc_vat"`
	ExitFeesExcVat         int     `json:"exit_fees_exc_vat"`
	ExitFeesIncVat         int     `json:"exit_fees_inc_vat"`
	Links                  []struct {
		Href   string `json:"href"`
		Method string `json:"method"`
		Rel    string
	} `json:"direct_debit_monthly"`
}

type SampleQuote struct {
	DirectDebitMonthly   SampleQuoteDirectDebit `json:"direct_debit_monthly"`
	DirectDebitQuarterly SampleQuoteDirectDebit `json:"direct_debit_quarterly"`
}

type SampleQuoteDirectDebit struct {
	ElectricitySingleRate SampleQuoteDirectDebitRate `json:"electricity_single_rate"`
	ElectricityDualRate   SampleQuoteDirectDebitRate `json:"electricity_dual_rate"`
	DualFuelSingleRate    SampleQuoteDirectDebitRate `json:"dual_fuel_single_rate"`
	DualFuelDualRate      SampleQuoteDirectDebitRate `json:"dual_fuel_dual_rate"`
}

type SampleQuoteDirectDebitRate struct {
	AnnualCostIncVat int `json:"annual_cost_inc_vat"`
	AnnualCostExcVat int `json:"annual_cost_exc_vat"`
}

type SampleConsumption struct {
	ElectricitySingleRate struct {
		ElectricityStandard int `json:"electricity_standard"`
	} `json:"electricity_single_rate"`
	ElectricityDualRate struct {
		ElectricityDay   int `json:"electricity_day"`
		ElectricityNight int `json:"electricity_night"`
	} `json:"electricity_dual_rate"`
	DualFuelSingleRate struct {
		ElectricityStandard int `json:"electricity_standard"`
		GasStandard         int `json:"gas_standard"`
	} `json:"dual_fuel_single_rate"`
	DualFuelDualRate struct {
		ElectricityDay   int `json:"electricity_day"`
		ElectricityNight int `json:"electricity_night"`
		GasStandard      int `json:"gas_standard"`
	} `json:"dual_fuel_dual_rate"`
}

// Get retrieves the details of a product (including all its tariffs) for a particular point in time.
// This endpoint is paginated, it will return next and previous links if returned data is larger than
// the set page size, you are responsible to request the next page if required.
func (s *ProductService) Get(options *ProductsGetOptions) (*ProductsGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context
func (s *ProductService) GetWithContext(ctx context.Context, options *ProductsGetOptions) (*ProductsGetOutput, error) {
	path := fmt.Sprintf("/v1/products/%s", options.ProductCode)
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

	res := ProductsGetOutput{}
	if err := s.client.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
