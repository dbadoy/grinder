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

	for _, opcode := range getOpCodes([]byte(hex)) {
		switch opcode.code {
		case vm.PUSH4:
			methods = append(methods, common.Bytes2Hex(opcode.value))
		case vm.PUSH32:
			events = append(events, common.Bytes2Hex(opcode.value))
		}
	}
	return
}
