package ruleenginelib

import (
	"encoding/json"
)

// Conditionals are the basic units of rules
type AstConditional struct {
	Fact     string        `json:"identifier"`
	Operator string        `json:"operator"`
	Value    []interface{} `json:"value"`
}

// A Condition is a group of conditionals within a binding context
// that determines how the group will be evaluated.
type AstCondition struct {
	Any []AstConditional `json:"any"`
	All []AstConditional `json:"all"`
}

// Fired when a identifier matches a rule
type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type RuleBlock struct {
	Type             string       `json:"ruleType,omiempty"`
	SubType          string       `json:"ruleSubType,omiempty"`
	Name             string       `json:"name,omiempty"`
	UUID             string       `json:"uuid,omiempty"`
	Description      string       `json:"description,omiempty"`
	LastModifiedTime int64        `json:"lastModifiedTime,omiempty"`
	State            bool         `json:"state"`
	RuleEntries      []*RuleEntry `json:"payload,omitempty"`
}

type RuleEntry struct {
	Condition AstCondition `json:"condition"`
	Actions   []Action     `json:"actions"`
}

// parse JSON string as Rule
func ParseJSON(j string) *RuleBlock {
	var rule *RuleBlock
	if err := json.Unmarshal([]byte(j), &rule); err != nil {
		panic("expected valid JSON")
	}
	return rule
}
