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

	products, err := client.Product.ListWithContext(ctx, &octopusenergy.ProductsListOptions{})

	if err != nil {
		log.Fatalf("failed to list products: %s", err.Error())
	}

	b, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		log.Fatalf("failed to pretty print: %s", err.Error())
	}
	fmt.Print(string(b))
}
