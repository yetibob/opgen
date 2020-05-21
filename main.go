package main

import (
	"encoding/json"
	"os"

	"github.com/yetibob/opgen/mod/i80"
	"github.com/yetibob/opgen/opcode"
)

func writeOpcodes(fileName string, codes []opcode.OpCode) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	err = enc.Encode(codes)
	if err != nil {
		panic(err)
	}
}

func main() {
	opcodes, err := i80.GenOpCodes()
	if err != nil {
		panic(err)
	}
	writeOpcodes("./out/opcodes.json", opcodes)
}
