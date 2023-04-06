package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

type Config struct {
	// PollInterval is the interval for polling that is performed if the
	// node does not support subscriptions. If it does support
	// subscriptions, this value is not used.
	PollInterval time.Duration
}

// Fetcher detects block generation, compares the block height on
// the blockchain with the checkpoint caller provide, and either
// forwards the block or performs a recovery process.
type Fetcher struct {
	eth ethclient.Client

	// cp only serves to return the last block number received by
	// the caller. Fetcher does not increment the checkpoint
	// directly to know that it has been 'passed', but rather after
	// the caller has received the block.
	//
	// note: See Server.loop() in server/server.go to see how it is
	// handled by the caller.
	cp checkpoint.CheckpointReader

	C    chan *types.Block
	quit chan struct{}

	cfg *Config
}

func New(client ethclient.Client, cp checkpoint.CheckpointReader, cfg *Config) *Fetcher {
	// TODO(dbadoy): Consider scaling out(e.g. redgla) in case we
	// need to perform recovery operations on many blocks.

	return &Fetcher{
		eth:  client,
		cp:   cp,
		C:    make(chan *types.Block),
		quit: make(chan struct{}),
		cfg:  cfg,
	}
}

func (f *Fetcher) Run() {
	go f.subscribe()
}

func (f *Fetcher) Stop() {
	f.quit <- struct{}{}
	f.quit = make(chan struct{})
}

// subscribe subscribes to events for new blocks. If the target
// node does not support subscription, perform the polling method.
func (f *Fetcher) subscribe() {
	ch := make(chan *types.Block)

	// In the subscription method, recovery does not occur until
	// the next block creation event comes in.
	sub, err := f.eth.SubscribeNewBlock(context.Background(), ch)
	if err == nil {
		for {
			select {
			case block := <-ch:
				f.handle(block.NumberU64())

			case <-f.quit:
				sub.Unsubscribe()
				return
			}
		}
	}

	close(ch)
	f.polling(f.cfg.PollInterval)
}

// polling checks the block number at the given interval time.
func (f *Fetcher) polling(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latest, err := f.eth.GetLatestBlockNumber(context.Background())
			if err != nil {
				continue
			}

			f.handle(latest)

		case <-f.quit:
			return
		}
	}
}

// handle compares the checkpoint to the blockchain's latest
// block number and performs a stateful action.
//
// 1. BN == Checkpoint
//		- Synchronized
//
// 2. BN == Checkpoint + 1
//		- Synchronized, new block detected
//
// 3. BN > Checkpoint
// 		- Asynchronous, data recovery between Checkpoint and BN
//
// 4. BN < Checkpoint
//		- Critical error (need to consider reconfiguration)
func (f *Fetcher) handle(latest uint64) {
	if latest == f.cp.Checkpoint() {
		return
	}

	if latest == f.cp.Checkpoint()+1 {
		block, err := f.eth.BlockByNumber(context.Background(), latest)
		if err != nil {
			return
		}
		f.C <- block
		return
	}

	f.recover(latest)
}

// recover fetches and sends the blocks between the checkpoint
// and the given latest block number in order. If the process
// of fetching a particular block number fails, it will retry.
func (f *Fetcher) recover(latest uint64) {
	if latest > f.cp.Checkpoint() {
		for latest != f.cp.Checkpoint() {
			// synchronization process is a blocking operation,
			// so we check each time if the user has ended the
			// polling.
			select {
			case <-f.quit:
				return
			default:
			}

			block, err := f.eth.BlockByNumber(context.Background(), f.cp.Checkpoint()+1)
			if err != nil {
				// Retry on failure until success.
				continue
			}

			f.C <- block
		}
	} else {
		panic(fmt.Errorf("occur critical error, blockchain latest: %d checkpoint: %d", latest, f.cp.Checkpoint()))
	}
}
