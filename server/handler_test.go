package server

import (
	"context"
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

func TestHandleContract(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 24 * time.Hour})
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

	txs, err := client.GetTransactionsByNumber(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}

	for _, tx := range txs {
		ca, err := contractAddress(tx)
		if err != nil {
			t.Fatal(err)
		}

		if err := s.handleContract(tx.Hash(), ca); err != nil {
			t.Fatal(err)
		}
	}

	meta := memdb.Get([]byte(ca.Hex())).(*dto.Contract)
	if len(meta.Candidates) != 5 {
		t.Fatalf("TestProcessContract, want: 5 got: %d", len(meta.Candidates))
	}
	if txs[0].Hash().Hex() != meta.TxHash {
		t.Fatalf("TestProcessContract, want: %s got: %s", txs[0].Hash().Hex(), meta.TxHash)
	}
}

func TestHandleABIRequest(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 24 * time.Hour})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	s.Run()
	defer s.Stop()

	input := &dto.ABI{
		MethodIDs: []string{"0x000001", "0x000002"},
		EventIDs:  []string{"0x00000a", "0x00000b"},
	}

	var (
		name = "test"
		errc = make(chan error, 1)
	)

	s.handleRequest(&ABIRequest{Name: name, ABI: input, errc: errc})
	<-errc

	abi := memdb.Get([]byte(name)).(*dto.ABI)
	if len(abi.MethodIDs) != 2 || len(abi.EventIDs) != 2 {
		t.Fatalf("TestHandleABIRequest, want: (MethodIDs 2 EventIDs 2) got: (MethodIDs %d EventIDs %d)", len(abi.MethodIDs), len(abi.EventIDs))
	}
}

func TestHandleContractRequest(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	var (
		cp        = checkpoint.New(checkpoint.DefaultBasePath, "handler")
		memdb     = memdb.New()
		fetcher   = fetcher.New(client, cp, &fetcher.Config{PollInterval: 24 * time.Hour})
		engine, _ = cft.NewSoloEngine(nil, memdb, cp)

		s, _ = New(client, fetcher, engine, cp, &Config{AllowProxyContract: false})
	)

	defer func() {
		os.RemoveAll(checkpoint.DefaultBasePath)
	}()

	s.Run()
	defer s.Stop()

	// Remix Storage.sol
	ca, err := mock.DeployContract(client, common.Hex2Bytes("608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220322c78243e61b783558509c9cc22cb8493dde6925aa5e89a08cdf6e22f279ef164736f6c63430008120033"))

	if err != nil {
		t.Fatal(err)
	}

	errc := make(chan error, 1)

	input := &ContractRequest{
		Address: ca,
		errc:    errc,
	}

	s.handleRequest(input)
	<-errc

	contract := memdb.Get([]byte(ca.Hex())).(*dto.Contract)
	if len(contract.Candidates) != 5 {
		t.Fatalf("TestHandleContractRequest, want: 5 got: %d", len(contract.Candidates))
	}
}
