package server

import (
	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/dbadoy/grinder/server/cft"
	"github.com/dbadoy/grinder/server/fetcher"
)

type Server struct {
	engine cft.Engine
	eth    ethclient.Client
	cp     checkpoint.CheckpointReader

	fetcher *fetcher.Fetcher

	// main loop
	quit chan struct{}

	cfg *Config
}

func New(eth ethclient.Client, fetcher *fetcher.Fetcher, engine cft.Engine, cp checkpoint.CheckpointReader, cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &Server{
		engine:  engine,
		eth:     eth,
		fetcher: fetcher,
		quit:    make(chan struct{}, 1),
		cfg:     cfg,
	}, nil
}

func (s *Server) Run() {
	s.fetcher.Run()
	go s.loop()
}

func (s *Server) Stop() {
	s.fetcher.Stop()
	s.quit <- struct{}{}
}

func (s *Server) EthClient() ethclient.Client {
	return s.eth
}

func (s *Server) Checkpoint() checkpoint.CheckpointReader {
	return s.cp
}

func (s *Server) loop() {}
