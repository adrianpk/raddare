package route

import (
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
	}

	// Get OSRM handler.
	oh, err := m.osrmHandler()
	if err != nil {
		m.errorResponse(w, r, err)
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

	// Choose best detination
	best := m.bestRoute(responses)

	// Output result.
	out := fmt.Sprintf("getRoutesHandler:\n\nBest route:\n\n%+v", best)
	w.Write([]byte(out))
}

func (m *Manager) osrmRequest(oh *osrm.Handler, points [][2]float64, ch chan<- channeledResponse) {
	chRes := channeledResponse{}
	res, err := oh.Routes(points)
	if err != nil {
		chRes.err = err
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

func (m *Manager) respDump(resps []*osrm.Response) string {
	var sb strings.Builder
	for _, r := range resps {
		sb.WriteString(r.Code)
		sb.WriteString(" ")
	}
	return sb.String()
}
