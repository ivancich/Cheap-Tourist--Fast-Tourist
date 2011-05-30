package flights

import (
	"fmt"
	"strconv"
	"io"
	"strings"
	"container/list"
)

type Flight struct {
	From   string
	To     string
	Depart uint
	Arrive uint
	Cost   float32
}

type FlightData struct {
	Flights    []*Flight
	Departures map[string]*list.List
}

func (fd *FlightData) GetDeparturesFrom(airport string) *list.List {
	return fd.Departures[airport]
}

func PrintDepartures(airport string, data *FlightData) {
	list := data.Departures[airport]
	for cursor := list.Front() ; cursor != nil ; cursor = cursor.Next() {
		f := cursor.Value.(Flight)
		fmt.Println(f.From, f.To, f.Cost)
	}
}

func MakeFlightSchedule(in io.Reader, flightCount uint) *FlightData {
	result := new(FlightData)

	result.Flights = make([]*Flight, flightCount)
	result.Departures = make(map[string] *list.List)

	var from, to, departure, arrival string
	var cost float32

	for i := uint(0); i < flightCount; i++ {
		fmt.Fscanln(in, &from, &to, &departure, &arrival, &cost)
		// fmt.Println(from, to, parseTime(departure), parseTime(arrival), cost)
		flight := Flight{from, to, parseTime(departure), parseTime(arrival), cost}
		result.Flights[i] = &flight
		if (result.Departures[from] == nil) {
			result.Departures[from] = list.New()
		}
		result.Departures[from].PushBack(flight)
	}

	return result
}

func parseTime(timeStr string) uint {
	pieces := strings.Split(timeStr, ":", 2)
	hours, _ := strconv.Atoui(pieces[0])
	minutes, _ := strconv.Atoui(pieces[1])
	return hours*uint(60) + minutes
}

