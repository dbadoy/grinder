package mock

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func DeployContract(m *Mock, bytecode []byte) (common.Address, error) {
	nonce, err := m.c.NonceAt(context.Background(), m.addr, nil)
	if err != nil {
		return common.Address{}, err
	}

	gasPrice, err := m.c.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Address{}, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Gas:      m.c.Blockchain().GasLimit(),
		GasPrice: gasPrice,
		Data:     bytecode,
	})

	signed, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), m.priv)
	if err != nil {
		return common.Address{}, err
	}

	if err := m.c.SendTransaction(context.Background(), signed); err != nil {
		return common.Address{}, err
	}

	m.c.Commit()

	receipt, err := m.c.TransactionReceipt(context.Background(), signed.Hash())
	if err != nil {
		return common.Address{}, err
	}

	if receipt.Status == 0 {
		return common.Address{}, errors.New("transaction failed")
	}

	return receipt.ContractAddress, nil
}
