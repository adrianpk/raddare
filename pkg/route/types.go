package route

import "fmt"

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

func (w *Waypoints) All() [][2]float64 {
	res := make([][2]float64, 0)
	res = append(res, w.Src)
	res = append(res, w.Dst...)

	fmt.Printf("\n--%+v--\n", res)

	return res
}
