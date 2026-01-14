package stdio

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"

	"github.com/73NN0/foe-hammer/internal/config/app"
	"github.com/73NN0/foe-hammer/internal/config/domain"
)

type Server struct {
	svc *app.Service
}

func NewServer(svc *app.Service) *Server {
	return &Server{svc: svc}
}

// Serve reads JSON requests from r and writes JSON responses to w, one per line.
func (s *Server) Serve(r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(bufio.NewReader(r))

	bw := bufio.NewWriter(w)
	enc := json.NewEncoder(bw)

	for {
		var req Request
		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				return bw.Flush()
			}
			_ = enc.Encode(Response{ID: req.ID, OK: false, Error: "invalid json: " + err.Error()})
			_ = bw.Flush()
			continue
		}

		resp := s.handle(req)
		_ = enc.Encode(resp)
		_ = bw.Flush()
	}
}

func (s *Server) handle(req Request) Response {
	switch req.Op {
	case "config.list":
		items, err := s.svc.ListConfigs()
		if err != nil {
			return Response{ID: req.ID, OK: false, Error: err.Error()}
		}
		return Response{ID: req.ID, OK: true, Result: items}

	case "config.get":
		var p GetParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return Response{ID: req.ID, OK: false, Error: "bad params: " + err.Error()}
		}
		cfg, err := s.svc.GetConfig(p.ID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return Response{ID: req.ID, OK: false, Error: "not found"}
			}
			return Response{ID: req.ID, OK: false, Error: err.Error()}
		}
		return Response{ID: req.ID, OK: true, Result: cfg}

	default:
		return Response{ID: req.ID, OK: false, Error: "unknown op: " + req.Op}
	}
}
