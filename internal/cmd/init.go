package cmd

import (
	"errors"
	"fmt"
	"isomer/internal/schema"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Initializes a durable domain modeling file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		ddmFile, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				fmt.Printf("file %s already exists\n", path)
				return
			} else {
				fmt.Fprintln(os.Stderr, "file system error: ", err)
			}
		}
		defer ddmFile.Close()

		var domain = schema.DomainRoot{}
		domainData, err := yaml.Marshal(domain)
		if err != nil {
			fmt.Fprintln(os.Stderr, "YAML processing error occurred: ", err)
			return
		}

		ddmFile.Write(domainData)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
