package server

import (
	"github.com/dbadoy/grinder/server/dto"
	"github.com/ethereum/go-ethereum/common"
)

const (
	abiRequestType = 0x01 + iota
	contractRequestType
)

var (
	_, _ request = (*ABIRequest)(nil), (*ContractRequest)(nil)
)

type request interface {
	Errorc() chan<- error
	Kind() byte
}

type ABIRequest struct {
	Name string
	ABI  *dto.ABI
	errc chan error
}

type ContractRequest struct {
	Address common.Address
	errc    chan error
}

func (a *ABIRequest) Errorc() chan<- error { return a.errc }
func (ABIRequest) Kind() byte              { return abiRequestType }

func (c *ContractRequest) Errorc() chan<- error { return c.errc }
func (ContractRequest) Kind() byte              { return contractRequestType }
