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

// OSRMRequest maps request data sent OSRM.
// Ref.: http://project-osrm.org/docs/v5.5.1/api/#nearest-service
type (
	OSRMRequest struct {
		OSRMRoutes    []OSRMRoutes    `json:"routes"`
		OSRMWaypoints []OSRMWaypoints `json:"waypoints"`
		OSRMCode      string          `json:"code"`
	}

	OSRMLegs struct {
		Summary  string        `json:"summary"`
		Weight   int           `json:"weight"`
		Duration float64       `json:"duration"`
		Steps    []interface{} `json:"steps"`
		Distance float64       `json:"distance"`
	}

	OSRMRoutes struct {
		Legs       []OSRMLegs `json:"legs"`
		WeightName string     `json:"weight_name"`
		Weight     int        `json:"weight"`
		Duration   float64    `json:"duration"`
		Distance   float64    `json:"distance"`
	}

	OSRMWaypoints struct {
		Hint     string    `json:"hint"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	}
)
