package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Scalar defines a built-in value type that domain models can reference.
type Scalar struct {
	Name    string `yaml:"name" json:"name"`
	Syntax  string `yaml:"syntax,omitempty" json:"syntax,omitempty"`
	Remarks string `yaml:"remarks,omitempty" json:"remarks,omitempty"`
}

// ExpressionClause describes a generic parameter or input expected by an expression.
type ExpressionClause struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

// Expression defines a reusable type expression, such as ref<T> or collection<T>.
type Expression struct {
	Name    string             `yaml:"name" json:"name"`
	Syntax  string             `yaml:"syntax" json:"syntax"`
	Where   []ExpressionClause `yaml:"where" json:"where"`
	Remarks string             `yaml:"remarks,omitempty" json:"remarks,omitempty"`
}

// OperatorClause describes an operator input or return value.
type OperatorClause struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

// Operator defines an operation that can be used in rules, constraints, and behavior contracts.
type Operator struct {
	Name    string           `yaml:"name" json:"name"`
	Syntax  string           `yaml:"syntax" json:"syntax"`
	Where   []OperatorClause `yaml:"where" json:"where"`
	Returns []OperatorClause `yaml:"returns" json:"returns"`
	Remarks string           `yaml:"remarks,omitempty" json:"remarks,omitempty"`
}

// PrimitiveFrom describes the source fields used to construct a primitive value.
type PrimitiveFrom struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

// PrimitiveMember describes a named value for enum-backed primitives.
type PrimitiveMember struct {
	Name    string `yaml:"name" json:"name"`
	Value   string `yaml:"value" json:"value"`
	Remarks string `yaml:"remarks,omitempty" json:"remarks,omitempty"`
}

// Primitive defines a constrained domain value type.
//
// String-like primitives generally use From and Constraints. Enum-like primitives
// generally use Members and can leave From empty in the YAML as "from: []".
type Primitive struct {
	Name        string            `yaml:"name" json:"name"`
	Base        string            `yaml:"base" json:"base"`
	From        []PrimitiveFrom   `yaml:"from" json:"from"`
	Constraints []string          `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	Members     []PrimitiveMember `yaml:"members,omitempty" json:"members,omitempty"`
}

// EntityProperty defines a typed field on an entity.
type EntityProperty struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
}

// EntityRelation defines a named relationship from one entity to another.
type EntityRelation struct {
	Name      string `yaml:"name" json:"name"`
	OwnsMany  string `yaml:"ownsMany,omitempty" json:"ownsMany,omitempty"`
	OwnsOne   string `yaml:"ownsOne,omitempty" json:"ownsOne,omitempty"`
	HasOne    string `yaml:"hasOne,omitempty" json:"hasOne,omitempty"`
	HasMany   string `yaml:"hasMany,omitempty" json:"hasMany,omitempty"`
	BelongsTo string `yaml:"belongsTo,omitempty" json:"belongsTo,omitempty"`
}

// Entity defines an aggregate or durable model in the domain.
type Entity struct {
	Name       string           `yaml:"name" json:"name"`
	Plural     string           `yaml:"plural" json:"plural"`
	Properties []EntityProperty `yaml:"properties" json:"properties"`
	// Relations is loaded from the "rels" key used by the domain model YAML.
	Relations []EntityRelation `yaml:"rels" json:"rels"`
	Rules     []string         `yaml:"rules,omitempty" json:"rules,omitempty"`
}

// BehaviorContext describes an entity or value made available to a behavior.
type BehaviorContext struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

// BehaviorContextList supports either "context: none" or a sequence of context values.
type BehaviorContextList struct {
	None   bool
	Values []BehaviorContext
}

// BehaviorInput defines a named input accepted by a behavior.
type BehaviorInput struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
}

// BehaviorReturn defines a named value returned by a behavior.
type BehaviorReturn struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

// BehaviorException defines a named exceptional condition for a behavior.
type BehaviorException struct {
	Name      string `yaml:"name" json:"name"`
	Condition string `yaml:"condition" json:"condition"`
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
	Name           string              `yaml:"name" json:"name"`
	Uses           []string            `yaml:"uses" json:"uses"`
	Context        BehaviorContextList `yaml:"context" json:"context"`
	Inputs         []BehaviorInput     `yaml:"inputs" json:"inputs"`
	Returns        []BehaviorReturn    `yaml:"returns" json:"returns"`
	Preconditions  []string            `yaml:"preconditions" json:"preconditions"`
	Postconditions []string            `yaml:"postconditions" json:"postconditions"`
	Exceptions     []BehaviorException `yaml:"exceptions" json:"exceptions"`
}

// ServiceCapabilityInput defines an input accepted by a service capability.
type ServiceCapabilityInput struct {
	Name     string `yaml:"name" json:"name"`
	Type     string `yaml:"type" json:"type"`
	Required bool   `yaml:"required" json:"required"`
}

// ServiceCapabilityReturn defines a value returned by a service capability.
type ServiceCapabilityReturn struct {
	Name    string `yaml:"name" json:"name"`
	Type    string `yaml:"type" json:"type"`
	Remarks string `yaml:"remarks,omitempty" json:"remarks,omitempty"`
}

// ServiceCapability defines one operation exposed by an external service boundary.
type ServiceCapability struct {
	Name     string                    `yaml:"name" json:"name"`
	Supports []string                  `yaml:"supports" json:"supports"`
	Inputs   []ServiceCapabilityInput  `yaml:"inputs" json:"inputs"`
	Returns  []ServiceCapabilityReturn `yaml:"returns" json:"returns"`
}

// Service groups related capabilities that the domain can depend on.
type Service struct {
	Name         string              `yaml:"name" json:"name"`
	Capabilities []ServiceCapability `yaml:"capabilities" json:"capabilities"`
}

// Domain contains the full modeled vocabulary loaded from the YAML domain node.
type Domain struct {
	Scalars     []Scalar     `yaml:"scalars" json:"scalars"`
	Expressions []Expression `yaml:"expr" json:"expr"`
	Operators   []Operator   `yaml:"operators" json:"operators"`
	Primitives  []Primitive  `yaml:"primitives" json:"primitives"`
	Entities    []Entity     `yaml:"entities" json:"entities"`
	Behaviors   []Behavior   `yaml:"behaviors,omitempty" json:"behaviors,omitempty"`
	Services    []Service    `yaml:"services,omitempty" json:"services,omitempty"`
}

// DomainRoot is the expected top-level YAML document shape.
type DomainRoot struct {
	Domain Domain `yaml:"domain" json:"domain"`
}
