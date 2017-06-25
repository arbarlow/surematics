package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arbarlow/surematics/postcode/geocoder"
)

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("postcode needs two args: i,e postcode 'SW1A 1AA' 'E8 4AA'")
	}

	p1, p2 := args[0], args[1]
	dist, err := geocoder.DistanceBetweenCodes(p1, p2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Distance between %v, %v: %.2f km\n", p1, p2, dist/1000)
}
