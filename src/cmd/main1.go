package main

import (
	"os"
	// "scanner"
	// "strconv"
	"fmt"
	"flights"
	"trips"
)

func main() {
	/*
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println(r)
			}
		}()
	*/

	dir, err := os.Getwd()
	fmt.Println(dir)

	filename := "../../input/sample-input.txt"
	// figure out what to do with error
	f, err := os.Open(filename, os.O_RDONLY, uint32(0))
	if err != nil {
		fmt.Println(err.String())
	}

	var testCases uint
	_, err = fmt.Fscan(f, &testCases)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(testCases)

	for i := uint(0); i < testCases; i++ {
		var flightCount uint
		fmt.Fscan(f, &flightCount)
		fmt.Println(flightCount)

		flightData := flights.MakeFlightData2(f, flightCount)
		// flights.PrintDepartures("A", flightData)

		flightHeap := trips.NewTripHeap("A", "Z", trips.LessCost)
		for !flightHeap.Done() {
			flightHeap.Process(flightData)
		}

		if flightHeap.Failed() {
			fmt.Println("could not find trip from A to Z")
		} else {
			trip := flightHeap.At(0).(*trips.Trip)
			fmt.Println(trip.TotalTime, trip.TotalCost)
		}
	}
}
