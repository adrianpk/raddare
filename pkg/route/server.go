package route

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Sample URL: http://your-service/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219

// initServer starts the server.
func (m *Manager) initServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Raddare is running!"))
	})

	r.Route("/routes", func(r chi.Router) {
		r.Use(m.routesCtx)
		r.Get("/", m.getRoutesHandler) // POST /routes
	})

	p := m.Cfg().ValOrDef("server.port", ":8080")

	err := http.ListenAndServe(p, r)
	fmt.Println(err)

	return err
}

func (m *Manager) routesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wps := &Waypoints{}

		// Process query string parameters
		keys := r.URL.Query()

		for k, v := range keys {

			// Process source (from)
			if k == "src" {

				latlng := strings.Split(v[0], ",")

				// Ensure source.
				if len(latlng) < 2 {
					continue
				}

				// Get source latitude.
				lat, okLat := toFloat(latlng[0])

				// Get source longitude.
				lng, okLng := toFloat(latlng[1])

				// Update source if tuple is complete
				if okLat && okLng {
					wps.Src = [2]float64{lat, lng}
				}
			}

			// Process destinations (to)
			if k == "dst" {
				// Ensure even number of values (lat & lng tuples)
				if len(v)%2 != 0 {
					v = v[:len(v)-1]
				}

				// Ensure at least one destination
				if len(v) < 2 {
					continue
				}

				dsts := make([][2]float64, 0)

				for i, _ := range v {
					//latlng := strings.FieldsFunc(v[i], split)

					latlng := strings.Split(v[i], ",")

					// Get source latitude.
					lat, okLat := toFloat(latlng[0])

					// Get source longitude.
					lng, okLng := toFloat(latlng[1])

					// Update destinations if tuple is complete
					if okLat && okLng {
						dsts = append(dsts, [2]float64{lat, lng})
					}
					wps.Dst = dsts
				}
			}
		}

		ctx := context.WithValue(r.Context(), waypointsCtxKey, wps)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func toFloat(floatVal string) (res float64, ok bool) {
	res, err := strconv.ParseFloat(floatVal, 64)
	ok = err == nil
	return res, ok
}
