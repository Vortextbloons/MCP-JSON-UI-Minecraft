package bedrockui

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateRootAtSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantPass bool
	}{
		{
			name:     "root @ inheritance syntax",
			input:    `{"namespace":"test","test_button@common.button":{"type":"button"}}`,
			wantPass: true,
		},
		{
			name:     "root no @ bare name",
			input:    `{"namespace":"test","some_custom_name":{"type":"button"}}`,
			wantPass: true,
		},
		{
			name:     "nested @ in controls array",
			input:    `{"namespace":"test","btn":{"type":"button","controls":[{"label@common.label":{"type":"label","text":"hi"}}]}}`,
			wantPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := ValidateDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected internal error: %v", err)
			}
			if tt.wantPass && len(errors) > 0 {
				t.Errorf("expected PASS but got errors: %v", errors)
			}
			if !tt.wantPass && len(errors) == 0 {
				t.Errorf("expected FAIL but got PASS")
			}
			if len(errors) > 0 {
				for _, e := range errors {
					fmt.Printf("  error: %s\n", e)
				}
			}
		})
	}
}

func TestValidateHints(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantHintCount int
		wantContains  string
	}{
		{
			name:          "root animation object hint",
			input:         `{"namespace":"test","pulse":{"anim_type":"alpha","duration":1}}`,
			wantHintCount: 1,
			wantContains:  "Root-level animation objects are not loaded",
		},
		{
			name:          "panel missing controls hint",
			input:         `{"namespace":"test","box":{"type":"panel","size":[10,10]}}`,
			wantHintCount: 1,
			wantContains:  "Container type 'panel' should include a 'controls' array",
		},
		{
			name:          "factory prop on non-factory hint",
			input:         `{"namespace":"test","grid":{"type":"grid","factory":{"name":"f","control_name":"c"}}}`,
			wantHintCount: 2,
			wantContains:  "'factory' property is only valid on elements with type 'factory'",
		},
		{
			name:          "binding missing binding_name hint",
			input:         `{"namespace":"test","btn":{"type":"button","controls":[{"lbl":{"type":"label","bindings":[{"binding_type":"view"}]}}]}}`,
			wantHintCount: 1,
			wantContains:  "binding is missing required 'binding_name'",
		},
		{
			name:          "valid panel with controls passes",
			input:         `{"namespace":"test","box":{"type":"panel","controls":[{"lbl":{"type":"label"}}]}}`,
			wantHintCount: 0,
			wantContains:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := ValidateDocument(tt.input)
			if err != nil {
				t.Fatalf("unexpected internal error: %v", err)
			}
			if len(errors) != tt.wantHintCount {
				t.Errorf("expected %d hints, got %d: %v", tt.wantHintCount, len(errors), errors)
			}
			if tt.wantContains != "" {
				found := false
				for _, e := range errors {
					if strings.Contains(e, tt.wantContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected hint containing %q, got %v", tt.wantContains, errors)
				}
			}
		})
	}
}
