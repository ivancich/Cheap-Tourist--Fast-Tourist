package main

import (
	"os"
	"fmt"
	"flights"
	"trips"
)

func output(t *trips.Trip) {
	fmt.Printf("%02d:%02d %02d:%02d %0.2f\n", t.BeganAt/60, t.BeganAt%60, t.CurrentTime/60, t.CurrentTime%60, t.TotalCost)
}

func main() {
	filename := "../../input/sample-input.txt"
	filename = "../../input/input.txt"

	var err os.Error

	in, err := os.Open(filename, os.O_RDONLY, uint32(0))
	if err != nil {
		panic(err.String())
	}
	defer func() {
		in.Close()
	}()

	var testCases uint
	if _, err = fmt.Fscan(in, &testCases); err != nil {
		panic(err)
	}

	for i := uint(0); i < testCases; i++ {
		var flightCount uint
		if _, err = fmt.Fscan(in, &flightCount); err != nil {
			panic(err)
		}

		flightSchedule := flights.MakeFlightSchedule(in, flightCount)

		cheap := trips.FindOptimal("A", "Z", flightSchedule, trips.LessCost)
		output(cheap)

		short := trips.FindOptimal("A", "Z", flightSchedule, trips.LessTime)
		output(short)

		fmt.Println() // blank after each test case
	}
}
