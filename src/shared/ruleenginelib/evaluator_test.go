package ruleenginelib

import (
	"testing"
)

func TestEvaluateConditional(t *testing.T) {
	tests := []struct {
		conditional *Conditional
		identifier  interface{}
		expected    bool
	}{
		{&Conditional{
			Fact:     "name",
			Operator: "eq",
			Value:    "Icheka",
		},
			"Icheka",
			true,
		},
		{&Conditional{
			Fact:     "name",
			Operator: "eq",
			Value:    "Icheka",
		},
			"Ronie",
			false,
		},
	}

	for i, tt := range tests {
		if ok := EvaluateConditional(tt.conditional, tt.identifier); ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateConditional to return %t, got=%t", i, tt.expected, ok)
		}
	}
}

func TestEvaluateAllCondition(t *testing.T) {
	tests := []struct {
		payload struct {
			conditions []Conditional
			identifier Data
		}
		expected bool
	}{
		{
			payload: struct {
				conditions []Conditional
				identifier Data
			}{
				conditions: []Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Neptune",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []Conditional
				identifier Data
			}{
				conditions: []Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		if ok := EvaluateAllCondition(&tt.payload.conditions, tt.payload.identifier); ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateAllCondition to be %t, got=%t", i, tt.expected, ok)
		}
	}
}

func TestEvaluateAnyCondition(t *testing.T) {
	tests := []struct {
		payload struct {
			conditions []Conditional
			identifier Data
		}
		expected bool
	}{
		{
			payload: struct {
				conditions []Conditional
				identifier Data
			}{
				conditions: []Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Neptune",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []Conditional
				identifier Data
			}{
				conditions: []Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "black",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: true,
		},
		{
			payload: struct {
				conditions []Conditional
				identifier Data
			}{
				conditions: []Conditional{
					{
						Fact:     "planet",
						Operator: "eq",
						Value:    "Saturn",
					},
					{
						Fact:     "colour",
						Operator: "eq",
						Value:    "white",
					},
				},
				identifier: Data{
					"planet": "Neptune",
					"colour": "black",
				},
			},
			expected: false,
		},
	}

	for i, tt := range tests {
		if ok := EvaluateAnyCondition(&tt.payload.conditions, tt.payload.identifier); ok != tt.expected {
			t.Errorf("tests[%d] - expected EvaluateAnyCondition to be %t, got=%t", i, tt.expected, ok)
		}
	}
}
