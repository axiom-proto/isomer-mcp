package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "isomer",
	Short: "Isomer MCP support for Axiom Protocol",
	Long:  "Isomer is a supporting tool for Axiom Protocol durable domain modeling files served over stdio for MCP support.",
}

func Execute() error {
	return rootCmd.Execute()
}
