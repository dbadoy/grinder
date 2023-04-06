package server

import (
	"bytes"
	"context"
	"errors"

	"github.com/dbadoy/grinder/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	emptySlot = make([]byte, 32)
)

func (s *Server) deployTransaction(tx *types.Transaction) (common.Address, error) {
	return contractAddress(tx)
}

// Transparent proxy pattern
func (s *Server) eip1967(tx *types.Transaction) (common.Address, common.Address, error) {
	ca, err := contractAddress(tx)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	admin, err := s.eth.StorageAt(context.Background(), ca, common.HexToHash(params.AdminAddressSlotEIP1967), nil)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	impl, err := s.eth.StorageAt(context.Background(), ca, common.HexToHash(params.ImplementationAddressSlotEIP1967), nil)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	if bytes.Equal(admin, emptySlot) || bytes.Equal(impl, emptySlot) {
		return common.Address{}, common.Address{}, errors.New("empty eip1967 contract address slot")
	}

	return common.BytesToAddress(admin), common.BytesToAddress(impl), nil
}

// UUPS proxy pattern
func (s *Server) eip1822(tx *types.Transaction) (common.Address, error) {
	ca, err := contractAddress(tx)
	if err != nil {
		return common.Address{}, err
	}

	// keccak256("PROXIABLE")
	impl, err := s.eth.StorageAt(context.Background(), ca, common.HexToHash(params.LogicAddressSlotEIP1822), nil)
	if err != nil {
		return common.Address{}, err
	}

	if bytes.Equal(impl, emptySlot) {
		return common.Address{}, errors.New("empty eip1822 logic contract address slot")
	}

	return common.BytesToAddress(impl), nil
}

func contractAddress(tx *types.Transaction) (common.Address, error) {
	if tx.To() != nil || tx.Data() == nil {
		return common.Address{}, errors.New("this is not deploy transaction")
	}

	return crypto.CreateAddress(calcFrom(tx), tx.Nonce()), nil
}

func calcFrom(tx *types.Transaction) common.Address {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		panic(err)
	}
	return from
}
