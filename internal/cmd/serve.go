package cmd

import (
	"fmt"
	"isomer/internal/mcp"
	"isomer/internal/schema"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var serveCmd = &cobra.Command{
	Use:   "serve [file]",
	Short: "Starts the MCP server on stdio for agent interaction",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		data, _ := os.ReadFile(path)
		var domain schema.DomainRoot
		err := yaml.Unmarshal(data, &domain)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to parse YAML file:", err)
			return
		}

		store := mcp.NewStore(domain)
		mcp.StartServer(store)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
