<p align=center>
  <img alt=logo src="https://github.com/danopstech/octopusenergy/raw/main/.docs/assets/workswith.png" height=150 />
  <h3 align=center>Octopus Energy Golang API client</h3>
</p>

---
[![PkgGoDev](https://pkg.go.dev/badge/github.com/danopstech/octopusenergy/)](https://pkg.go.dev/github.com/danopstech/octopusenergy/)
[![License](https://img.shields.io/github/license/danopstech/octopusenergy)](/LICENSE)
[![Release](https://img.shields.io/github/release/danopstech/octopusenergy.svg)](https://github.com/danopstech/octopusenergy/releases/latest)
[![tests](https://github.com/danopstech/octopusenergy/actions/workflows/build.yaml/badge.svg)](https://github.com/danopstech/octopusenergy/actions/workflows/build.yaml)

This package provides a Golang client to [Octopus Energy's API](https://developer.octopus.energy/docs/api/). Octopus Energy provides a REST API for customers to interact with our platform. Amongst other things, it provides functionality for:

- Browsing energy products, tariffs and their charges.
- Retrieving details about a UK electricity meter-point.
- Browsing the half-hourly consumption of an electricity or gas meter.
- Determining the grid-supply-point (GSP) for a UK postcode.

If you are an Octopus Energy customer, you can generate an API key from your [online dashboard](https://octopus.energy/dashboard/developer/).

### Authentication
Authentication is required for all API end-points when using this API client. This is performed via [HTTP Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication). This is configured when you instantiate a new client with a config object.
**Warning: Do not share your secret API keys with anyone.**

### Not an Octopus Energy customer?
Please read about the Octopus tariffs and ensure they are right for you, if you think they are, then please use my [referral link](https://share.octopus.energy/dusk-shark-465). (at the time of writing this we will both receive £50 credit)

### Usage
More in the examples folder

```golang
ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
defer cancel()

var netClient = http.Client{
    Timeout: time.Second * 10,
}

client := octopusenergy.NewClient(octopusenergy.NewConfig().
    WithApiKeyFromEnvironments().
    WithHTTPClient(netClient),
)

consumption, err := client.Consumption.GetPagesWithContext(ctx, &octopusenergy.ConsumptionGetOptions{
    MPN:          "1111111111", // <--- replace
    SerialNumber: "1111111111", // <--- replace
    FuelType:     octopusenergy.FuelTypeElectricity,
    PeriodFrom:   octopusenergy.Time(time.Now().Add(-48 * time.Hour)),
})

if err != nil {
    log.Fatalf("failed to getting consumption: %s", err.Error())
}
```

### Links
- [Octopus Energy API Docs](https://developer.octopus.energy/docs/api/)
- [Get API Key](https://octopus.energy/dashboard/developer/)
