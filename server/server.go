package server

import (
	"errors"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/dbadoy/grinder/server/cft"
	"github.com/dbadoy/grinder/server/fetcher"
)

var (
	ErrServerTooBusy = errors.New("main loop is in a blocked state")
)

type Server struct {
	engine cft.Engine
	eth    ethclient.Client
	cp     checkpoint.CheckpointReader

	fetcher *fetcher.Fetcher

	journals []journalObject

	// main loop
	req  chan request
	quit chan struct{}

	cfg *Config
}

func New(eth ethclient.Client, fetcher *fetcher.Fetcher, engine cft.Engine, cp checkpoint.CheckpointReader, cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &Server{
		engine:   engine,
		eth:      eth,
		fetcher:  fetcher,
		journals: make([]journalObject, 0),
		req:      make(chan request),
		quit:     make(chan struct{}),
		cfg:      cfg,
	}, nil
}

func (s *Server) Run() {
	s.fetcher.Run()
	go s.loop()
}

func (s *Server) Stop() {
	s.fetcher.Stop()
	s.quit <- struct{}{}
	s.quit = make(chan struct{})
}

func (s *Server) EthClient() ethclient.Client {
	return s.eth
}

func (s *Server) Checkpoint() checkpoint.CheckpointReader {
	return s.cp
}
func (s *Server) AddABI(req *ABIRequest) error {
	req.errc = make(chan error)
	select {
	case s.req <- req:
		return <-req.errc
	default:
		return ErrServerTooBusy
	}
}

func (s *Server) AddABISync(req *ABIRequest) error {
	req.errc = make(chan error)
	s.req <- req
	return <-req.errc
}

func (s *Server) AddContract(req *ContractRequest) error {
	req.errc = make(chan error)
	select {
	case s.req <- req:
		return <-req.errc
	default:
		return ErrServerTooBusy
	}
}

func (s *Server) AddContractSync(req *ContractRequest) error {
	req.errc = make(chan error)
	s.req <- req
	return <-req.errc
}

func (s *Server) loop() {
	for {
		select {
		case block := <-s.fetcher.C:
			if s.engine.Checkpoint()+1 == block.NumberU64() {
				if err := s.handleBlock(block); err == nil {
					s.engine.Increase()
				}
			}

		case req := <-s.req:
			s.handleRequest(req)

		case <-s.quit:
			return
		}
	}
}
