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

	b, err := json.MarshalIndent(consumption, "", "  ")
	if err != nil {
		log.Fatalf("failed to pretty print: %s", err.Error())
	}
	fmt.Print(string(b))
}
