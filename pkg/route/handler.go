package route

import (
	"errors"
	"fmt"
	"net/http"
)

func (m *Manager) getRoutesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wps, ok := ctx.Value(waypointsCtxKey).(*Waypoints)
	if !ok {
		err := errors.New("incomplete coords data")
		m.errorResponse(w, r, err)
	}

	out := fmt.Sprintf("getRoutesHandler:\n\n%+v", wps)

	w.Write([]byte(out))
}

func (m *Manager) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	str := fmt.Sprintf(`{"data":"","error":"%s"}`, err.Error())
	fmt.Fprint(w, str)
	m.Log().Error(err, err.Error())
}
