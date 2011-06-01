package trips

import (
	"container/vector"
	"container/heap"
	"flights"
	"fmt"
	"io"
)

type Trip struct {
	From           string
	To             string
	BeganAt        uint
	CurrentAirport string
	CurrentTime    uint
	TotalTime      uint
	TotalCost      float32
	Flights        *vector.Vector
}

func (t *Trip) Done() bool {
	return t.CurrentAirport == t.To
}

func (t *Trip) Print(out io.Writer) {
	t.Flights.Do(func(v interface{}) {
		f := v.(*flights.Flight)
		fmt.Fprint(out, "-- ")
		f.Print(out)
		fmt.Fprintln(out)
	})
}

func NewTrip(from, to string) *Trip {
	result := new(Trip)
	result.From = from
	result.To = to
	result.CurrentAirport = from
	result.Flights = new(vector.Vector)
	return result
}

type TripHeap struct {
	vector.Vector
	seenSet map[string]bool
	less    func(x, y interface{}) bool
}

func (h *TripHeap) Less(i, j int) bool {
	return h.less(h.At(i), h.At(j))
}

func NewTripHeap(from, to string, lessFunc func(x, y interface{}) bool) *TripHeap {
	result := new(TripHeap)
	result.seenSet = make(map[string]bool)
	result.less = lessFunc
	result.Push(NewTrip(from, to))
	heap.Init(result)
	return result
}

func (h *TripHeap) Process(flightData *flights.FlightData) {
	top := heap.Pop(h).(*Trip)
	h.seenSet[top.CurrentAirport] = true
	departures := flightData.Departures[top.CurrentAirport]
	if departures == nil {
		// panic("we got a nil departures for " + top.CurrentAirport + " somehow")
		return
	}
	for flightElem := departures.Front(); flightElem != nil; flightElem = flightElem.Next() {
		flight := flightElem.Value.(flights.Flight)
		if flight.Depart >= top.CurrentTime && !h.seenSet[flight.To] {
			newFlights := top.Flights.Copy()
			newFlights.Push(&flight)
			t := Trip{top.From, top.To, 0, flight.To, flight.Arrive, 0, top.TotalCost + flight.Cost, &newFlights}
			if top.TotalTime == 0 {
				t.BeganAt = flight.Depart
				t.TotalTime = flight.Arrive - flight.Depart
			} else {
				t.BeganAt = top.BeganAt
				t.TotalTime = top.TotalTime + (flight.Arrive - top.CurrentTime)
			}
			heap.Push(h, &t)
		}
	}
}

func FindOptimal(from, to string, flightSchedule *flights.FlightData, lessFunc func(x, y interface{}) bool) *Trip {
	h := NewTripHeap(from, to, lessFunc)
	for !h.Done() {
		h.Process(flightSchedule)
	}

	if h.Failed() {
		return nil
	}
	
	trip := heap.Pop(h).(*Trip)
	return trip
}

func (h *TripHeap) Done() bool {
	return h.Len() == 0 || h.At(0).(*Trip).Done()
}

func (h *TripHeap) Failed() bool {
	return h.Len() == 0
}

func LessCost(x, y interface{}) bool {
	t1 := x.(*Trip)
	t2 := y.(*Trip)
	switch {
	case t1.TotalCost < t2.TotalCost:
		return true
	case t1.TotalCost == t2.TotalCost:
		return t1.TotalTime < t2.TotalTime
	}
	return false
}

func LessTime(x, y interface{}) bool {
	t1 := x.(*Trip)
	t2 := y.(*Trip)
	switch {
	case t1.TotalTime < t2.TotalTime:
		return true
	case t1.TotalTime == t2.TotalTime:
		return t1.TotalCost < t2.TotalCost
	}
	return false
}

func (h *TripHeap) Print() {
	for i := 0; i < h.Len(); i++ {
		t := h.At(i).(*Trip)
		fmt.Println(i, t.From, t.To, t.CurrentAirport, t.CurrentTime, t.TotalTime, t.TotalCost)
	}
}
