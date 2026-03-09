package main

import (
	"strings"
	"testing"

	"github.com/expr-lang/expr"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		outputExpr  string
		opts        []expr.Option
		input       string
		wantOutput  string
		wantMatched bool
		wantErr     bool
	}{
		{
			name:        "matching line",
			expression:  "age > 30",
			input:       `{"age": 35, "name": "alice"}`,
			wantOutput:  "{\"age\": 35, \"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "non-matching line",
			expression:  "age > 30",
			input:       `{"age": 20, "name": "bob"}`,
			wantMatched: false,
		},
		{
			name:       "multiple lines some match",
			expression: "active == true",
			input: strings.Join([]string{
				`{"active": true, "name": "alice"}`,
				`{"active": false, "name": "bob"}`,
				`{"active": true, "name": "carol"}`,
			}, "\n"),
			wantOutput:  "{\"active\": true, \"name\": \"alice\"}\n{\"active\": true, \"name\": \"carol\"}\n",
			wantMatched: true,
		},
		{
			name:        "invalid json line is skipped",
			expression:  "name == \"alice\"",
			input:       "not-json\n{\"name\": \"alice\"}",
			wantOutput:  "{\"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "invalid expression returns error",
			expression:  "!!!",
			input:       `{"name": "alice"}`,
			wantErr:     true,
		},
		{
			name:        "missing field without flag skips line",
			expression:  "age > 30",
			input:       `{"name": "alice"}`,
			wantMatched: false,
		},
		{
			name:        "allow-missing-fields: missing field is nil",
			expression:  "age == nil",
			opts:        []expr.Option{expr.AllowUndefinedVariables()},
			input:       `{"name": "alice"}`,
			wantOutput:  "{\"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "allow-missing-fields: present field still works",
			expression:  "age > 30",
			opts:        []expr.Option{expr.AllowUndefinedVariables()},
			input:       `{"age": 35, "name": "alice"}`,
			wantOutput:  "{\"age\": 35, \"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "non-bool truthy output matches",
			expression:  "name",
			input:       `{"name": "alice"}`,
			wantOutput:  "{\"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "empty lines are skipped",
			expression:  "name == \"alice\"",
			input:       "\n\n{\"name\": \"alice\"}\n\n",
			wantOutput:  "{\"name\": \"alice\"}\n",
			wantMatched: true,
		},
		{
			name:        "large line within buffer limit",
			expression:  "len(data) > 0",
			input:       "{\"data\": \"" + strings.Repeat("x", 512*1024) + "\"}",
			wantMatched: true,
		},
		{
			name:        "--output prints extracted field",
			expression:  "age > 30",
			outputExpr:  "name",
			input:       "{\"age\": 35, \"name\": \"alice\"}",
			wantOutput:  "alice\n",
			wantMatched: true,
		},
		{
			name:        "--output invalid expression returns error",
			expression:  "age > 30",
			outputExpr:  "!!!",
			input:       "{\"age\": 35}",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			var w strings.Builder
			matched, err := run(tt.expression, tt.outputExpr, tt.opts, r, &w)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if matched != tt.wantMatched {
				t.Errorf("run() matched = %v, want %v", matched, tt.wantMatched)
			}
			if tt.wantOutput != "" && w.String() != tt.wantOutput {
				t.Errorf("run() output = %q, want %q", w.String(), tt.wantOutput)
			}
		})
	}
}
