# Isomer

Isomer is a small MCP stdio server for Axiom Protocol durable domain modeling files.

It loads a domain model from YAML, indexes the modeled vocabulary, and exposes each item as an MCP resource. Agents can then inspect a domain by reading focused resources such as `isomer://entities/tenant` or `isomer://primitives/identifier` instead of loading the entire model into context.

## Intent

Durable domain models are useful because they describe the language, rules, entities, behaviors, and service boundaries of a system in one coherent document. Large model files can be expensive for agents to consume directly, though. Isomer turns that document into a navigable resource graph.

The project is designed to make domain inspection:

- Targeted: agents read only the scalar, primitive, entity, behavior, or service they need.
- Deterministic: resources are addressed by stable URIs.
- Friendly to authoring: the source YAML remains a domain model, not an MCP-specific document.
- Easy to test: the server speaks MCP over stdio and can be checked with an MCP inspector.

## Resource Model

Isomer serves one resource per modeled item. The current resource kinds are:

- `scalars`
- `expressions`
- `operators`
- `primitives`
- `entities`
- `behaviors`
- `services`

Resource URIs use this form:

```text
isomer://<kind>/<name>
```

For example:

```text
isomer://entities/tenant
isomer://primitives/identifier
isomer://operators/validFor
```

`resources/list` returns the available resources, and `resources/read` returns one modeled item serialized as YAML.

## Requirements

- Go 1.26.2 or newer, matching `go.mod`
- An MCP-compatible client or inspector for interactive testing

## Quick Start

Build the server:

```sh
mkdir -p build
go build -o build/isomer .
```

Run the sample model over stdio:

```sh
./build/isomer serve ./samples/decide.yml
```

The server reads JSON-RPC messages from stdin and writes JSON-RPC responses to stdout. For a quick manual smoke test:

```sh
printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' \
  '{"jsonrpc":"2.0","id":2,"method":"resources/list"}' \
  | ./build/isomer serve ./samples/decide.yml
```

Read a single resource:

```sh
printf '%s\n' \
  '{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"isomer://entities/tenant"}}' \
  | ./build/isomer serve ./samples/decide.yml
```

## VS Code MCP Configuration

This repository includes an example `.vscode/mcp.json` configuration:

```json
{
  "servers": {
    "isomer-test": {
      "type": "stdio",
      "command": "./build/isomer",
      "args": ["serve", "./samples/decide.yml"]
    }
  },
  "inputs": []
}
```

Build `./build/isomer` first, then start or refresh the MCP server from your client.

## Commands

Show available commands:

```sh
go run . --help
```

Serve a domain model over MCP stdio:

```sh
go run . serve ./samples/decide.yml
```

Validate that a model can be parsed:

```sh
go run . validate ./samples/decide.yml
```

Print the current version:

```sh
go run . version
```

## Project Layout

```text
main.go                    CLI entrypoint
internal/cmd/              Cobra commands
internal/schema/           YAML and JSON-RPC/MCP data structures
internal/mcp/              MCP stdio server and resource store
samples/decide.yml         Example durable domain model
samples/init.json          Minimal initialize request sample
.vscode/mcp.json           Example local MCP server configuration
```

## Development Checks

Format Go files:

```sh
gofmt -w main.go internal
```

Run tests:

```sh
go test ./...
```

Build the local MCP binary:

```sh
mkdir -p build
go build -o build/isomer .
```

## License

This project is licensed under MPL 2.0. See `LICENSE.md`.
