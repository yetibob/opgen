package mod

import (
	"github.com/yetibob/opgen/opcode"
)

type Chip8 struct{
}

// GenOpCodes retrieves and generates a json file containing the Opcodes for the Intel 8080
func (module Chip8) GenOpCodes() ([]opcode.OpCode, error) {
	return []opcode.OpCode{}, nil
}
