package mod

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/yetibob/opgen/opcode"
)

type Chip8 struct{
}

// GenOpCodes retrieves and generates a json file containing the Opcodes for the Intel 8080
func (module Chip8) GenOpCodes() ([]opcode.OpCode, error) {
	url := "http://www.emulator101.com/8080-by-opcode.html"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// doc, err := html.Parse(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	s, _ := regexp.Compile("<td[^>]*>(.*?)<\\/td>")

	ops := s.FindAllStringSubmatch(string(body), -1)

	var opcodes []opcode.OpCode

	for i := 0; i < len(ops); i += 5 {
		op := strings.TrimSpace(ops[i][1])
		name := strings.TrimSpace(ops[i+1][1])
		size := strings.TrimSpace(ops[i+2][1])
		flags := strings.TrimSpace(ops[i+3][1])
		desc := strings.TrimSpace(ops[i+4][1])

		parsedOp, _ := strconv.ParseInt(op[2:], 16, 64)
		parsedSize, _ := strconv.ParseInt(size, 10, 8)
		parsedFlags := strings.Split(flags, ", ")

		if parsedSize <= 1 {
			parsedSize = 1
		}

		opcode := opcode.OpCode{
			Code:  int(parsedOp),
			Desc:  desc,
			Size:  int(parsedSize),
			Flags: parsedFlags,
			Name:  name,
		}

		opcodes = append(opcodes, opcode)
	}
	return opcodes, nil
}
