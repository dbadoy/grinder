package grinder

import (
	"github.com/ethereum/go-ethereum/core/vm"
)

type opCode struct {
	code  vm.OpCode
	value []byte
}

func getOpCodes(bytecode []byte) []*opCode {
	opcodes := make([]*opCode, 0)

	for index := 0; index < len(bytecode); index++ {
		code := vm.OpCode(bytecode[index])
		if code.IsPush() {
			pLen := bytecode[index] - 0x5f
			pData := bytecode[index+1 : index+int(pLen)+1]

			opcodes = append(opcodes, &opCode{code, pData})
			index += int(pLen)
		} else {
			opcodes = append(opcodes, &opCode{code, nil})
		}
	}

	return opcodes
}
