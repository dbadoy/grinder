package api

import (
	"fmt"
	"net/http"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient"
)

// server.Server
type Backend interface {
	EthClient() ethclient.Client
	Checkpoint() checkpoint.CheckpointReader
}

type service interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	path() string
	get(w http.ResponseWriter, r *http.Request)
	post(w http.ResponseWriter, r *http.Request)
	put(w http.ResponseWriter, r *http.Request)
	delete(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	port uint16
	b    Backend
}

func New(port uint16, b Backend) *Server {
	return &Server{port, b}
}

func (s *Server) Listen() {
	for _, api := range SupportAPIs(s.b) {
		http.Handle(api.path(), api)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}

func SupportAPIs(b Backend) []service {
	return []service{
		&status{b},
	}
}
