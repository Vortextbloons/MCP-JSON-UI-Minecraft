package mcpserver

import (
	"strings"
	"testing"

	mcp "github.com/metoro-io/mcp-golang"
	"bedrock-json-ui-mcp/internal/bedrockui"
)

func TestServerNew(t *testing.T) {
	server, err := New()
	if err != nil {
		t.Fatalf("failed to construct MCP server: %v", err)
	}
	if server == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestValidateUIContractResponsePreface(t *testing.T) {
	input := `{"namespace":"test","box":{"type":"panel","size":[10,10]}}`
	resp, err := validateUIContract(ValidateUIContractArgs{UIJSON: input})
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if resp == nil || len(resp.Content) == 0 {
		t.Fatal("expected non-empty response content")
	}

	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, "bedrock://docs/custom_json_ui_design") {
		t.Error("response missing custom_json_ui_design resource link")
	}
	if !strings.Contains(text, "Validation Hints") {
		t.Error("response missing validation hints section")
	}
	if !strings.Contains(text, "Container type 'panel' should include a 'controls' array") {
		t.Error("response missing expected panel controls hint")
	}
}

func TestValidateUIContractSuccess(t *testing.T) {
	input := `{"namespace":"test","box":{"type":"panel","controls":[{"lbl":{"type":"label"}}]}}`
	resp, err := validateUIContract(ValidateUIContractArgs{UIJSON: input})
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, "Success") {
		t.Errorf("expected success message, got:\n%s", text)
	}
}

func TestDesignDocResource(t *testing.T) {
	resp, err := designDocResource()
	if err != nil {
		t.Fatalf("unexpected resource error: %v", err)
	}
	if resp == nil || len(resp.Contents) == 0 {
		t.Fatal("expected non-empty resource content")
	}

	text := resp.Contents[0].TextResourceContents.Text
	if !strings.Contains(text, "Server Form Hijack") {
		t.Error("design doc missing 'Server Form Hijack' architecture mention")
	}
	if !strings.Contains(text, "ui/_ui_defs.json") {
		t.Error("design doc missing ui_defs registration step")
	}
}

func TestDesignDocURIConstant(t *testing.T) {
	if bedrockui.DesignDocURI != "bedrock://docs/custom_json_ui_design" {
		t.Errorf("unexpected DesignDocURI: %s", bedrockui.DesignDocURI)
	}
}

func TestResourceHandlersReturnContent(t *testing.T) {
	tests := []struct {
		name     string
		handler  func() (*mcp.ResourceResponse, error)
		expected string
	}{
		{"ui_reference", uiReferenceResource, "MINECRAFT BEDROCK JSON UI REFERENCE"},
		{"full_reference", fullReferenceResource, "ELEMENT TYPES"},
		{"design_doc", designDocResource, "Implementation Steps"},
		{"leaf_button_templates", leafButtonTemplatesResource, "LEAF BUTTON TEMPLATES"},
		{"leaf_dialog_layouts", leafDialogLayoutsResource, "LEAF DIALOG LAYOUTS"},
		{"binding_patterns", bindingPatternsResource, "BINDING PATTERNS"},
		{"custom_form_controls", customFormControlsResource, "CUSTOM FORM CONTROLS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.handler()
			if err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if resp == nil || len(resp.Contents) == 0 {
				t.Fatal("expected resource content")
			}
			if !strings.Contains(resp.Contents[0].TextResourceContents.Text, tt.expected) {
				t.Errorf("expected content containing %q", tt.expected)
			}
		})
	}
}

// --- New tool handler tests ---

func TestGenerateActionForm(t *testing.T) {
	resp, err := generateActionForm(GenerateActionFormArgs{Namespace: "test_ui", Variant: "default"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"namespace": "test_ui"`) {
		t.Error("generated action form missing expected namespace")
	}
	if !strings.Contains(text, "third_party_server_screen") {
		t.Error("generated action form missing screen registration")
	}
}

func TestGenerateModalForm(t *testing.T) {
	resp, err := generateModalForm(GenerateModalFormArgs{
		Namespace: "test_ui",
		Controls: []bedrockui.ModalFormControl{
			{Type: "label", Label: "Hello"},
			{Type: "toggle", Label: "Enable"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"namespace": "test_ui"`) {
		t.Error("generated modal form missing expected namespace")
	}
	if !strings.Contains(text, "custom_label") {
		t.Error("generated modal form missing custom_label")
	}
	if !strings.Contains(text, "custom_toggle") {
		t.Error("generated modal form missing custom_toggle")
	}
}

func TestGenerateTabbedForm(t *testing.T) {
	resp, err := generateTabbedForm(GenerateTabbedFormArgs{
		Namespace: "test_ui",
		TabNames:  []string{"Shop", "Stats", "Settings"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `§t§a§bShop`) {
		t.Error("generated tabbed form missing tab marker")
	}
}

func TestGenerateGridForm(t *testing.T) {
	resp, err := generateGridForm(GenerateGridFormArgs{
		Namespace:   "test_ui",
		Columns:     3,
		Rows:        4,
		IconSupport: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"type": "grid"`) {
		t.Error("generated grid form missing grid type")
	}
}

func TestCherryThemePackPack(t *testing.T) {
	resp, err := cherryThemePack(CherryThemePackArgs{
		Mode:  "pack",
		Title: "My Title",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.HasPrefix(text, "\u00a7f\u00a70\u00a70") {
		t.Error("packed string missing §f§0§0 prefix")
	}
}

func TestCherryThemePackUnpack(t *testing.T) {
	packed, _ := bedrockui.CherryThemePack(bedrockui.CherryThemePackParams{Mode: "pack", Title: "My Title", Header: "textures/header"})
	resp, err := cherryThemePack(CherryThemePackArgs{
		Mode:         "unpack",
		PackedString: packed,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"title":`) {
		t.Error("unpacked result missing title field")
	}
}

func TestAnalyzeBindings(t *testing.T) {
	input := `{"namespace":"test","btn":{"type":"button","controls":[{"lbl":{"type":"label","bindings":[{"binding_name":"#text","binding_type":"global"}]}}]}}`
	resp, err := analyzeBindings(AnalyzeBindingsArgs{UIJSON: input})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"total_bindings": 1`) {
		t.Error("analyze_bindings missing expected total_bindings count")
	}
}

func TestExpandInheritance(t *testing.T) {
	input := `{"namespace":"test","parent":{"type":"panel","size":[100,100]},"child@parent":{"type":"button"}}`
	resp, err := expandInheritance(ExpandInheritanceArgs{UIJSON: input})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	text := resp.Content[0].TextContent.Text
	if !strings.Contains(text, `"namespace": "test"`) {
		t.Error("expand_inheritance missing expected namespace")
	}
}
