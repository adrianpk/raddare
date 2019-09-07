package route

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Sample URL: http://your-service/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219

// startServer starts the server.
func (m *Manager) startServer() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Regio is running!"))
	})

	r.Route("/routes", func(r chi.Router) {
		r.Post("/", m.getRoutesHandler) // POST /routes
	})

	err := http.ListenAndServe(":8080", r)
	fmt.Println(err)
	return err
}

func (m *Manager) routesCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		someParameter := chi.URLParam(r, "someParameter")
		ctx := context.WithValue(r.Context(), routesCtxKey, someParameter)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
