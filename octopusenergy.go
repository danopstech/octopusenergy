//go:generate stringer -linecomment -output=gen_stringer.go -type=FuelType,Rate

// Package octopusenergy proves an interface for Octopus Energy REST APIs.
package octopusenergy

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const (
	// ProductCodeAgile180221 is the product code to Octopus current Agile tariff
	ProductCodeAgile180221 = "AGILE-18-02-21"

	defaultBaseURL = "https://api.octopus.energy/"
	userAgent      = "octopus-energy-api-client-go/0.0.0"
	apiKeyEnvKey   = "OCTOPUS_ENERGY_API_KEY"
)

type FuelType int

const (
	FuelTypeElectricity FuelType = iota // electricity
	FuelTypeGas                         // gas
)

type service struct {
	client *Client
}

type Client struct {
	// Base URL for API requests. Defaults to the public Octopus API.
	BaseURL url.URL

	//
	auth string

	// User agent used when communicating with the GitHub API.
	userAgent string

	// HTTP client used to communicate with the API.
	HTTPClient *http.Client

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Octopus API.
	TariffCharge    *TariffChargeService
	MeterPoint      *MeterPointService
	Product         *ProductService
	GridSupplyPoint *GridSupplyPointService
	Consumption     *ConsumptionService
}

// NewClient accepts a config object and returns an initiated client ready to use.
func NewClient(cfg *Config) *Client {
	url, _ := url.Parse(defaultBaseURL)
	httpClient := http.DefaultClient
	var auth string

	if cfg.Endpoint != nil {
		url, _ = url.Parse(*cfg.Endpoint)
	}

	if cfg.HTTPClient != nil {
		httpClient = cfg.HTTPClient
	}

	if cfg.ApiKey != nil {
		auth = base64.StdEncoding.EncodeToString([]byte(*cfg.ApiKey + ":"))
	}

	c := &Client{
		BaseURL:    *url,
		auth:       auth,
		userAgent:  userAgent,
		HTTPClient: httpClient,
	}

	c.common.client = c
	c.TariffCharge = (*TariffChargeService)(&c.common)
	c.MeterPoint = (*MeterPointService)(&c.common)
	c.Product = (*ProductService)(&c.common)
	c.GridSupplyPoint = (*GridSupplyPointService)(&c.common)
	c.Consumption = (*ConsumptionService)(&c.common)

	return c
}

func addParameters(url *url.URL, parameters interface{}) (*url.URL, error) {
	v := reflect.ValueOf(parameters)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return url, nil
	}

	vs, err := query.Values(parameters)
	if err != nil {
		return url, err
	}
	url.RawQuery = vs.Encode()
	return url, nil
}

func (c *Client) sendRequest(req *http.Request, authed bool, castTo interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", c.userAgent)

	if authed && c.auth != "" {
		req.Header.Set("Authorization", "Basic "+c.auth)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			// TODO: return common custom errors types
			return errors.New(errRes.Detail)
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&castTo); err != nil {
		return err
	}

	return nil
}
