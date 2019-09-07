package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/raddare/internal/osrm"
	"github.com/raddare/pkg/route"
	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

type contextKey string

var (
	s svc.Service
)

func main() {
	cfg := config.Load("grn")
	log := initLog(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go checkSigTerm(cancel)

	// Create service
	s = newService(ctx, cfg, log, cancel)

	// Add service handlers
	osrm, err := osrm.NewHandler(ctx, cfg, log, "osrm-handler")
	s.AddHandler(osrm)

	// Set service worker
	auth := route.NewWorker(ctx, cfg, log, "route-worker")
	s.SetWorker(auth)

	// Initialize handlers and workers
	err = s.Init()
	if err != nil {
		log.Error(err)
		cancel()
	}

	// Start service
	s.Start()

	log.Error(err, "Service stoped")
}

// newService creates a service instance.
func newService(ctx context.Context, cfg *config.Config, log *log.Logger, cancel context.CancelFunc) svc.Service {
	sn := cfg.ValOrDef("svc.name", "raddare")
	sr := cfg.ValOrDef("svc.revision", "n/a")
	s := svc.NewService(ctx, cfg, log, cancel, sn, sr)
	return s
}

func initLog(cfg *config.Config) *log.Logger {
	ll := int(cfg.ValAsInt("log.level", 1))
	sn := cfg.ValOrDef("svc.name", "raddare")
	sr := cfg.ValOrDef("svc.revision", "n/a")
	return log.NewLogger(ll, sn, sr)
}

// checkSigTerm - Listens to sigterm events.
func checkSigTerm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	s.Stop()
	cancel()
	os.Exit(0)
}
