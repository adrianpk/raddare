package route

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/raddare/internal/osrm"
)

func (m *Manager) getRoutesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wps, ok := ctx.Value(waypointsCtxKey).(*Waypoints)
	if !ok {
		err := errors.New("incomplete coordinates data")
		m.errorResponse(w, r, err)
	}

	osrm, err := m.osrmHandler()
	if err != nil {
		m.errorResponse(w, r, err)
	}

	routes, err := osrm.Routes(wps.All())
	if err != nil {
		m.errorResponse(w, r, err)
	}

	out := fmt.Sprintf("getRoutesHandler:\n\n%+v", routes)
	w.Write([]byte(out))
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
