package osrm

import (
	"context"
	"time"

	"gitlab.com/mikrowezel/config"
	"gitlab.com/mikrowezel/log"
	svc "gitlab.com/mikrowezel/service"
)

type (
	// Handler is an OSRM handler.
	Handler struct {
		*svc.BaseHandler
		Client client
	}
)

// NewHandler creates and returns a new repo handler.
func NewHandler(ctx context.Context, cfg *config.Config, log *log.Logger, name string) (*Handler, error) {
	// name := fmt.Sprintf("osrm-handler")
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

// Routes returns a list of all possible routes from some point
// to another one ordered by driving time and distante.
// Sample URL: 'http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407?overview=false'
func (h *Handler) Routes(points [][2]float64) (*Response, error) {
	var res Response
	ctx, _ := context.WithTimeout(context.Background(), h.reqTimeout())
	req := h.newRoutesRequest(toPointSet(points))

	err := h.query(ctx, req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (h *Handler) query(ctx context.Context, req *Request, res *Response) error {
	err := h.Client.MakeRequest(ctx, req, res)
	if err != nil {
		return err
	}
	return res.Error()
}

func (h *Handler) reqTimeout() time.Duration {
	to := h.Cfg().ValAsInt("osrm.req.timeout.sec", int64(5))
	return time.Duration(to) * time.Second
}
