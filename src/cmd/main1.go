package main

import (
	"os"
	// "scanner"
	// "strconv"
	"fmt"
	"flights"
	"trips"
)

func display(t *trips.Trip) {
	fmt.Printf("%02d:%02d %02d:%02d %0.2f\n", t.BeganAt/60, t.BeganAt%60, t.CurrentTime/60, t.CurrentTime%60, t.TotalCost)
}

func main() {
	filename := "../../input/sample-input.txt"
	filename = "../../input/input.txt"
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

	// fmt.Println(testCases)

	for i := uint(0); i < testCases; i++ {
		var flightCount uint
		fmt.Fscan(f, &flightCount)
		// fmt.Println(flightCount)

		flightSchedule := flights.MakeFlightSchedule(f, flightCount)
		// flights.PrintDepartures("A", flightData)

		cheap := trips.FindOptimal("A", "Z", flightSchedule, trips.LessCost)
		short := trips.FindOptimal("A", "Z", flightSchedule, trips.LessTime)
		
		display(cheap)
		display(short)
		
		fmt.Println()
	}

	f.Close()
}
