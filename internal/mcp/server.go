package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"isomer/internal/schema"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	jsonRPCVersion  = "2.0"
	uriScheme       = "isomer://"
	yamlMimeType    = "application/x-yaml"
	protocolVersion = "2024-11-05"
)

type readResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type readResourceResult struct {
	Contents []readResourceContent `json:"contents"`
}

// handleInitialize responds to the MCP initialize request with the server's
// supported capabilities. Isomer currently exposes read-only resources.
func handleInitialize(req schema.Request) {
	response := schema.Response{
		JsonRPC: jsonRPCVersion,
		Id:      req.Id,
		Result: &schema.HandshakeResponse{
			ProtocolVersion: protocolVersion,
			Capabilities: map[string]interface{}{
				"resources": map[string]interface{}{
					"subscribe":   false,
					"listChanged": false,
				},
				"tools": map[string]interface{}{},
			},
			ServerInfo: map[string]string{
				"name":    "isomer",
				"version": "0.1.3",
			},
		},
	}

	render(response)
}

// handleListResources returns every modeled item as an individual MCP resource.
func handleListResources(req schema.Request, store *IsomerStore) {
	var resourcesList []schema.ResourceListResponseItem

	appendResourceItems(&resourcesList, "scalars", "Scalar", store.Scalars)
	appendResourceItems(&resourcesList, "expressions", "Expression", store.Expressions)
	appendResourceItems(&resourcesList, "operators", "Operator", store.Operators)
	appendResourceItems(&resourcesList, "primitives", "Primitive", store.Primitives)
	appendResourceItems(&resourcesList, "entities", "Entity", store.Entities)
	appendResourceItems(&resourcesList, "behaviors", "Behavior", store.Behaviors)
	appendResourceItems(&resourcesList, "services", "Service", store.Services)

	resp := schema.Response{
		JsonRPC: jsonRPCVersion,
		Id:      req.Id,
		Result: &schema.ResourceListResponse{
			Resources: resourcesList,
		},
	}

	render(resp)
}

// handleReadResource resolves a single isomer://<kind>/<name> URI and returns
// the matching model item serialized back to YAML.
func handleReadResource(req schema.Request, store *IsomerStore) {
	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		renderError(req.Id, -32602, "Invalid params", err.Error())
		return
	}

	content, ok, err := lookupResourceYAML(params.URI, store)
	if err != nil {
		renderError(req.Id, -32602, "Invalid params", err.Error())
		return
	}
	if !ok {
		renderError(req.Id, -32004, "Resource not found", params.URI)
		return
	}

	render(schema.Response{
		JsonRPC: jsonRPCVersion,
		Id:      req.Id,
		Result: readResourceResult{
			Contents: []readResourceContent{
				{
					URI:      params.URI,
					MimeType: yamlMimeType,
					Text:     content,
				},
			},
		},
	})
}

// lookupResourceYAML finds the resource addressed by uri and marshals it to YAML.
func lookupResourceYAML(uri string, store *IsomerStore) (string, bool, error) {
	kind, name, err := parseResourceURI(uri)
	if err != nil {
		return "", false, err
	}

	switch kind {
	case "scalars":
		item, ok := store.Scalars[name]
		return marshalJSON(item, ok)
	case "operators":
		item, ok := store.Operators[name]
		return marshalJSON(item, ok)
	case "primitives":
		item, ok := store.Primitives[name]
		return marshalJSON(item, ok)
	case "expressions":
		item, ok := store.Expressions[name]
		return marshalJSON(item, ok)
	case "entities":
		item, ok := store.Entities[name]
		return marshalJSON(item, ok)
	case "behaviors":
		item, ok := store.Behaviors[name]
		return marshalJSON(item, ok)
	case "services":
		item, ok := store.Services[name]
		return marshalJSON(item, ok)
	default:
		return "", false, fmt.Errorf("unsupported resource kind: %s", kind)
	}
}

// parseResourceURI validates and splits an Isomer resource URI.
//
// Collection-level reads are not supported here. Callers must provide a concrete
// item URI in the form isomer://<kind>/<name>, for example isomer://entities/tenant.
func parseResourceURI(uri string) (string, string, error) {
	if !strings.HasPrefix(uri, uriScheme) {
		return "", "", fmt.Errorf("uri must use %s scheme", uriScheme)
	}

	trimmed := strings.TrimPrefix(uri, uriScheme)
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("uri must be in the format isomer://<kind>/<name>")
	}

	return parts[0], parts[1], nil
}

// marshalYAML serializes a found resource, preserving the not-found case for
// handleReadResource to translate into an MCP error response.
func marshalYAML(value interface{}, ok bool) (string, bool, error) {
	if !ok {
		return "", false, nil
	}

	buffer, err := yaml.Marshal(value)
	if err != nil {
		return "", false, err
	}

	return string(buffer), true, nil
}

func marshalJSON(value interface{}, ok bool) (string, bool, error) {
	if !ok {
		return "", false, nil
	}

	buffer, err := json.Marshal(value)
	if err != nil {
		return "", false, err
	}

	return string(buffer), true, nil
}

// appendResourceItems appends sorted list entries for one resource kind.
func appendResourceItems[T any](items *[]schema.ResourceListResponseItem, kind string, label string, values map[string]T) {
	for _, name := range sortedKeys(values) {
		*items = append(*items, schema.ResourceListResponseItem{
			Uri:         fmt.Sprintf("%s%s/%s", uriScheme, kind, name),
			Name:        name,
			Description: fmt.Sprintf("%s '%s'", label, name),
		})
	}
}

// sortedKeys returns deterministic resource ordering for stable inspector output.
func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// render writes one JSON-RPC response frame to stdout.
func render(res schema.Response) {
	err := json.NewEncoder(os.Stdout).Encode(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding response: %v\n", err)
	}
}

// renderError writes a JSON-RPC error response using the provided request id.
func renderError(id interface{}, code int, message string, data interface{}) {
	render(schema.Response{
		JsonRPC: jsonRPCVersion,
		Id:      id,
		Error: &schema.RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	})
}

// StartServer starts the MCP server and listens for incoming requests.
func StartServer(store *IsomerStore) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		buffer := scanner.Bytes()
		var req schema.Request
		if err := json.Unmarshal(buffer, &req); err != nil {
			renderError(nil, -32700, "Parse error", err.Error())
			continue
		}

		switch req.Method {
		case "initialize":
			handleInitialize(req)
		case "resources/list":
			handleListResources(req, store)
		case "resources/read":
			handleReadResource(req, store)
		default:
			renderError(req.Id, -32601, "Method not found", req.Method)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
