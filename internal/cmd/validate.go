package cmd

import (
	"fmt"
	"isomer/internal/schema"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// validateCmd Runs a quick validation pass by loading the target YAML and unmarshalling.
var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a DDM meets Ogham-level compliance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return
		}

		var domain schema.Domain
		err = yaml.Unmarshal(data, &domain)
		if err != nil {
			fmt.Printf("error unmarshalling YAML: %v\n", err)
			return
		}

		fmt.Printf("Validation successful for domain: %s\n", path)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
