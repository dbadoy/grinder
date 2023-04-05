package fetcher

import (
	"time"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

type Config struct {
	PollInterval time.Duration
}

type Fetcher struct {
	eth ethclient.Client
	cp  checkpoint.CheckpointReader

	C    chan *types.Block
	quit chan struct{}

	cfg *Config
}

func New(client ethclient.Client, cp checkpoint.CheckpointReader, cfg *Config) *Fetcher {
	// redgla
	//
	return &Fetcher{
		eth:  client,
		cp:   cp,
		C:    make(chan *types.Block),
		quit: make(chan struct{}, 1),
		cfg:  cfg,
	}
}

func (f *Fetcher) Run() {
}

func (f *Fetcher) Stop() {
	f.quit <- struct{}{}
}
