package api

import (
	"context"
	"fmt"
	"net/http"
)

var (
	_ = service(&status{})
)

type status struct {
	b Backend
}

func (s *status) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.get(w, r)
	case http.MethodPost:
		s.post(w, r)
	case http.MethodPut:
		s.put(w, r)
	case http.MethodDelete:
		s.delete(w, r)
	default:
		setMethodNotAllowed(w)
	}
}

func (status) path() string { return "/status" }

func (s *status) get(w http.ResponseWriter, r *http.Request) {
	latest, err := s.b.EthClient().GetLatestBlockNumber(context.Background())
	if err != nil {
		setInternalServerError(w, []byte(err.Error()))
		return
	}

	cp := s.b.Checkpoint().Checkpoint()

	// description
	if r.URL.Query().Get("v") != "" {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", "blockchian", "checkpoint", "progress")))
	}

	w.Write([]byte(fmt.Sprintf("%d\t\t%d\t\t%.4f / 1.000\n", latest, cp, (float32(cp) / float32(latest)))))
}

func (s *status) post(w http.ResponseWriter, r *http.Request) { setMethodNotAllowed(w) }

func (s *status) put(w http.ResponseWriter, r *http.Request) { setMethodNotAllowed(w) }

func (s *status) delete(w http.ResponseWriter, r *http.Request) { setMethodNotAllowed(w) }
