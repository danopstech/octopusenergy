package octopusenergy

import (
	"log"
	"net/http"
	"os"
)

// Config provides service configuration for client.
type Config struct {
	// Required for Authentication on non public API end-points when using this client.
	// If you are an Octopus Energy customer, you can generate an API key from your online dashboard
	// https://octopus.energy/dashboard/developer/.
	ApiKey *string

	// All API requests will use this base URL.
	Endpoint *string

	// The HTTP client to use when sending requests. Defaults to `http.DefaultClient`.
	HTTPClient *http.Client
}

// NewConfig returns a new Config pointer that can be chained with builder
// methods to set multiple configuration values inline without using pointers.
//
//     client := octopusenergy.NewClient(octopusenergy.NewConfig().
//         WithApiKey("your-api-key"),
//     ))
func NewConfig() *Config {
	return &Config{}
}

// WithApiKey sets a config ApiKey value returning a Config pointer for chaining.
func (c *Config) WithApiKey(apiKey string) *Config {
	c.ApiKey = &apiKey
	return c
}

// WithApiKeyFromEnvironments sets a config ApiKey value from environments valuable
// returning a Config pointer for chaining.
func (c *Config) WithApiKeyFromEnvironments() *Config {
	apiKey, ok := os.LookupEnv(apiKeyEnvKey)
	if !ok {
		log.Fatalln("could not find api key in environment variable 'OCTOPUS_ENERGY_API_KEY'")
	}
	if apiKey == "" {
		log.Fatalln("the api key in environment variable 'OCTOPUS_ENERGY_API_KEY' is blank")
	}

	c.ApiKey = &apiKey
	return c
}

// WithEndpoint sets a config Endpoint value returning a Config pointer for chaining.
func (c *Config) WithEndpoint(endpoint string) *Config {
	c.Endpoint = &endpoint
	return c
}

// WithHTTPClient sets a config HTTP Client value returning a Config pointer for chaining.
func (c *Config) WithHTTPClient(HTTPClient http.Client) *Config {
	c.HTTPClient = &HTTPClient
	return c
}
