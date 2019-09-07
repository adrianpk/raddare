package osrm

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

type (
	// Handler is an OSRM handler.
	Handler struct {
		*svc.BaseHandler
	}
)

// NewHandler creates and returns a new repo handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger) (*Handler, error) {
	name := fmt.Sprintf("osrm-handler-%s", nameSufix())
	log.Info("New handler", "name", name)

	h := &Handler{
		BaseHandler: svc.NewBaseHandler(ctx, cfg, log, name),
	}

	return h, nil
}

// Init a new OSRM handler.
func (h *Handler) Init(s svc.Service) chan bool {
	ok := make(chan bool)
	go func() {
		defer close(ok)
		s.Lock()
		s.AddHandler(h)
		s.Unlock()
		h.Log().Info("Handler initializated", "name", h.Name())
		ok <- true
	}()
	return ok
}

func nameSufix() string {
	digest := hash(time.Now().String())
	return digest[len(digest)-8:]
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}
