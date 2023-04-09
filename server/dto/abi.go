package dto

import "github.com/ethereum/go-ethereum/accounts/abi"

type ABI struct {
	MethodIDs  []string
	MethodSigs []string
	EventIDs   []string
	EventSigs  []string
}

func (ABI) Index() string {
	return "abis"
}

func PackABI(abi *abi.ABI) *ABI {
	var (
		methods = abi.Methods
		events  = abi.Events

		mids  = make([]string, len(methods))
		msigs = make([]string, len(methods))
		eids  = make([]string, len(methods))
		esigs = make([]string, len(methods))
	)

	for _, method := range methods {
		mids = append(mids, string(method.ID))
		msigs = append(msigs, method.Sig)
	}

	for _, event := range events {
		eids = append(eids, event.ID.Hex())
		esigs = append(esigs, event.Sig)
	}

	return &ABI{
		MethodIDs:  mids,
		MethodSigs: msigs,
		EventIDs:   eids,
		EventSigs:  esigs,
	}
}
