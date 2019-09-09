package route

import (
	"fmt"

	"github.com/raddare/internal/osrm"
)

// FIX: Improve naming consitency
// i.e.: Src/Source, singular name for types, etc...

// RoutingReq  maps service request.
type (
	RoutingReq struct {
		Source string `json:"source"`
	}

	Point struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	Waypoints struct {
		Src [2]float64
		Dst [][2]float64
	}
)

// RoutingRes maps service response.
type (
	RoutingRes struct {
		Source string  `json:"source"`
		Routes []Route `json:"routes"`
	}

	Route struct {
		Destination string  `json:"destination"`
		Duration    float64 `json:"duration"`
		Distance    float64 `json:"distance"`
	}
)

type (
	channeledResponse struct {
		res *osrm.Response
		err error
	}
)

// Return source coordinates as a string.
func (w *Waypoints) SrcAsString() (coords string, ok bool) {
	if len(w.Src) != 2 {
		return "0.0,0.0", false
	}

	return fmt.Sprintf("%f,%f", w.Src[0], w.Src[1]), true
}

// Combinations consolidate origin and destination
// coordinates in a single slice.
func (w *Waypoints) Combinations() [][][2]float64 {
	res := make([][][2]float64, 0)
	for _, d := range w.Dst {
		comb := [][2]float64{}
		comb = append(comb, w.Src)
		comb = append(comb, d)
		res = append(res, comb)
	}
	return res
}

// Combinations consolidate origin and destination
// coordinates in a single slice.
func (w *Waypoints) Combinations2() [][][2]float64 {
	res := make([][][2]float64, 0)
	comb := [][2]float64{}
	for _, d := range w.Dst {
		comb = append(comb, w.Src)
		comb = append(comb, d)
		res = append(res, comb)
	}
	return res
}

// All consolidate origin and destination
// coordinates in a single slice.
func (w *Waypoints) All() [][2]float64 {
	res := make([][2]float64, 0)
	res = append(res, w.Src)
	res = append(res, w.Dst...)
	return res
}
