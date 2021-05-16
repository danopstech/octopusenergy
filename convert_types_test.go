package octopusenergy_test

import (
	"fmt"
	"time"

	"github.com/danopstech/octopusenergy"
)

func ExampleBool() {
	boolPtr := octopusenergy.Bool(true)
	fmt.Println(*boolPtr)
	// Output: true
}

func ExampleString() {
	stringPtr := octopusenergy.String("example")
	fmt.Println(*stringPtr)
	// Output: example
}

func ExampleInt() {
	intPtr := octopusenergy.Int(5)
	fmt.Println(*intPtr)
	// Output: 5
}

func ExampleTime() {
	wayBack := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	timePtr := octopusenergy.Time(wayBack)
	fmt.Println(*timePtr)
	// Output: 1974-05-19 01:02:03.000000004 +0000 UTC
}
