package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/raddare/internal/osrm"
)

// TODO: Move busines logic to a service.
// Only for a matter of time all logic
// is code directly in this handler.
func (m *Manager) getRoutesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get points from context.
	wps, ok := ctx.Value(waypointsCtxKey).(*Waypoints)
	if !ok {
		err := errors.New("incomplete coordinates data")
		m.errorResponse(w, r, err)
		return
	}

	// Get OSRM handler.
	oh, err := m.osrmHandler()
	if err != nil {
		m.errorResponse(w, r, err)
		return
	}

	// Call API
	// One reques for each origin-destination
	// waypoint tuple.
	responses := []*osrm.Response{}
	respCh := make(chan channeledResponse)
	combs := wps.Combinations()
	qty := len(combs)

	// OSRM concurrent API calls.
	for _, points := range combs {
		go m.osrmRequest(oh, points, respCh)
	}

	// Collect responses.
	for i := 0; i < qty; i++ {
		ch := <-respCh
		if ch.err != nil {
			m.Log().Error(err)
			continue
		}
		responses = append(responses, ch.res)
	}

	// Sort results.
	m.sortRoutes(responses)

	// Transform output
	json, err := m.genJSONResponse(wps, responses)
	if err != nil {
		m.errorResponse(w, r, err)
		return
	}

	// Output result.
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (m *Manager) genJSONResponse(wp *Waypoints, ors []*osrm.Response) ([]byte, error) {
	rr := m.genRoutingRes(wp, ors)
	return json.Marshal(rr)
}

func (m *Manager) genRoutingRes(wp *Waypoints, ors []*osrm.Response) RoutingRes {
	routes := []Route{}

	src, _ := wp.SrcAsString() // TODO: Check ok

	for _, or := range ors {
		orr := or.Routes[0]       // Routing data
		orw := or.Waypoints[1]    // Destination
		dest, _ := orw.AsString() // TODO: Check ok

		r := Route{
			Destination: dest,
			Duration:    orr.Duration,
			Distance:    orr.Distance,
		}
		routes = append(routes, r)
	}

	return RoutingRes{
		Source: src,
		Routes: routes,
	}
}

func (m *Manager) osrmRequest(oh *osrm.Handler, points [][2]float64, ch chan<- channeledResponse) {
	chRes := channeledResponse{}

	res, err := oh.Routes(points)
	if err != nil {
		chRes.err = err
	}

	// Discard non "Ok" responses
	if res.Code != "Ok" {
		chRes.err = errors.New("non 'Ok' code response")
	}

	chRes.res = res
	ch <- chRes
}

func (m *Manager) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	str := fmt.Sprintf(`{"data":"","error":"%s"}`, err.Error())
	fmt.Fprint(w, str)
	m.Log().Error(err, err.Error())
}

func (m *Manager) osrmHandler() (*osrm.Handler, error) {
	h, ok := m.Handler("osrm-handler")
	if !ok {
		return nil, errors.New("OSRM handler not available")
	}

	osrm, ok := h.(*osrm.Handler)
	if !ok {
		return nil, errors.New("invalidad OSRM handler")
	}

	return osrm, nil
}

func (m *Manager) responsesDump(resps []*osrm.Response) string {
	var sb strings.Builder
	for _, r := range resps {
		sb.WriteString(fmt.Sprintf("%+v\n", r.Routes))
	}
	return sb.String()
}
