package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	claudeMCPConfigPath = ".mcp.json"
	isomerServerName    = "isomer"
	vscodeMCPConfigPath = ".vscode/mcp.json"
)

var useVsCode bool
var useClaude bool
var createMissing bool

var registerCmd = &cobra.Command{
	Use:   "register path",
	Short: "Registers the MCP provider using a durable domain modeling file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		if _, err := exec.LookPath("isomer"); err != nil {
			fmt.Println("isomer command not in your local path")
			return
		}

		if useVsCode {
			handleVsCodeRegistration(path, createMissing)
		} else if useClaude {
			handleClaudeRegistration(path, createMissing)
		} else {
			fmt.Println("No registration method specified.")
		}
	},
}

func init() {
	registerCmd.Flags().BoolVar(&useVsCode, "vscode", false, "integrate with VS Code registration")
	registerCmd.Flags().BoolVar(&useClaude, "claude-code", false, "integrate with Claude Code registration")
	registerCmd.Flags().BoolVar(&createMissing, "create-missing", false, "create registration files if missing")

	rootCmd.AddCommand(registerCmd)
}

func handleVsCodeRegistration(path string, createIfMissing bool) {
	if err := registerMCPServer(vscodeMCPConfigPath, "servers", path, createIfMissing); err != nil {
		fmt.Fprintln(os.Stderr, "VS Code registration failed:", err)
		return
	}

	fmt.Printf("Registered isomer MCP server in %s\n", vscodeMCPConfigPath)
}

func handleClaudeRegistration(path string, createIfMissing bool) {
	if err := registerMCPServer(claudeMCPConfigPath, "mcpServers", path, createIfMissing); err != nil {
		fmt.Fprintln(os.Stderr, "Claude Code registration failed:", err)
		return
	}

	fmt.Printf("Registered isomer MCP server in %s\n", claudeMCPConfigPath)
}

func registerMCPServer(configPath string, serverCollectionKey string, domainModelPath string, createIfMissing bool) error {
	config, err := readJSONConfig(configPath, createIfMissing)
	if err != nil {
		return err
	}

	commandPath := "isomer"

	serverCollections, err := getServerCollection(config, serverCollectionKey)
	if err != nil {
		return err
	}

	serverCollections[isomerServerName] = map[string]any{
		"type":    "stdio",
		"command": commandPath,
		"args":    []string{"serve", domainModelPath},
	}
	config[serverCollectionKey] = serverCollections

	return writeJSONConfig(configPath, config)
}

func readJSONConfig(configPath string, createIfMissing bool) (map[string]any, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %s: %w", configPath, err)
		}

		if !createIfMissing {
			return nil, fmt.Errorf("%s does not exist; pass --create-missing to create it", configPath)
		}

		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return nil, fmt.Errorf("create parent directory for %s: %w", configPath, err)
		}

		return map[string]any{}, nil
	}

	if len(data) == 0 {
		return map[string]any{}, nil
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse %s: %w", configPath, err)
	}

	if config == nil {
		config = map[string]any{}
	}

	return config, nil
}

func getServerCollection(config map[string]any, key string) (map[string]any, error) {
	rawServers, ok := config[key]
	if !ok {
		return map[string]any{}, nil
	}

	serverCollections, ok := rawServers.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%q must be a JSON object", key)
	}

	return serverCollections, nil
}

func writeJSONConfig(configPath string, config map[string]any) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize %s: %w", configPath, err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", configPath, err)
	}

	return nil
}
