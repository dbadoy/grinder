package fetcher

import (
	"os"
	"testing"
	"time"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/ethclient/mock"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestSubscribeFetcherRealTime(t *testing.T) {
	c, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	c.SupportSubscribe = true

	cp := checkpoint.New(checkpoint.DefaultBasePath, "fetcher")
	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	fetcher := New(c, cp, &Config{PollInterval: 50 * time.Millisecond})

	mined := make([]*types.Block, 0)
	go func() {
		for {
			block := <-fetcher.C
			if block.NumberU64() == fetcher.cp.Checkpoint()+1 {
				mined = append(mined, block)
				cp.Increase()
			}
		}
	}()

	fetcher.Run()

	for i := 1; i <= 10; i++ {
		c.Core().Commit()
		time.Sleep(100 * time.Millisecond)

		if len(mined) != i {
			t.Fatalf("TestSubscribeFetcherRealTime, mined: %d want: %d", len(mined), i)
		}
	}
}

func TestSubscribeFetcherRecover(t *testing.T) {
	c, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	c.SupportSubscribe = true

	cp := checkpoint.New(checkpoint.DefaultBasePath, "fetcher")
	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	fetcher := New(c, cp, &Config{PollInterval: 50 * time.Millisecond})

	mined := make([]*types.Block, 0)
	go func() {
		for {
			block := <-fetcher.C
			if block.NumberU64() == fetcher.cp.Checkpoint()+1 {
				mined = append(mined, block)
				cp.Increase()
			}
		}
	}()

	want := 10

	for i := 1; i <= want; i++ {
		c.Core().Commit()
	}

	fetcher.Run()

	// If the fetcher starts, and the next block is not created,
	// recover is not automatically performed.
	c.Core().Commit()

	time.Sleep(1 * time.Second)

	if len(mined) != want+1 /* Include one block for recovery triggers */ {
		t.Fatalf("TestSubscribeFetcherRecover, mined: %d want: %d", len(mined), want)
	}
}

func TestPollingFetcherRealTime(t *testing.T) {
	c, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	cp := checkpoint.New(checkpoint.DefaultBasePath, "fetcher")
	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	fetcher := New(c, cp, &Config{PollInterval: 50 * time.Millisecond})

	mined := make([]*types.Block, 0)
	go func() {
		for {
			block := <-fetcher.C
			if block.NumberU64() == fetcher.cp.Checkpoint()+1 {
				mined = append(mined, block)
				cp.Increase()
			}
		}
	}()

	fetcher.Run()

	for i := 1; i <= 10; i++ {
		c.Core().Commit()
		time.Sleep(100 * time.Millisecond)

		if len(mined) != i {
			t.Fatalf("TestFetcherRealTime, mined: %d want: %d", len(mined), i)
		}
	}
}

func TestPollingFetcherRecover(t *testing.T) {
	c, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	cp := checkpoint.New(checkpoint.DefaultBasePath, "fetcher")
	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	fetcher := New(c, cp, &Config{PollInterval: 50 * time.Millisecond})

	mined := make([]*types.Block, 0)
	go func() {
		for {
			block := <-fetcher.C
			if block.NumberU64() == fetcher.cp.Checkpoint()+1 {
				mined = append(mined, block)
				cp.Increase()
			}
		}
	}()

	want := 10

	for i := 1; i <= want; i++ {
		c.Core().Commit()
	}

	fetcher.Run()

	time.Sleep(1 * time.Second)

	if len(mined) != want {
		t.Fatalf("TestFetcherRecover, mined: %d want: %d", len(mined), want)
	}
}
