package bedrockui

import (
	"encoding/json"
	"testing"
)

func TestGeneratedFormsAreDocuments(t *testing.T) {
	forms := map[string]string{}
	for _, variant := range []string{"default", "fullscreen", "grid", "tabs", "light"} {
		forms[variant] = GenerateActionForm(ActionFormParams{Namespace: "test_ui", Variant: variant, IconSupport: true, IncludeScroll: true})
	}
	forms["modal"] = GenerateModalForm(ModalFormParams{Namespace: "test_ui", Controls: []ModalFormControl{{Type: "toggle"}, {Type: "input"}}})
	for name, raw := range forms {
		t.Run(name, func(t *testing.T) {
			var doc map[string]any
			if err := json.Unmarshal([]byte(raw), &doc); err != nil {
				t.Fatalf("expected a JSON object: %v", err)
			}
			issues, err := ValidateDocument(raw)
			if err != nil || len(issues) != 0 {
				t.Fatalf("invalid generated document: %v, %v", err, issues)
			}
			factory := doc["main_content"].(map[string]any)["factory"].(map[string]any)
			if factory["name"] != "server_form_factory" {
				t.Fatal("screen factory must use the engine's server form factory name")
			}
			ids := factory["control_ids"].(map[string]any)
			if ids["long_form"] == nil || ids["custom_form"] == nil {
				t.Fatal("must retain the other vanilla form type")
			}
		})
	}
}

func TestAnalyzeValidBindingForms(t *testing.T) {
	result, err := AnalyzeBindings(`{"namespace":"test","box":{"bindings":[{"binding_name":"#title_text"},{"binding_type":"view","source_property_name":"#title_text","target_property_name":"#text"},{"binding_type":"collection_details","binding_collection_name":"form_buttons"}]}}`)
	if err != nil || len(result.Warnings) != 0 || result.ByType["global"] != 1 {
		t.Fatalf("valid vanilla binding patterns must not warn: %+v, %v", result, err)
	}
}

func TestLocalInheritance(t *testing.T) {
	result, err := ExpandInheritance(`{"namespace":"test","base":{"type":"label","size":[10,10],"controls":[{"child":{"type":"image"}}]},"middle@test.base":{"text":"hello"},"derived@test.middle":{},"external@other.base":{},"cycle_a@test.cycle_b":{},"cycle_b@test.cycle_a":{}}`)
	if err != nil {
		t.Fatal(err)
	}
	var derived, external InheritanceNode
	for _, node := range result.Controls {
		if node.Name == "derived" {
			derived = node
		}
		if node.Name == "external" {
			external = node
		}
	}
	if derived.Type != "label" || derived.Properties["text"] != "hello" || len(derived.Children) != 1 {
		t.Fatalf("lost inherited content: %+v", derived)
	}
	if external.Type != "" {
		t.Fatal("external namespace must not resolve to a same-named local template")
	}
	if len(result.Warnings) < 3 {
		t.Fatalf("expected unresolved external and cycle warnings: %v", result.Warnings)
	}
}

func TestVariablePropertyValues(t *testing.T) {
	issues, err := ValidateDocument(`{"namespace":"test","label":{"type":"label","size":"$size","visible":"(#value = 1)","alpha":"$alpha","anchor_from":"$anchor"}}`)
	if err != nil || len(issues) != 0 {
		t.Fatalf("variable/expression values must pass: %v %v", err, issues)
	}
}
