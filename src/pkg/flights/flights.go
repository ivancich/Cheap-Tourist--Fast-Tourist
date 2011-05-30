package flights

import (
	"scanner"
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

func scanTime(s *scanner.Scanner) uint {
	s.Scan()
	hour, _ := strconv.Atoui(s.TokenText())
	s.Scan()
	s.Scan()
	minute, _ := strconv.Atoui(s.TokenText())
	return hour*uint(60) + minute
}


func MakeFlightData(s *scanner.Scanner, flightCount uint) *FlightData {
	result := new(FlightData)

	oldMode := s.Mode
	s.Mode = scanner.ScanInts | scanner.ScanStrings | scanner.ScanFloats

	for i := uint(0); i < flightCount; i++ {
		s.Scan()
		from := s.TokenText()
		s.Scan()
		to := s.TokenText()
		departure := scanTime(s)
		arrival := scanTime(s)
		_ = s.Scan()
		cost, _ := strconv.Atof32(s.TokenText())
		fmt.Println(from, to, departure, arrival, cost)
	}

	s.Mode = oldMode
	return result
}

func parseTime(timeStr string) uint {
	pieces := strings.Split(timeStr, ":", 2)
	hours, _ := strconv.Atoui(pieces[0])
	minutes, _ := strconv.Atoui(pieces[1])
	return hours*uint(60) + minutes
}

func PrintDepartures(airport string, data *FlightData) {
	list := data.Departures[airport]
	for cursor := list.Front() ; cursor != nil ; cursor = cursor.Next() {
		f := cursor.Value.(Flight)
		fmt.Println(f.From, f.To, f.Cost)
	}
}

func MakeFlightData2(in io.Reader, flightCount uint) *FlightData {
	result := new(FlightData)

	result.Flights = make([]*Flight, flightCount)
	result.Departures = make(map[string] *list.List)

	var from, to, departure, arrival string
	var cost float32

	for i := uint(0); i < flightCount; i++ {
		fmt.Fscanln(in, &from, &to, &departure, &arrival, &cost)
		fmt.Println(from, to, parseTime(departure), parseTime(arrival), cost)
		flight := Flight{from, to, parseTime(departure), parseTime(arrival), cost}
		result.Flights[i] = &flight
		if (result.Departures[from] == nil) {
			result.Departures[from] = list.New()
		}
		result.Departures[from].PushBack(flight)
	}

	return result
}
