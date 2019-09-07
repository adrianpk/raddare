package route

import (
	"fmt"
	"net/http"
)

func (m *Manager) getRoutesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getRoutesHandler"))
}

func (m *Manager) errorResponse(w http.ResponseWriter, r *http.Request, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	str := `{"data":"","error":"Invalid request"}`
	fmt.Fprint(w, str)
	m.Log().Error(err, message)
}
