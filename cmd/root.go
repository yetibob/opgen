package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/yetibob/opgen/mod/i80"
	"github.com/yetibob/opgen/opcode"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "opgen",
		Short: "A generator for emulator opcodes",
		// 		Long: `Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			format, err := cmd.PersistentFlags().GetString("format")
			if err != nil {
				panic(err)
			}
			filePath, err := cmd.PersistentFlags().GetString("out")
			var file *os.File
			if filePath == "" {
				file = os.Stdout
			} else {
				file, err = os.Create(filePath)
				if err != nil {
					panic(err)
				}
			}
			if err != nil {
				panic(err)
			}
			in, err := cmd.PersistentFlags().GetString("in")
			if err != nil {
				panic(err)
			}
			if format == "json" {
				opcodes, err := i80.GenOpCodes()
				if err != nil {
					panic(err)
				}
				writeOpcodes(file, opcodes)
			} else if format == "go" {
				var opcodes []opcode.OpCode
				if in == "" {
					opcodes, err = i80.GenOpCodes()
					if err != nil {
						panic(err)
					}
				} else {
					opcodes = readOpJSON(in)
				}
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
	if err != nil {
		panic(err)
	}
	br := bytes.NewReader(b)
	err = json.NewDecoder(br).Decode(&ops)
	if err != nil {
		panic(err)
	}
	return ops
}

func writeGo(file *os.File, codes []opcode.OpCode) {
	goCode := "package opcode\n\nimport (\n\t\"fmt\"\n)\n\n// HandleOp handles opcodes\nfunc HandleOp(op byte) int {\n\tvar opbytes int\n\tswitch op {\n"
	for _, op := range codes {
		goCode += fmt.Sprintf("\tcase 0x%02X:\n\t\tfmt.Println(\"Handling OpCode: 0x%02X\")\n\t\topbytes = %v\n", op.Code, op.Code, op.Size)
	}
	goCode += "\tdefault:\n\t\tfmt.Printf(\"Unrecognized OpCode: 0x%02X\", op)\n\t}\n\n\treturn opbytes\n}\n"
	io.WriteString(file, goCode)
}

func writeOpcodes(file *os.File, codes []opcode.OpCode) {
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	err := enc.Encode(codes)
	if err != nil {
		panic(err)
	}
}

// Execute t
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("format", "f", "json", "output format")
	rootCmd.PersistentFlags().StringP("in", "i", "", "input file")
	rootCmd.PersistentFlags().StringP("out", "o", "", "name of output file. defaults to std out")
}

func initConfig() {
	return
}