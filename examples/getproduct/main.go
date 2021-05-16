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

	product, err := client.Product.GetWithContext(ctx, &octopusenergy.ProductsGetOptions{
		ProductCode: octopusenergy.ProductCodeAgile180221,
	})

	if err != nil {
		log.Fatalf("failed to get product: %s", err.Error())
	}

	b, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		log.Fatalf("failed to pretty print: %s", err.Error())
	}
	fmt.Print(string(b))
}
