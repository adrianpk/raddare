package route

import "github.com/raddare/internal/osrm"

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
		Source string   `json:"source"`
		Routes []Routes `json:"routes"`
	}

	Routes struct {
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
