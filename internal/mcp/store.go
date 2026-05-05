package mcp

import "isomer/internal/schema"

// IsomerStore indexes a parsed domain by resource kind and item name.
//
// The MCP server serves one resource per modeled item, so the YAML slices in
// schema.Domain are normalized into maps for lookup by URIs such as
// isomer://primitives/label.
type IsomerStore struct {
	Scalars     map[string]schema.Scalar
	Operators   map[string]schema.Operator
	Primitives  map[string]schema.Primitive
	Expressions map[string]schema.Expression
	Entities    map[string]schema.Entity
	Behaviors   map[string]schema.Behavior
	Services    map[string]schema.Service
}

// NewStore creates an IsomerStore from a parsed DomainRoot.
//
// Duplicate names in the same resource kind are resolved by the last item in
// the YAML document because each item is assigned into a map by name.
func NewStore(d schema.DomainRoot) *IsomerStore {
	store := &IsomerStore{
		Scalars:     make(map[string]schema.Scalar),
		Operators:   make(map[string]schema.Operator),
		Primitives:  make(map[string]schema.Primitive),
		Expressions: make(map[string]schema.Expression),
		Entities:    make(map[string]schema.Entity),
		Behaviors:   make(map[string]schema.Behavior),
		Services:    make(map[string]schema.Service),
	}

	for _, e := range d.Domain.Scalars {
		store.Scalars[e.Name] = e
	}

	for _, e := range d.Domain.Operators {
		store.Operators[e.Name] = e
	}

	for _, e := range d.Domain.Primitives {
		store.Primitives[e.Name] = e
	}

	for _, e := range d.Domain.Expressions {
		store.Expressions[e.Name] = e
	}

	for _, e := range d.Domain.Entities {
		store.Entities[e.Name] = e
	}

	for _, b := range d.Domain.Behaviors {
		store.Behaviors[b.Name] = b
	}

	for _, s := range d.Domain.Services {
		store.Services[s.Name] = s
	}

	return store
}
