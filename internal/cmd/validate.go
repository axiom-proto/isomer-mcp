package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/ddm_1.0.0-ogham.json
var schemaBytes []byte

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

		var source interface{}
		err = yaml.Unmarshal(data, &source)
		if err != nil {
			fmt.Printf("error unmarshalling YAML: %v\n", err)
			return
		}

		jsonData, err := json.Marshal(source)
		if err != nil {
			fmt.Printf("error handling JSON marshalling: %v\n", err)
			return
		}

		var normalized interface{}
		if err := json.Unmarshal(jsonData, &normalized); err != nil {
			fmt.Printf("error handling JSON unmarshalling: %v\n", err)
			return
		}

		compiler := jsonschema.NewCompiler()
		compiler.AddResource("schema.json", bytes.NewReader(schemaBytes))
		compiledSchema, _ := compiler.Compile("schema.json")

		if err := compiledSchema.Validate(normalized); err != nil {
			if validationErr, ok := errors.AsType[*jsonschema.ValidationError](err); ok {
				for _, cause := range validationErr.BasicOutput().Errors {
					fmt.Printf("\tfield: %s - %s\n", cause.InstanceLocation, cause.Error)
				}
			}
		} else {
			fmt.Printf("Validation successful for domain: %s\n", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
