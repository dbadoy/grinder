package grinder

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Grinder extracts the MethodID and EventID from the bytecode
// and returns them.
//
// The return values are just a 'candidates'.
func Grinde(bytecode []byte) (methods []string, events []string, err error) {
	hex := common.Bytes2Hex(bytecode)
	if !strings.Contains(hex, "60806040") {
		return nil, nil, errors.New("must be contract")
	}

	for _, opcode := range getOpCodes(bytecode) {
		switch opcode.code {
		case vm.PUSH4:
			method := common.Bytes2Hex(opcode.value)
			if method == "ffffffff" {
				continue
			}

			methods = append(methods, method)
		case vm.PUSH32:
			events = append(events, common.Bytes2Hex(opcode.value))
		}
	}

	return removeDuplicateString(methods), removeDuplicateString(events), nil
}

func removeDuplicateString(arr []string) (res []string) {
	keys := make(map[string]struct{})
	for _, s := range arr {
		if _, ok := keys[s]; ok {
			continue
		}
		keys[s] = struct{}{}
		res = append(res, s)
	}

	return
}
