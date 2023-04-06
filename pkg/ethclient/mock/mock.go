package mock

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	// DO NOT USE IT PERSONALLY!!
	// PrivateKey(Public) Generated by Ganache
	DefaultPrivateKey = "4eeebf001564499721307427f37733571dbda44976f68acaa44bd6fac3d820c7"

	// simulated backend always uses chainID 1337.
	chainID = big.NewInt(1337)

	_ ethclient.Client = (*Mock)(nil)
)

type Mock struct {
	c *backends.SimulatedBackend

	addr common.Address
	priv *ecdsa.PrivateKey

	SupportSubscribe bool
}

func New(hexPriv string) (*Mock, error) {
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
	return &Mock{
		c:                backends.NewSimulatedBackend(genesisAlloc, blockGasLimit),
		priv:             private,
		addr:             address,
		SupportSubscribe: false,
	}, nil
}

func (m *Mock) Backend() *backends.SimulatedBackend {
	return m.c
}

func (m *Mock) Endpoint() string                          { return "mock" }
func (m *Mock) Close()                                    {}
func (m *Mock) ChainID(context.Context) (*big.Int, error) { return chainID, nil }

func (m *Mock) SubscribeNewBlock(ctx context.Context, ch chan<- *types.Block) (ethereum.Subscription, error) {
	if !m.SupportSubscribe {
		return nil, rpc.ErrNotificationsUnsupported
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

func (m *Mock) BlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	return m.c.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
}

func (m *Mock) BlockLogsByNumber(ctx context.Context, blockNumber uint64) ([]types.Log, error) {
	bn := big.NewInt(int64(blockNumber))
	return m.c.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: bn, ToBlock: bn})
}

func (m *Mock) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return m.c.Blockchain().CurrentBlock().Number.Uint64(), nil
}

func (m *Mock) GetTransactionsByNumber(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error) {
	block, err := m.c.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

func (m *Mock) GetTransactionHashesByNumber(ctx context.Context, blockNumber uint64) ([]common.Hash, error) {
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

func (m *Mock) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return m.c.TransactionReceipt(ctx, txHash)
}

func (m *Mock) GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	receipt, err := m.c.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt.Logs, nil
}

func (m *Mock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return m.c.CodeAt(ctx, account, blockNumber)
}

func (m *Mock) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return m.c.StorageAt(ctx, account, key, blockNumber)
}
