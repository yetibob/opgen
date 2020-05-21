package opcode

import "fmt"

// OpCode represent a CPU Operation and associated metadata
type OpCode struct {
	Code  int      `json:"code"`
	Name  string   `json:"name"`
	Desc  string   `json:"description"`
	Size  int      `json:"size"`
	Flags []string `json:"flags"`
}

// ToCase generates a dummy case statement based on the given opcode
func (c OpCode) ToCase() string {
	return fmt.Sprintf("case 0x%02X:", c.Code)
}
