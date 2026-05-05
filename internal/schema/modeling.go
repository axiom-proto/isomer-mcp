package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Scalar struct {
	Name    string `yaml:"name"`
	Syntax  string `yaml:"syntax,omitempty"`
	Remarks string `yaml:"remarks,omitempty"`
}

type ExpressionClause struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Expression struct {
	Name    string             `yaml:"name"`
	Syntax  string             `yaml:"syntax"`
	Where   []ExpressionClause `yaml:"where"`
	Remarks string             `yaml:"remarks,omitempty"`
}

type OperatorClause struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type Operator struct {
	Name    string           `yaml:"name"`
	Syntax  string           `yaml:"syntax"`
	Where   []OperatorClause `yaml:"where"`
	Returns []OperatorClause `yaml:"returns"`
	Remarks string           `yaml:"remarks,omitempty"`
}

type PrimitiveFrom struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type PrimitiveMember struct {
	Name    string `yaml:"name"`
	Value   string `yaml:"value"`
	Remarks string `yaml:"remarks,omitempty"`
}

type Primitive struct {
	Name        string            `yaml:"name"`
	Base        string            `yaml:"base"`
	From        []PrimitiveFrom   `yaml:"from"`
	Constraints []string          `yaml:"constraints,omitempty"`
	Members     []PrimitiveMember `yaml:"members,omitempty"`
}

type EntityProperty struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

type EntityRelation struct {
	Name     string `yaml:"name"`
	OwnsMany string `yaml:"ownsMany,omitempty"`
	OwnsOne  string `yaml:"ownsOne,omitempty"`
	HasOne   string `yaml:"hasOne,omitempty"`
	HasMany  string `yaml:"hasMany,omitempty"`
}

type Entity struct {
	Name       string           `yaml:"name"`
	Plural     string           `yaml:"plural"`
	Properties []EntityProperty `yaml:"properties"`
	Relations  []EntityRelation `yaml:"relations"`
	Rules      []string         `yaml:"rules,omitempty"`
}

type BehaviorContext struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type BehaviorContextList struct {
	None   bool
	Values []BehaviorContext
}

type BehaviorInput struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

type BehaviorReturn struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type BehaviorException struct {
	Name      string `yaml:"name"`
	Condition string `yaml:"condition"`
}

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

type ServiceCapabilityInput struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

type ServiceCapabilityReturn struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Remarks string `yaml:"remarks,omitempty"`
}

type ServiceCapability struct {
	Name     string                    `yaml:"name"`
	Supports []string                  `yaml:"supports"`
	Inputs   []ServiceCapabilityInput  `yaml:"inputs"`
	Returns  []ServiceCapabilityReturn `yaml:"returns"`
}

type Service struct {
	Name         string              `yaml:"name"`
	Capabilities []ServiceCapability `yaml:"capabilities"`
}

type Domain struct {
	Scalars     []Scalar     `yaml:"scalars"`
	Expressions []Expression `yaml:"expr"`
	Operators   []Operator   `yaml:"operators"`
	Primitives  []Primitive  `yaml:"primitives"`
	Entities    []Entity     `yaml:"entities"`
	Behaviors   []Behavior   `yaml:"behaviors,omitempty"`
	Services    []Service    `yaml:"services,omitempty"`
}

type DomainRoot struct {
	Domain Domain `yaml:"domain"`
}
