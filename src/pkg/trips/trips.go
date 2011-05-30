package trips

import (
	"container/vector"
	"flights"
	"fmt"
)

type Trip struct {
	From           string
	To             string
	BeganAt        uint
	CurrentAirport string
	CurrentTime    uint
	TotalTime      uint
	TotalCost      float32
}

func (t *Trip) Done() bool {
	fmt.Println("comparing " + t.CurrentAirport + " and " + t.To)
	return t.CurrentAirport == t.To
}

func NewTrip(from, to string) *Trip {
	result := new(Trip)
	result.From = from
	result.To = to
	result.CurrentAirport = from
	return result
}

type TripHeap struct {
	vector.Vector
	less func(x, y interface{}) bool
}

func NewTripHeap(from, to string, lessFunc func (x, y interface{}) bool) *TripHeap {
	result := new(TripHeap)
	result.less = lessFunc
	result.Push(NewTrip(from, to))
	return result
}

func (h *TripHeap) Process(flightData *flights.FlightData) {
	top := h.Pop().(*Trip)
	departures := flightData.Departures[top.CurrentAirport]
	if departures == nil {
		// panic("we got a nil departures for " + top.CurrentAirport + " somehow")
		return
	}
	for flightElem := departures.Front() ; flightElem != nil ; flightElem = flightElem.Next() {
		flight := flightElem.Value.(flights.Flight)
		if (flight.Depart >= top.CurrentTime) {
			t := Trip{top.From, top.To, 0, flight.To, flight.Arrive, 0, top.TotalCost + flight.Cost}
			if top.TotalTime == 0 {
				t.BeganAt = flight.Depart
				t.TotalTime = flight.Arrive - flight.Depart
			} else {
				t.BeganAt = top.BeganAt
				t.TotalTime = top.TotalTime + (flight.Arrive - top.CurrentTime)
			}
			h.Push(&t)
		}
	}
}

func (h *TripHeap) Less(i, j int) bool {
	return h.less(h.At(i), h.At(j))
}

func (h *TripHeap) Done() bool {
	return h.Len() == 0 || h.At(0).(*Trip).Done()
}

func (h *TripHeap) Failed() bool {
	return h.Len() == 0
}

func LessCost(x, y interface{}) bool {
	t1 := x.(Trip)
	t2 := y.(Trip)
	switch {
	case t1.TotalCost < t2.TotalCost:
		return true
	case t1.TotalCost == t2.TotalCost:
		return t1.TotalTime < t2.TotalTime
	}
	return false
}

func LessTime(x, y interface{}) bool {
	t1 := x.(Trip)
	t2 := y.(Trip)
	switch {
	case t1.TotalTime < t2.TotalTime:
		return true
	case t1.TotalTime == t2.TotalTime:
		return t1.TotalCost < t2.TotalCost
	}
	return false
}
