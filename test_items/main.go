package main

import (
	"fmt"
	"github.com/hadnu/onatar/internal/etl5e"
	"github.com/hadnu/onatar/internal/etl5e/filter"
)

func main() {
	data, err := etl5e.LoadFrom5eTools("vendor/5etools-src/data", filter.Filter2024Only)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Items parsed: %d\n", len(data.Items))
	for i, item := range data.Items {
		if i < 5 {
			fmt.Printf("  %s (%s) - %s\n", item.Name, item.Type, item.Rarity)
		}
	}
}
