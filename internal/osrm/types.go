package osrm

import "fmt"

// OSRMRequest maps request data sent OSRM.
// Ref.: http://project-osrm.org/docs/v5.5.1/api/#nearest-service
type (
	// Request is a container for
	// generic OSRM query parameters.
	Request struct {
		Service string
		Version string
		Profile string
		Points  []Point
		Options Options
	}

	// Options is an alias for a map of type [string]string
	// used to store options for OSRM requests.
	Options map[string]string

	// Point is an alias for a two float64 elements array.
	// used to store geo location coordinates.
	Point [2]float64
)

type (
	// Response for OSRM request.
	Response struct {
		Routes    []Route    `json:"routes"`
		Waypoints []Waypoint `json:"waypoints"`
		Code      string     `json:"code"`
	}

	// Legs values.
	Leg struct {
		Summary  string        `json:"summary"`
		Weight   float64       `json:"weight"`
		Duration float64       `json:"duration"`
		Steps    []interface{} `json:"steps"`
		Distance float64       `json:"distance"`
	}

	// Routes values
	Route struct {
		Legs       []Leg   `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     float64 `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	}

	// Waypoints values
	Waypoint struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	}
)

// Return source coordinates as a string.
func (w *Waypoint) AsString() (coords string, ok bool) {
	if len(w.Location) != 2 {
		return "0.0,0.0", false
	}

	return fmt.Sprintf("%f,%f", w.Location[0], w.Location[1]), true
}
