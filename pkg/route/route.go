package route

import (
	"context"
	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

// Manager is a route worker.
type Manager struct {
	*svc.BaseWorker
}

// NewWorker creates a new route calculator worker instance.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *Manager {
	w := &Manager{
		BaseWorker: svc.NewWorker(ctx, cfg, log, "raddare-route-manager"),
	}
	return w
}
