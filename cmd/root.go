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
		"i80": mod.I80{},
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
				err = fmt.Errorf("Modtype %v not supported", modType)
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
	goCode := "package opcode\n\nimport (\n\t\"fmt\"\n)\n\n// HandleOp handles opcodes\nfunc HandleOp(buf []byte) int {\n\topbytes := 1\n\tswitch buf[0] {\n"
	for _, op := range codes {
		goCode += fmt.Sprintf("\tcase 0x%02X:\n\t\tfmt.Println(\"Handling OpCode: 0x%02X\")\n", op.Code, op.Code)
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
				goCode += fmt.Sprintf("\t\tfmt.Printf(\"%v\\n\")\n", name[0])
			} else {
				goCode += fmt.Sprintf("\t\tfmt.Printf(\"%v%v%v\\n\")\n", name[0], spaces, name[1])
			}
		} else if op.Size == 2 {
			goCode += fmt.Sprintf("\t\tfmt.Printf(\"%v%v%v,%v\\t$%%02X\\n\", buf[1])\n", name[0], spaces, "B", "D8")
		} else if op.Size == 3 {
			goCode += fmt.Sprintf("\t\tfmt.Printf(\"%v%v%v,%v\\t$%%02X%%02X\\n\", buf[2], buf[1])\n", name[0], spaces, "B", "D16")
		}
		if op.Size > 1 {
			goCode += fmt.Sprintf("\t\topbytes = %v\n", op.Size)
		}
	}
	goCode += "\tdefault:\n\t\tfmt.Printf(\"Unknown OpCode: 0x%02X\\n\", buf[0])\n\t}\n\n\treturn opbytes\n}\n"
	io.WriteString(file, goCode)
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

func initConfig() {
	return
}
