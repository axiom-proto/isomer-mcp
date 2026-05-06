# Isomer

Isomer is an MCP stdio server that makes large domain model documents navigable for AI agents without loading the entire document into context.

## Problem

Durable domain models describe the complete language of a system — its value types, constrained primitives, aggregate entities, domain behaviors, and external service boundaries — in a single coherent document. That document is authoritative but expensive: loading it wholesale into an agent's context window on every request is wasteful when the agent needs only one entity or one behavior at a time.

## Solution

Isomer indexes a domain model file and exposes each item in it as an individually addressable MCP resource. An agent reads only what it needs, identified by a stable URI of the form `isomer://<kind>/<name>`. The source document stays a pure domain model; Isomer adds no MCP-specific structure to it.

## Domain Model Vocabulary

The schema layer defines six kinds of modeled items:

**Scalar** — a named built-in value type such as a string, integer, decimal, or datetime. Scalars are the leaves of the type system.

**Expression** — a parameterized type form such as a reference, collection, or optional wrapper. Expressions are composed from scalars and other named types.

**Operator** — a named domain operation used in rules and constraints, such as membership tests, equality checks, or temporal validity.

**Primitive** — a constrained domain value type. A primitive is either a scalar with validation constraints or an enumeration of named members with fixed values.

**Entity** — a durable aggregate with typed properties, typed relations to other entities, and a set of invariant rules expressed using operators.

**Behavior** — a domain action with a named context, typed inputs, typed return values, preconditions, postconditions, and a set of named exception cases.

**Service** — an external service boundary, grouping named capabilities that the domain depends on.

## Data Flow

```
YAML domain model file
  → DomainRoot        parsed into nested Go structs
  → IsomerStore       slices normalized into name-keyed maps for O(1) lookup
  → MCP stdio server  JSON-RPC 2.0 over stdin/stdout

resources/list   returns a sorted list of all resource URIs
resources/read   returns one item serialized back to YAML
```

## Commands

```sh
isomer serve <file>      # index and serve a domain model over MCP stdio
isomer validate <file>   # parse and validate a domain model file
isomer version           # print the server version
```

## Project Layout

```text
main.go                    CLI entry point
internal/cmd/              Cobra commands: serve, validate, version
internal/schema/           Domain model types and JSON-RPC/MCP message structures
internal/mcp/              IsomerStore (index) and MCP stdio protocol handler
```

## Build

MCP clients register servers as executable paths, not as `go run` invocations. Building a binary gives the client a stable path to launch the server from and eliminates the compile step from each startup.

```sh
mkdir -p build
go build -o build/isomer .
```

## Development Checks

```sh
gofmt -w main.go internal   # format
```

## License

MPL 2.0. See `LICENSE.md`.
