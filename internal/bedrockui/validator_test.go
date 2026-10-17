package bedrockui

import (
	"fmt"
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

func TestDocumentContracts(t *testing.T) {
	tests := []struct {
		name, input string
		pass        bool
	}{
		{"root animation", `{"namespace":"test","pulse":{"anim_type":"alpha","duration":1},"box":{"type":"panel","anims":["@test.pulse"]}}`, true},
		{"empty panel", `{"namespace":"test","box":{"type":"panel"}}`, true},
		{"collection factory", `{"namespace":"test","list":{"type":"stack_panel","factory":{"name":"buttons","control_name":"test.button"}}}`, true},
		{"view binding", `{"namespace":"test","box":{"type":"panel","bindings":[{"binding_type":"view","source_property_name":"#value","target_property_name":"#visible"}]}}`, true},
		{"collection details", `{"namespace":"test","box":{"bindings":[{"binding_type":"collection_details","binding_collection_name":"form_buttons"}]}}`, true},
		{"RGB label", `{"namespace":"test","label":{"type":"label","color":[1,1,1]}}`, true},
		{"nested invalid binding", `{"namespace":"test","box":{"controls":[{"inner":{"controls":[{"label":{"bindings":[{"binding_type":"global"}]}}]}}]}}`, false},
		{"nested invalid size", `{"namespace":"test","box":{"controls":[{"inner":{"size":[1]}}]}}`, false},
		{"unnamed child", `{"namespace":"test","box":{"controls":[{"type":"label"}]}}`, false},
		{"view missing target", `{"namespace":"test","box":{"bindings":[{"binding_type":"view","source_property_name":"#value"}]}}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := ValidateDocument(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if (len(issues) == 0) != tt.pass {
				t.Fatalf("pass=%v, issues=%v", tt.pass, issues)
			}
		})
	}
}
