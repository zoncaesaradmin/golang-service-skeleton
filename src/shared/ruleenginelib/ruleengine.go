package ruleenginelib

import (
	"sync"
)

type MatchedResults []Action

type EvaluatorOptions struct {
	AllowUndefinedVars bool
	FirstMatch         bool
}

var defaultOptions = &EvaluatorOptions{
	AllowUndefinedVars: false,
	FirstMatch:         true,
}

// RuleEngine represents the main rule engine with its configuration and state
type RuleEngine struct {
	EvaluatorOptions
	RuleMap map[string]RuleBlock
	Results MatchedResults
	Mutex   sync.Mutex
	//Logger *Logger
	RuleTypess []string
}

// EvaluateStruct evaluates a single rule against the provided data
func (re *RuleEngine) EvaluateStruct(jsonText *RuleEntry, identifier Data) bool {
	return EvaluateRule(jsonText, identifier, &Options{
		AllowUndefinedVars: re.AllowUndefinedVars,
	})
}

// AddRule adds a new rule to the engine
func (re *RuleEngine) AddRule(rule string) *RuleEngine {
	ruleBlock := ParseJSON(rule)
	re.Mutex.Lock()
	defer re.Mutex.Unlock()
	re.RuleMap[ruleBlock.UUID] = *ruleBlock
	return re
}

// DeleteRule removes a rule from the engine
func (re *RuleEngine) DeleteRule(rule string) {
	ruleBlock := ParseJSON(rule)
	re.Mutex.Lock()
	defer re.Mutex.Unlock()
	delete(re.RuleMap, ruleBlock.UUID)
}

// EvaluateRules evaluates all rules against the provided data
func (re *RuleEngine) EvaluateRules(data Data) (bool, string, *RuleEntry) {
	re.Mutex.Lock()
	defer re.Mutex.Unlock()

	for _, ruleBlock := range re.RuleMap {
		for _, rule := range ruleBlock.RuleEntries {
			if re.EvaluateStruct(rule, data) {
				if defaultOptions.FirstMatch {
					return true, ruleBlock.UUID, rule
				}
			}
		}
	}
	return false, "", nil
}

// NewRuleEngineInstance creates a new instance of the RuleEngine with the given options
func NewRuleEngineInstance(options *EvaluatorOptions) *RuleEngine {
	opts := options
	if opts == nil {
		opts = defaultOptions
	}

	return &RuleEngine{
		EvaluatorOptions: *opts,
	}
}
