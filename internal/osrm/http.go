package osrm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

// newClient creates a custom HTTP client.
func newClient(url string) client {
	var c = &http.Client{
		Timeout: time.Second * 10,
	}
	return client{c, url}
}

// newRoutesRequest creates a request.
func (h Handler) newRoutesRequest(points []Point) *Request {
	return &Request{
		Service: "route",
		Version: h.Cfg().ValOrDef("osrm.api.ver", "v1"),
		Profile: "driving",
		Points:  points,
		Options: Options(map[string]string{"overview": "false"}),
	}
}

// Lat returns point latitude value.
func (p Point) Lat() float64 {
	return p[0]
}

// Lng returns point longitude value.
func (p Point) Lng() float64 {
	return p[1]
}

// toPoints convert an slice of float64 two elements arrays
// into a slice of Point
func toPoints(points [][2]float64) []Point {
	ps := make([]Point, 0)
	for _, p := range points {
		ps = append(ps, Point{p[0], p[1]})
	}
	return ps
}

// ToQueryParams returns a string representation of
// a slice of points formated as a query string
// to be use in OSRM queries.
func toQueryParams(ps []Point) string {
	var qs strings.Builder

	last := len(ps) - 1
	for i, p := range ps {
		qs.WriteString(fmt.Sprintf("%f,%f", p.Lat(), p.Lng()))
		if i < last {
			qs.WriteString(";")
		}
	}

	return qs.String()
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
		hostURL,                 // host
		r.Service,               // service
		r.Version,               // version
		r.Profile,               // profile
		toQueryParams(r.Points), // coords
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
	fmt.Printf("\n%s\n", string(bytes)) // TODO: Remove this line
	err = json.Unmarshal(bytes, res)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error()) // TODO: Remove this line
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

// EncodeOptions into a string.
func (r *Request) EncodeOptions() string {
	opts := r.Options
	var qs strings.Builder

	i := 0
	last := len(opts) - 1
	for k, v := range opts {
		qs.WriteString(fmt.Sprintf("%s=%s", k, v))
		if i < last {
			qs.WriteString(";")
		}
		i++
	}

	return qs.String()
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

func closeBody(c io.Closer) {
	_ = c.Close()
}
