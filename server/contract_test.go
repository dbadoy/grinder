package server

import (
	"context"
	"testing"

	"github.com/dbadoy/grinder/pkg/ethclient/mock"
	"github.com/ethereum/go-ethereum/common"
)

type data struct {
	bytecode   string
	isDeployTx bool
	eip1822    bool
	eip1967    bool
}

var (
	// Using bytecode with deploy
	testset = []data{
		/* Remix Storage.sol */ {
			bytecode:   "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220322c78243e61b783558509c9cc22cb8493dde6925aa5e89a08cdf6e22f279ef164736f6c63430008120033",
			isDeployTx: true,
			eip1822:    false,
			eip1967:    false,
		},
	}
)

func TestContractHandle(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	s := &Server{eth: client}

	for i, elem := range testset {
		var (
			blockNumber = i + 1
			bytecode    = common.Hex2Bytes(elem.bytecode)
		)
		ca, err := mock.DeployContract(client, bytecode)
		if err != nil {
			t.Fatal(err)
		}

		txs, err := client.GetTransactionsByNumber(context.Background(), uint64(blockNumber))
		if err != nil {
			t.Fatal(err)
		}

		if _, err := contractAddress(txs[0]); err != nil && elem.isDeployTx {
			t.Fatalf("TestContractHandle - deploy transaction, want: success got: failed (%v)", err)
		}

		if _, err := s.eip1822(ca); err != nil && elem.eip1822 {
			t.Fatalf("TestContractHandle - eip1822, want: success got: failed (%v)", err)
		}

		if _, _, err := s.eip1967(ca); err != nil && elem.eip1967 {
			t.Fatalf("TestContractHandle - eip1967, want: success got: failed (%v)", err)
		}
	}
}

func TestEIP1822(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	s := &Server{eth: client}

	if _, err := s.eip1822(common.HexToAddress(mock.PrecompiledContractEIP1822)); err != nil {
		t.Fatalf("TestEIP1822, want: sucess got: failed (%v)", err)
	}
}

func TestEIP1967(t *testing.T) {
	client, err := mock.New(mock.DefaultPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	s := &Server{eth: client}

	if _, _, err := s.eip1967(common.HexToAddress(mock.PrecompiledContractEIP1967)); err != nil {
		t.Fatalf("TestEIP1967, want: sucess got: failed (%v)", err)
	}
}
