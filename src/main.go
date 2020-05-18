package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

func main() {
	url := "http://www.emulator101.com/8080-by-opcode.html"
	resp, _ := http.Get(url)

	// doc, err := html.Parse(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	s, _ := regexp.Compile("<td[^>]*>(.*?)<\\/td>")

	ops := s.FindAllStringSubmatch(string(body), -1)

	var opcodes []OpCode

	for i := 0; i < len(ops); i += 5 {
		op := ops[i]
		ins := ops[i+1]
		size := ops[i+2]
		flags := ops[i+3]
		f := ops[i+4]
		parsedOp, _ := strconv.ParseInt(strings.TrimSpace(op[1][2:]), 16, 8)
		name := strings.TrimSpace(ins[1])
		parsedSize, _ := strconv.ParseInt(strings.TrimSpace(size[1]), 10, 8)
		parsedFlags := strings.Split(strings.TrimSpace(flags[1]), ", ")
		parsedDesc := strings.TrimSpace(f[1])
		opcode := OpCode{
			Code:  int(parsedOp),
			Desc:  parsedDesc,
			Size:  int(parsedSize),
			Flags: parsedFlags,
			Name:  name,
		}
		opcodes = append(opcodes, opcode)
	}

	fmt.Println(opcodes)
}
