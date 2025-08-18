package ruleenginelib

type results []Event

type EvaluatorOptions struct {
	AllowUndefinedVars bool
}

var defaultOptions = &EvaluatorOptions{
	AllowUndefinedVars: false,
}

type RuleEngine struct {
	EvaluatorOptions
	Rules   []string
	Results results
}

func (re *RuleEngine) EvaluateStruct(jsonText *Rule, identifier Data) bool {
	return EvaluateRule(jsonText, identifier, &Options{
		AllowUndefinedVars: re.AllowUndefinedVars,
	})
}

func (re *RuleEngine) AddRule(rule string) *RuleEngine {
	re.Rules = append(re.Rules, rule)
	return re
}

func (re *RuleEngine) AddRules(rules ...string) *RuleEngine {
	re.Rules = append(re.Rules, rules...)
	return re
}

func (re *RuleEngine) EvaluateRules(data Data) results {
	for _, j := range re.Rules {
		rule := ParseJSON(j)

		if re.EvaluateStruct(rule, data) {
			re.Results = append(re.Results, rule.Event)
		}
	}
	return re.Results
}

func New(options *EvaluatorOptions) *RuleEngine {
	opts := options
	if opts == nil {
		opts = defaultOptions
	}

	return &RuleEngine{
		EvaluatorOptions: *opts,
	}
}
