package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DefaultHeartbeat(ctx context.Context, endpoint string) error {
	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ChainID(ctx)
	return err
}

type Client interface {
	Endpoint() string
	Close()
	ChainID(ctx context.Context) (*big.Int, error)
	SubscribeNewBlock(ctx context.Context, ch chan<- *types.Block) (ethereum.Subscription, error)
	BlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
	BlockLogsByNumber(ctx context.Context, blockNumber uint64) ([]types.Log, error)
	GetLatestBlockNumber(ctx context.Context) (uint64, error)
	GetTransactionsByNumber(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error)
	GetTransactionHashesByNumber(ctx context.Context, blockNumber uint64) ([]common.Hash, error)
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error)
}

var _ = Client(&client{})

type client struct {
	endpoint string
	eth      *ethclient.Client
}

func New(endpoint string) (Client, error) {
	eth, err := ethclient.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	return &client{endpoint, eth}, nil
}

func (c *client) Endpoint() string {
	return c.endpoint
}

func (c *client) Close() {
	c.eth.Close()
}

func (c *client) ChainID(ctx context.Context) (*big.Int, error) {
	return c.eth.ChainID(ctx)
}

func (c *client) SubscribeNewBlock(ctx context.Context, ch chan<- *types.Block) (ethereum.Subscription, error) {
	hch := make(chan *types.Header)
	sub, err := c.eth.SubscribeNewHead(ctx, hch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case header := <-hch:
				block, err := c.eth.BlockByNumber(context.Background(), header.Number)
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

func (c *client) BlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	return c.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

func (c *client) BlockLogsByNumber(ctx context.Context, blockNumber uint64) ([]types.Log, error) {
	bn := big.NewInt(int64(blockNumber))
	return c.eth.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: bn, ToBlock: bn})
}

func (c *client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

func (c *client) GetTransactionsByNumber(ctx context.Context, blockNumber uint64) ([]*types.Transaction, error) {
	block, err := c.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}

func (c *client) GetTransactionHashesByNumber(ctx context.Context, blockNumber uint64) ([]common.Hash, error) {
	hashes := make([]common.Hash, 0)

	block, err := c.eth.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}

	for _, hash := range block.Transactions() {
		hashes = append(hashes, hash.Hash())
	}

	return hashes, nil
}

func (c *client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.eth.TransactionReceipt(ctx, txHash)
}

func (c *client) GetTransactionLogs(ctx context.Context, txHash common.Hash) ([]*types.Log, error) {
	receipt, err := c.eth.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt.Logs, nil
}

func (c *client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return c.eth.CodeAt(ctx, account, blockNumber)
}

func (c *client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return c.eth.StorageAt(ctx, account, key, blockNumber)
}
