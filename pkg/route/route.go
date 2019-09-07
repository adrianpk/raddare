package route

import (
	"context"
	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

type (
	contextKey string
)

// Manager is a route worker.
type Manager struct {
	*svc.BaseWorker
}

const (
	routesCtxKey contextKey = "routes"
)

// NewWorker creates a new route calculator worker instance.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *Manager {
	w := &Manager{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "raddare-route-manager"),
	}
	return w
}

// Init service worker.
func (m *Manager) Init() bool {
	err := m.initServer()
	return err != nil
}
