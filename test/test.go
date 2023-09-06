package main

import (
	"fmt"
	"strconv"
)

func main() {
	var max float64
	max, err := strconv.ParseFloat("0.5", 32)
	if err != nil {
		max = 4
	}
	fmt.Printf("%v", max)
}
