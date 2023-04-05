package mock

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// simulated backend always uses chainID 1337.
	chainID = big.NewInt(1337)

	_ ethclient.Client = (*mock)(nil)
)

type mock struct {
	c          *backends.SimulatedBackend
	privateKey *ecdsa.PrivateKey

	supportSubscribe bool
}

func New(hexPriv string) (*mock, error) {
	private, err := crypto.HexToECDSA(hexPriv)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(private, chainID)
	if err != nil {
		return nil, err
	}

	address := auth.From
	balance, _ := new(big.Int).SetString("10000000000000000000", 10) // 10 eth in wei

	genesisAlloc := map[common.Address]core.GenesisAccount{
		address: {
			Balance: balance,
			Nonce:   1,
		},
	}

	blockGasLimit := uint64(15000000)
	return &mock{
		c:                backends.NewSimulatedBackend(genesisAlloc, blockGasLimit),
		privateKey:       private,
		supportSubscribe: false,
	}, nil
}

func (m *mock) Endpoint() string                          { return "mock" }
func (m *mock) Close()                                    {}
func (m *mock) ChainID(context.Context) (*big.Int, error) { return chainID, nil }

func (m *mock) SubscribeNewBlock(ctx context.Context, ch chan<- *types.Block) (ethereum.Subscription, error) {
	if !m.supportSubscribe {
		return nil, errors.New("not support Subscribe")
	}

	hch := make(chan *types.Header)
	sub, err := m.c.SubscribeNewHead(ctx, hch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case header := <-hch:
				block, err := m.c.BlockByNumber(context.Background(), header.Number)
				if err != nil {
					return
				}
				ch <- block

			case <-sub.Err():
				return
			}
		}
	}()

	return sub, nil
}

func (m *mock) BlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	return m.c.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
}

func (m *mock) BlockLogsByNumber(ctx context.Context, blockNumber uint64) ([]types.Log, error) {
	bn := big.NewInt(int64(blockNumber))
	return m.c.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: bn, ToBlock: bn})
}

func (m *mock) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return m.c.Blockchain().CurrentBlock().Number.Uint64(), nil
}

func (m *mock) GetTransactionsByNumber(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error) {
	block, err := m.c.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

func (m *mock) GetTransactionHashesByNumber(ctx context.Context, blockNumber uint64) ([]common.Hash, error) {
	hashes := make([]common.Hash, 0)

	block, err := m.c.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}

	for _, hash := range block.Transactions() {
		hashes = append(hashes, hash.Hash())
	}

	return hashes, nil
}

func (m *mock) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return m.c.TransactionReceipt(ctx, txHash)
}

func (m *mock) GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	receipt, err := m.c.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt.Logs, nil
}

func (m *mock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return m.c.CodeAt(ctx, account, blockNumber)
}

func (m *mock) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return m.c.StorageAt(ctx, account, key, blockNumber)
}
