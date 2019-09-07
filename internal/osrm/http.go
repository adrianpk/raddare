package osrm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type (
	client struct {
		httpClient httpClient
		URL        string
	}

	httpClient interface {
		Do(*http.Request) (*http.Response, error)
	}
)

// OSRMRequest maps request data sent OSRM.
// Ref.: http://project-osrm.org/docs/v5.5.1/api/#nearest-service
type (
	// Request is a container for
	// generic OSRM query parameters.
	Request struct {
		Service string
		Version string
		Profile string
		Points  PointSet
		Options Options
	}

	// Options is an alias for a map of type [string]string
	// used to store options for OSRM requests.
	Options map[string]string

	// PointSet is an alias for a slice of Point.
	// used to store geo location geometries..
	PointSet []Point

	// Point is an alias for a two float64 elements array.
	// used to store geo location coordinates.
	Point [2]float64

	// Response struct is used to unmarshall an OSRM response.
	Response struct {
		Routes    []Route    `json:"routes"`
		Waypoints []Waypoint `json:"waypoints"`
		Code      string     `json:"code"`
	}

	Legs struct {
		Summary  string        `json:"summary"`
		Weight   int           `json:"weight"`
		Duration float64       `json:"duration"`
		Steps    []interface{} `json:"steps"`
		Distance float64       `json:"distance"`
	}

	Route struct {
		Legs       []Legs  `json:"legs"`
		WeightName string  `json:"weight_name"`
		Weight     int     `json:"weight"`
		Duration   float64 `json:"duration"`
		Distance   float64 `json:"distance"`
	}

	Waypoint struct {
		Hint  string  `json:"hint"`
		Name  string  `json:"name"`
		Point []Point `json:"location"`
	}
)

// newClient creates a custom HTTP client.
func newClient(c httpClient, url string) client {
	return client{c, url}
}

// newRoutesRequest creates a request.
func (h Handler) newRoutesRequest(points PointSet) *Request {
	return &Request{
		Service: "route",
		Version: h.Cfg().ValOrDef("osrm.api.ver", "v1"),
		Profile: "driving",
		Points:  points,
		Options: Options(map[string]string{"overview": "false"}),
	}
}

// Lat returns point latitude value.
func (p *Point) Lat() float64 {
	return p[0]
}

// Lng returns point longitude value.
func (p *Point) Lng() float64 {
	return p[1]
}

// URL generates a url for OSRM request
func (r *Request) URL(hostURL string) (reqURL string, err error) {
	if r.Service == "" {
		return "", errors.New("no service declared")
	}

	if r.Version == "" {
		return "", errors.New("no API version declared")
	}

	if r.Profile == "" {
		return "", errors.New("no profile declared")
	}

	if r.CountPoints() == 0 {
		return "", errors.New("no coordinates  declared")
	}

	// http://{host}/{service}/{version}/{profile}/{coords}[.{format}]?option=value&option2=value2
	// i.e.: 'http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407?overview=false'
	reqURL = strings.Join([]string{
		hostURL,   // host
		r.Service, // service
		r.Version, // version
		r.Profile, // profile
		"polyline(" + url.PathEscape(r.EncodePoints()) + ")", // coords
	}, "/")

	if r.CountOptions() > 0 {
		reqURL += "?" + r.EncodeOptions() // options
	}
	return reqURL, nil
}

// MakeRequest sends a request to OSRM.
// and unmarshall JSON response into an appropriate struct.
func (c client) MakeRequest(ctx context.Context, req *Request, res interface{}) error {
	url, err := req.URL(c.URL)
	if err != nil {
		return err
	}

	r, err := c.get(ctx, url)
	if err != nil {
		return err
	}
	defer closeBody(r.Body)

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("cannot read response body")
	}

	// Non body response (non 200 or 400 status from OSRM)
	if r.StatusCode != http.StatusOK && r.StatusCode != http.StatusBadRequest {
		return errors.New("non valid response status")
	}

	// Unmarshall response.
	err = json.Unmarshal(bytes, res)
	if err != nil {
		return fmt.Errorf("response body cannot be unmarshalled")
	}

	return nil
}

func (c client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req.WithContext(ctx))
}

// PolylineFactor value used for encoding.
func (r *Request) PolylineFactor() float64 {
	return 1.0e6
}

// EncodePoints values into a string.
func (r *Request) EncodePoints() string {
	f := r.PolylineFactor()

	var pLat int
	var pLng int

	var result bytes.Buffer
	mark1 := make([]byte, 0, 50)
	mark2 := make([]byte, 0, 50)

	for _, p := range r.Points {
		latCorr := int(math.Floor(p.Lat()*f + 0.5))
		lngCorr := int(math.Floor(p.Lng()*f + 0.5))

		deltaLat := latCorr - pLat
		deltaLng := lngCorr - pLng

		pLat = latCorr
		pLng = lngCorr

		result.Write(append(encodeSignedNumber(deltaLat, mark1), encodeSignedNumber(deltaLng, mark2)...))

		mark1 = mark1[:0]
		mark2 = mark2[:0]
	}

	return result.String()
}

// EncodeOptions into a string.
func (r *Request) EncodeOptions() string {
	opts := r.Options

	if opts == nil {
		return ""
	}

	keys := make([]string, 0, len(opts))
	for k := range opts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf []byte
	for _, k := range keys {
		if len(buf) > 0 {
			buf = append(buf, '&')
		}

		buf = append(buf, url.QueryEscape(k)...)
		buf = append(buf, '=')

		for n, val := range opts[k] {
			if n > 0 {
				buf = append(buf, ';')
			}

			buf = append(buf, url.QueryEscape(string(val))...)
		}
	}
	return string(buf)
	return ""
}

// CountPoints in request.
func (r *Request) CountPoints() int {
	return len(r.Points)
}

// CountOptions in request.
func (r *Request) CountOptions() int {
	return len(r.Options)
}

func (r Response) Error() error {
	if r.Code != codeOK {
		return fmt.Errorf("OSRM API error code: %s", r.Code)
	}
	return nil
}

// toPointSet converts a slice of float64 two elements array
// into a PointSet struct.
func toPointSet(points [][2]float64) PointSet {
	var pset PointSet
	for _, point := range points {
		p := [2]float64{point[0], point[1]}
		pset = append(pset, p)
	}
	return pset
}

func encodeSignedNumber(num int, result []byte) []byte {
	shifted := num << 1

	if num < 0 {
		shifted = ^shifted
	}

	for shifted >= 0x20 {
		result = append(result, byte(0x20|(shifted&0x1f)+63))
		shifted >>= 5
	}

	return append(result, byte(shifted+63))
}

func closeBody(c io.Closer) {
	_ = c.Close()
}
