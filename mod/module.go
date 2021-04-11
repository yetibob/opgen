package mod

import (
	"github.com/yetibob/opgen/opcode"
)

type Module interface {
	GenOpCodes() ([]opcode.OpCode, error)
}
