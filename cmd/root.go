package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yetibob/opgen/mod"
	"github.com/yetibob/opgen/opcode"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	mods = map[string]mod.Module{
		"i80":   mod.I80{},
		"chip8": mod.Chip8{},
	}
	// Used for flags.
	rootCmd = &cobra.Command{
		Use:   "opgen",
		Short: "A generator for emulator opcodes",
		Run: func(cmd *cobra.Command, args []string) {
			format, err := cmd.PersistentFlags().GetString("format")
			panicErr(err)

			filePath, err := cmd.PersistentFlags().GetString("out")
			panicErr(err)

			modType, err := cmd.PersistentFlags().GetString("mod")
			panicErr(err)

			mod, ok := mods[modType]
			if !ok {
				err = fmt.Errorf("mod %v not supported", modType)
				panicErr(err)
			}

			var file *os.File
			if filePath == "" {
				file = os.Stdout
			} else {
				file, err = os.Create(filePath)
				panicErr(err)
			}

			in, err := cmd.PersistentFlags().GetString("in")
			panicErr(err)

			var opcodes []opcode.OpCode
			if in == "" {
				opcodes, err = mod.GenOpCodes()
				panicErr(err)

			} else {
				opcodes = readOpJSON(in)
			}
			if format == "json" {
				writeOpcodes(file, opcodes)
			} else if format == "go" {
				writeGo(file, opcodes)
			} else if format == "cpp" {
				writeCpp(file, opcodes)
			} else {
				fmt.Println("Currently only supports JSON and Go output")
			}
		},
	}
)

func readOpJSON(fileName string) []opcode.OpCode {
	var ops []opcode.OpCode
	b, err := ioutil.ReadFile(fileName)
	panicErr(err)

	br := bytes.NewReader(b)
	err = json.NewDecoder(br).Decode(&ops)
	panicErr(err)

	return ops
}

func writeGo(file *os.File, codes []opcode.OpCode) {
	gocode := "package opcode\n\nimport (\n\t\"fmt\"\n)\n\n// handleop handles opcodes\nfunc handleop(buf []byte) int {\n\topbytes := 1\n\tswitch buf[0] {\n"
	for _, op := range codes {
		gocode += fmt.Sprintf("\tcase 0x%02x:\n\t\tfmt.println(\"handling opcode: 0x%02x\")\n", op.Code, op.Code)
		name := strings.Split(op.Name, " ")
		if name[0] == "-" {
			name[0] = "NOP"
		}

		spaces := "    "
		if l := len(name[0]); l == 3 {
			spaces += " "
		} else if l == 2 {
			spaces += "  "
		}

		if op.Size == 1 {
			if len(name) == 1 {
				gocode += fmt.Sprintf("\t\tfmt.printf(\"%v\\n\")\n", name[0])
			} else {
				gocode += fmt.Sprintf("\t\tfmt.printf(\"%v%v%v\\n\")\n", name[0], spaces, name[1])
			}
		} else if op.Size == 2 {
			gocode += fmt.Sprintf("\t\tfmt.printf(\"%v%v%v,%v\\t$%%02x\\n\", buf[1])\n", name[0], spaces, "B", "D8")
		} else if op.Size == 3 {
			gocode += fmt.Sprintf("\t\tfmt.printf(\"%v%v%v,%v\\t$%%02x%%02x\\n\", buf[2], buf[1])\n", name[0], spaces, "B", "D16")
		}
		if op.Size > 1 {
			gocode += fmt.Sprintf("\t\topbytes = %v\n", op.Size)
		}
	}
	gocode += "\tdefault:\n\t\tfmt.printf(\"unknown opcode: 0x%02x\\n\", buf[0])\n\t}\n\n\treturn opbytes\n}\n"
	io.WriteString(file, gocode)
}

func writeCpp(file *os.File, codes []opcode.OpCode) {
	cppcode := "#include <fmt/core.h>\n#include <vector>\n\nint handleOp(std::vector<uint8_t> &buf, int pc) {\n\tint opbytes = 1;\n\tuint8_t *code    = &buf[pc];\n\tfmt::print(\"{:04X} : \", pc);\n\n\tswitch (*code) {\n"
	for _, op := range codes {
		cppcode += fmt.Sprintf("\tcase 0x%02X:\n", op.Code)
		name := strings.Split(op.Name, " ")
		if name[0] == "-" {
			name[0] = "NOP"
		}

		if op.Size == 1 {
			if len(name) == 1 {
				cppcode += fmt.Sprintf("\t\tfmt::print(\"0x%02X : %-5v\\n\");\n", op.Code, name[0])
			} else {
				cppcode += fmt.Sprintf("\t\tfmt::print(\"0x%02X : %-5v%-7v\\n\");\n", op.Code, name[0], name[1])
			}
		} else if op.Size == 2 {
			cppcode += fmt.Sprintf("\t\tfmt::print(\"0x%02X : %-5v%-7v${:02X}\\n\", code[1]);\n", op.Code, name[0], name[1])
		} else if op.Size == 3 {
			cppcode += fmt.Sprintf("\t\tfmt::print(\"0x%02X : %-5v%-7v${:02X}{:02X}\\n\", code[2], code[1]);\n", op.Code, name[0], name[1])
		}
		if op.Size > 1 {
			cppcode += fmt.Sprintf("\t\topbytes = %v;\n", op.Size)
		}
		cppcode += "\t\tbreak;\n"
	}
	cppcode += "\tdefault:\n\t\tfmt::print(\"unknown opcode: {:#2X}\\n\", *code);\n\t\tbreak;\n\t}\n\n\treturn opbytes;\n}\n"
	io.WriteString(file, cppcode)
}

func writeOpcodes(file *os.File, codes []opcode.OpCode) {
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	err := enc.Encode(codes)
	panicErr(err)
}

// Execute t
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("mod", "m", "i80", "mod to run. defaults to i80")
	rootCmd.PersistentFlags().StringP("format", "f", "json", "output format")
	rootCmd.PersistentFlags().StringP("in", "i", "", "input file")
	rootCmd.PersistentFlags().StringP("out", "o", "", "name of output file. defaults to std out")
}

func initConfig() {}
