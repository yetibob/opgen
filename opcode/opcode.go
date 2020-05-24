package opcode

// OpCode represent a CPU Operation and associated metadata
type OpCode struct {
	Code  int      `json:"code"`
	Name  string   `json:"name"`
	Desc  string   `json:"description"`
	Size  int      `json:"size"`
	Flags []string `json:"flags"`
}
