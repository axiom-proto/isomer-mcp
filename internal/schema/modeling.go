package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Scalar defines a built-in value type that domain models can reference.
type Scalar struct {
	Name    string `yaml:"name"`
	Syntax  string `yaml:"syntax,omitempty"`
	Remarks string `yaml:"remarks,omitempty"`
}

// ExpressionClause describes a generic parameter or input expected by an expression.
type ExpressionClause struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Expression defines a reusable type expression, such as ref<T> or collection<T>.
type Expression struct {
	Name    string             `yaml:"name"`
	Syntax  string             `yaml:"syntax"`
	Where   []ExpressionClause `yaml:"where"`
	Remarks string             `yaml:"remarks,omitempty"`
}

// OperatorClause describes an operator input or return value.
type OperatorClause struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// Operator defines an operation that can be used in rules, constraints, and behavior contracts.
type Operator struct {
	Name    string           `yaml:"name"`
	Syntax  string           `yaml:"syntax"`
	Where   []OperatorClause `yaml:"where"`
	Returns []OperatorClause `yaml:"returns"`
	Remarks string           `yaml:"remarks,omitempty"`
}

// PrimitiveFrom describes the source fields used to construct a primitive value.
type PrimitiveFrom struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// PrimitiveMember describes a named value for enum-backed primitives.
type PrimitiveMember struct {
	Name    string `yaml:"name"`
	Value   string `yaml:"value"`
	Remarks string `yaml:"remarks,omitempty"`
}

// Primitive defines a constrained domain value type.
//
// String-like primitives generally use From and Constraints. Enum-like primitives
// generally use Members and can leave From empty in the YAML as "from: []".
type Primitive struct {
	Name        string            `yaml:"name"`
	Base        string            `yaml:"base"`
	From        []PrimitiveFrom   `yaml:"from"`
	Constraints []string          `yaml:"constraints,omitempty"`
	Members     []PrimitiveMember `yaml:"members,omitempty"`
}

// EntityProperty defines a typed field on an entity.
type EntityProperty struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

// EntityRelation defines a named relationship from one entity to another.
type EntityRelation struct {
	Name      string `yaml:"name"`
	OwnsMany  string `yaml:"ownsMany,omitempty"`
	OwnsOne   string `yaml:"ownsOne,omitempty"`
	HasOne    string `yaml:"hasOne,omitempty"`
	HasMany   string `yaml:"hasMany,omitempty"`
	BelongsTo string `yaml:"belongsTo,omitempty"`
}

// Entity defines an aggregate or durable model in the domain.
type Entity struct {
	Name       string           `yaml:"name"`
	Plural     string           `yaml:"plural"`
	Properties []EntityProperty `yaml:"properties"`
	// Relations is loaded from the "rels" key used by the domain model YAML.
	Relations []EntityRelation `yaml:"rels"`
	Rules     []string         `yaml:"rules,omitempty"`
}

// BehaviorContext describes an entity or value made available to a behavior.
type BehaviorContext struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// BehaviorContextList supports either "context: none" or a sequence of context values.
type BehaviorContextList struct {
	None   bool
	Values []BehaviorContext
}

// BehaviorInput defines a named input accepted by a behavior.
type BehaviorInput struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

// BehaviorReturn defines a named value returned by a behavior.
type BehaviorReturn struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// BehaviorException defines a named exceptional condition for a behavior.
type BehaviorException struct {
	Name      string `yaml:"name"`
	Condition string `yaml:"condition"`
}

// UnmarshalYAML accepts the compact "none" form or a list of behavior context values.
func (bcl *BehaviorContextList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		if value.Value == "none" {
			bcl.None = true
			bcl.Values = nil
			return nil
		}
		return fmt.Errorf("unexpected string value: %s", value.Value)
	}

	if value.Kind == yaml.SequenceNode {
		bcl.None = false
		return value.Decode(&bcl.Values)
	}

	return fmt.Errorf("unexpected node kind: %v", value.Kind)
}

// Behavior defines a domain action with inputs, outputs, contracts, and exceptions.
type Behavior struct {
	Name           string              `yaml:"name"`
	Uses           []string            `yaml:"uses"`
	Context        BehaviorContextList `yaml:"context"`
	Inputs         []BehaviorInput     `yaml:"inputs"`
	Returns        []BehaviorReturn    `yaml:"returns"`
	Preconditions  []string            `yaml:"preconditions"`
	Postconditions []string            `yaml:"postconditions"`
	Exceptions     []BehaviorException `yaml:"exceptions"`
}

// ServiceCapabilityInput defines an input accepted by a service capability.
type ServiceCapabilityInput struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

// ServiceCapabilityReturn defines a value returned by a service capability.
type ServiceCapabilityReturn struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Remarks string `yaml:"remarks,omitempty"`
}

// ServiceCapability defines one operation exposed by an external service boundary.
type ServiceCapability struct {
	Name     string                    `yaml:"name"`
	Supports []string                  `yaml:"supports"`
	Inputs   []ServiceCapabilityInput  `yaml:"inputs"`
	Returns  []ServiceCapabilityReturn `yaml:"returns"`
}

// Service groups related capabilities that the domain can depend on.
type Service struct {
	Name         string              `yaml:"name"`
	Capabilities []ServiceCapability `yaml:"capabilities"`
}

// Domain contains the full modeled vocabulary loaded from the YAML domain node.
type Domain struct {
	Scalars     []Scalar     `yaml:"scalars"`
	Expressions []Expression `yaml:"expr"`
	Operators   []Operator   `yaml:"operators"`
	Primitives  []Primitive  `yaml:"primitives"`
	Entities    []Entity     `yaml:"entities"`
	Behaviors   []Behavior   `yaml:"behaviors,omitempty"`
	Services    []Service    `yaml:"services,omitempty"`
}

// DomainRoot is the expected top-level YAML document shape.
type DomainRoot struct {
	Domain Domain `yaml:"domain"`
}
