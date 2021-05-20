package octopusenergy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// AccountService handles communication with the tariff related Octopus API.
type AccountService service

// AccountGetOptions is the options for GetTariffCharges.
type AccountGetOptions struct {
	// The octopus account number to be retrieved.
	AccountNumber string `url:"-"`
}

// AccountGetOutput is the returned struct from GetTariffCharges.
type AccountGetOutput struct {
	Number     string `json:"number"`
	Properties []struct {
		ID                     int                      `json:"id"`
		MovedInAt              time.Time                `json:"moved_in_at"`
		MovedOutAt             *time.Time               `json:"moved_out_at"`
		AddressLine1           string                   `json:"address_line_1"`
		AddressLine2           string                   `json:"address_line_2"`
		AddressLine3           string                   `json:"address_line_3"`
		Town                   string                   `json:"town"`
		County                 string                   `json:"county"`
		Postcode               string                   `json:"postcode"`
		ElectricityMeterPoints []ElectricityMeterPoints `json:"electricity_meter_points"`
		GasMeterPoints         []GasMeterPoints         `json:"gas_meter_points"`
	} `json:"properties"`
}

type ElectricityMeterPoints struct {
	MPAN                string `json:"mpan"`
	ProfileClass        int    `json:"profile_class"`
	ConsumptionStandard int    `json:"consumption_standard"`
	Meters              []struct {
		SerialNumber string `json:"serial_number"`
		Registers    []struct {
			Identifier           string `json:"identifier"`
			Rate                 string `json:"rate"`
			IsSettlementRegister bool   `json:"is_settlement_register"`
		} `json:"registers"`
	} `json:"meters"`
	Agreements []struct {
		TariffCode string     `json:"tariff_code"`
		ValidFrom  time.Time  `json:"valid_from"`
		ValidTo    *time.Time `json:"valid_to"`
	} `json:"agreements"`
}

type GasMeterPoints struct {
	MPRN                string `json:"mprn"`
	ProfileClass        int    `json:"profile_class"`
	ConsumptionStandard int    `json:"consumption_standard"`
	Meters              []struct {
		SerialNumber string `json:"serial_number"`
		Registers    []struct {
			Identifier           string `json:"identifier"`
			Rate                 string `json:"rate"`
			IsSettlementRegister bool   `json:"is_settlement_register"`
		} `json:"registers"`
	} `json:"meters"`
	Agreements []struct {
		TariffCode string     `json:"tariff_code"`
		ValidFrom  time.Time  `json:"valid_from"`
		ValidTo    *time.Time `json:"valid_to"`
	} `json:"agreements"`
}

// Get retrieves the details of an account.
func (s *AccountService) Get(options *AccountGetOptions) (*AccountGetOutput, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext same as Get except it takes a Context
func (s *AccountService) GetWithContext(ctx context.Context, options *AccountGetOptions) (*AccountGetOutput, error) {
	path := fmt.Sprintf("/v1/accounts/%s/", options.AccountNumber)
	rel := &url.URL{Path: path}
	url := s.client.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	res := AccountGetOutput{}
	if err := s.client.sendRequest(req, true, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
