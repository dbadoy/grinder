package server

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/database/memdb"
	"github.com/dbadoy/grinder/pkg/ethclient/mock"
	"github.com/dbadoy/grinder/server/cft"
	"github.com/dbadoy/grinder/server/dto"
	"github.com/dbadoy/grinder/server/fetcher"
	"github.com/ethereum/go-ethereum/common"
)

func TestAddABI(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	client.SupportSubscribe = false

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 1 * time.Second})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	input := &dto.ABI{
		MethodIDs: []string{"0x000001", "0x000002"},
		EventIDs:  []string{"0x00000a", "0x00000b"},
	}

	// It must fail because it's before the loop starts.
	if err := s.AddABI(&ABIRequest{Name: "test", ABI: input}); err == nil {
		t.Fatal("TestAddABI, want: failed got: success")
	}

	s.Run()
	defer s.Stop()

	time.Sleep(10 * time.Millisecond)

	if err := s.AddABI(&ABIRequest{Name: "test", ABI: input}); err != nil {
		t.Fatal(err)
	}

	s.AddABI(&ABIRequest{Name: "test-1", ABI: input})
	s.AddABI(&ABIRequest{Name: "test-2", ABI: input})
	s.AddABI(&ABIRequest{Name: "test-3", ABI: input})

	abi := memdb.Get([]byte("test")).(*dto.ABI)
	if len(abi.MethodIDs) != 2 || len(abi.EventIDs) != 2 {
		t.Fatalf("TestAddABI, want: (MethodIDs 2 EventIDs 2) got: (MethodIDs %d EventIDs %d)", len(abi.MethodIDs), len(abi.EventIDs))
	}
}

func TestMustAddABI(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	client.SupportSubscribe = false

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 1 * time.Second})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	input := &dto.ABI{
		MethodIDs: []string{"0x000001", "0x000002"},
		EventIDs:  []string{"0x00000a", "0x00000b"},
	}

	go func() {
		if err := s.MustAddABI(&ABIRequest{Name: "test", ABI: input}); err != nil {
			panic(err)
		}

		abi := memdb.Get([]byte("test")).(*dto.ABI)
		if len(abi.MethodIDs) != 2 || len(abi.EventIDs) != 2 {
			panic(fmt.Errorf("TestAddABI, want: (MethodIDs 2 EventIDs 2) got: (MethodIDs %d EventIDs %d)", len(abi.MethodIDs), len(abi.EventIDs)))
		}
	}()

	time.Sleep(100 * time.Millisecond)
	s.Run()
	defer s.Stop()
}

func TestAddContract(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	client.SupportSubscribe = false

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 1 * time.Second})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	// Remix Storage.sol
	ca, err := mock.DeployContract(client, common.Hex2Bytes("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220322c78243e61b783558509c9cc22cb8493dde6925aa5e89a08cdf6e22f279ef164736f6c63430008120033"))

	if err != nil {
		t.Fatal(err)
	}

	// It must fail because it's before the loop starts.
	if err := s.AddContract(&ContractRequest{Address: ca}); err == nil {
		t.Fatal("TestAddContract, want: failed got: success")
	}

	s.Run()
	defer s.Stop()

	time.Sleep(10 * time.Millisecond)

	if err := s.AddContract(&ContractRequest{Address: ca}); err != nil {
		t.Fatal(err)
	}

	contract := memdb.Get([]byte(ca.Hex())).(*dto.Contract)
	if len(contract.Candidates) != 5 {
		t.Fatalf("TestAddContract, want: 5 got: %d", len(contract.Candidates))
	}
}

func TestMustAddContract(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	client.SupportSubscribe = false

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 1 * time.Second})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	// Remix Storage.sol
	ca, err := mock.DeployContract(client, common.Hex2Bytes("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220322c78243e61b783558509c9cc22cb8493dde6925aa5e89a08cdf6e22f279ef164736f6c63430008120033"))

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		if err := s.MustAddContract(&ContractRequest{Address: ca}); err != nil {
			panic(err)
		}

		contract := memdb.Get([]byte(ca.Hex())).(*dto.Contract)
		if len(contract.Candidates) != 5 {
			panic(fmt.Errorf("TestAddContract, want: 5 got: %d", len(contract.Candidates)))
		}
	}()

	time.Sleep(100 * time.Millisecond)
	s.Run()
	defer s.Stop()
}
