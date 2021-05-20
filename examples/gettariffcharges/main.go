package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danopstech/octopusenergy"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var netClient = http.Client{
		Timeout: time.Second * 10,
	}

	client := octopusenergy.NewClient(octopusenergy.NewConfig().
		WithHTTPClient(netClient),
	)

	consumption, err := client.TariffCharge.GetWithContext(ctx, &octopusenergy.TariffChargesGetOptions{
		ProductCode: "AGILE-18-02-21",
		TariffCode:  "E-1R-AGILE-18-02-21-H",
		FuelType:    octopusenergy.FuelTypeElectricity,
		Rate:        octopusenergy.RateStandardUnit,
	})

	if err != nil {
		log.Fatalf("failed to getting consumption: %s", err.Error())
	}

	b, err := json.MarshalIndent(consumption, "", "  ")
	if err != nil {
		log.Fatalf("failed to pretty print: %s", err.Error())
	}
	fmt.Print(string(b))
}
