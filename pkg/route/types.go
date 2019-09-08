package route

import ()

// RoutingReq  maps service request.
type (
	RoutingReq struct {
		Source string `json:"source"`
	}

	Point struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	ReqParams struct {
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
