package mcp

import "isomer/internal/schema"

// IsomerStore represents a store for entities and behaviors
type IsomerStore struct {
	Scalars     map[string]schema.Scalar
	Operators   map[string]schema.Operator
	Primitives  map[string]schema.Primitive
	Expressions map[string]schema.Expression
	Entities    map[string]schema.Entity
	Behaviors   map[string]schema.Behavior
	Services    map[string]schema.Service
}

// NewStore creates a new IsomerStore from a DomainRoot
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
