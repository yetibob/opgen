package cmd

import (
	"encoding/json"
	"fmt"
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
			if format == "json" {
				opcodes, err := i80.GenOpCodes()
				if err != nil {
					panic(err)
				}
				writeOpcodes("./out/opcodes.json", opcodes)
			} else {
				fmt.Println("Currently only supports JSON output")
			}
		},
	}
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

// Execute t
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "json", "output format")
}
