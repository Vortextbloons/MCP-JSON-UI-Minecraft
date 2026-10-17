package mcpserver

import (
	"encoding/json"
	"fmt"
	"strings"

	"bedrock-json-ui-mcp/internal/bedrockui"
	mcp "github.com/metoro-io/mcp-golang"
)

const validatePreface = `📚 Resources: bedrock://docs/ui_reference (quick) | bedrock://docs/json_ui_full_reference (full) | bedrock://docs/custom_json_ui_design (implementation guide)

💡 Validation Hints:
  - Root namespace is required
  - controls contains named child objects; empty and inherited containers are valid
  - View bindings use source_property_name and target_property_name
  - Collection-details bindings do not require binding_name
  - Declare animations at the document root and reference them via anims or animated properties
  - Collection containers can use a factory property
  - This is a partial structural check; verify templates, bindings, and input behavior in Minecraft`

// ValidateUIContractArgs defines the tool input payload.
type ValidateUIContractArgs struct {
	UIJSON string `json:"ui_json" jsonschema:"required,description=The fully stringified Minecraft Bedrock JSON UI document to validate."`
}

func validateUIContract(args ValidateUIContractArgs) (*mcp.ToolResponse, error) {
	var response strings.Builder
	response.WriteString(validatePreface)
	response.WriteString("\n\n---\n\n")

	errors, err := bedrockui.ValidateDocument(args.UIJSON)
	if err != nil {
		response.WriteString(fmt.Sprintf("Internal Validation Error: %v", err))
		return mcp.NewToolResponse(mcp.NewTextContent(response.String())), nil
	}

	if len(errors) == 0 {
		response.WriteString("✅ Success: UI code passes the local partial schema. Runtime behavior and external references are not verified.")
		return mcp.NewToolResponse(mcp.NewTextContent(response.String())), nil
	}

	response.WriteString("❌ Validation Failed with the following errors:\n")
	for i, msg := range errors {
		fmt.Fprintf(&response, "%d. %s\n", i+1, msg)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(response.String())), nil
}

// --- Resources ---

func uiReferenceResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.ReferenceURI,
			bedrockui.UIReference(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func fullReferenceResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.FullReferenceURI,
			bedrockui.FullReference(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func designDocResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.DesignDocURI,
			bedrockui.DesignDoc(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func leafButtonTemplatesResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.LeafButtonTemplatesURI,
			bedrockui.LeafButtonTemplates(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func leafDialogLayoutsResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.LeafDialogLayoutsURI,
			bedrockui.LeafDialogLayouts(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func bindingPatternsResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.BindingPatternsURI,
			bedrockui.BindingPatterns(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

func customFormControlsResource() (*mcp.ResourceResponse, error) {
	return mcp.NewResourceResponse(
		mcp.NewTextEmbeddedResource(
			bedrockui.CustomFormControlsURI,
			bedrockui.CustomFormControls(),
			bedrockui.ReferenceMIMEType,
		),
	), nil
}

// --- Tool: generate_action_form ---

type GenerateActionFormArgs struct {
	Namespace     string `json:"namespace" jsonschema:"description=UI namespace to use (default: custom_ui)"`
	Variant       string `json:"variant" jsonschema:"description=Layout variant: default, fullscreen, grid, tabs, light (default: default)"`
	IconSupport   bool   `json:"icon_support" jsonschema:"description=Whether to include icon texture bindings for each button (default: false)"`
	IncludeScroll bool   `json:"include_scroll" jsonschema:"description=Whether to wrap buttons in a scrolling_panel (default: false)"`
}

func generateActionForm(args GenerateActionFormArgs) (*mcp.ToolResponse, error) {
	jsonStr := bedrockui.GenerateActionForm(bedrockui.ActionFormParams{
		Namespace:     args.Namespace,
		Variant:       args.Variant,
		IconSupport:   args.IconSupport,
		IncludeScroll: args.IncludeScroll,
	})
	return generatedFormResponse(jsonStr), nil
}

// --- Tool: generate_modal_form ---

type GenerateModalFormArgs struct {
	Namespace string                       `json:"namespace" jsonschema:"description=UI namespace to use (default: custom_ui)"`
	Controls  []bedrockui.ModalFormControl `json:"controls" jsonschema:"required,description=Array of modal form control definitions."`
}

func generateModalForm(args GenerateModalFormArgs) (*mcp.ToolResponse, error) {
	jsonStr := bedrockui.GenerateModalForm(bedrockui.ModalFormParams{
		Namespace: args.Namespace,
		Controls:  args.Controls,
	})
	return generatedFormResponse(jsonStr), nil
}

// --- Tool: generate_tabbed_form ---

type GenerateTabbedFormArgs struct {
	Namespace string   `json:"namespace" jsonschema:"description=UI namespace to use (default: custom_ui)"`
	TabNames  []string `json:"tab_names" jsonschema:"required,description=List of tab names."`
}

func generateTabbedForm(args GenerateTabbedFormArgs) (*mcp.ToolResponse, error) {
	jsonStr := bedrockui.GenerateTabbedForm(bedrockui.TabbedFormParams{
		Namespace: args.Namespace,
		TabNames:  args.TabNames,
	})
	return generatedFormResponse(jsonStr), nil
}

// --- Tool: generate_grid_form ---

type GenerateGridFormArgs struct {
	Namespace   string `json:"namespace" jsonschema:"description=UI namespace to use (default: custom_ui)"`
	Columns     int    `json:"columns" jsonschema:"description=Number of grid columns (default: 3)"`
	Rows        int    `json:"rows" jsonschema:"description=Number of grid rows (default: 4)"`
	IconSupport bool   `json:"icon_support" jsonschema:"description=Whether to include icon texture bindings (default: false)"`
}

func generateGridForm(args GenerateGridFormArgs) (*mcp.ToolResponse, error) {
	jsonStr := bedrockui.GenerateGridForm(bedrockui.GridFormParams{
		Namespace:   args.Namespace,
		Columns:     args.Columns,
		Rows:        args.Rows,
		IconSupport: args.IconSupport,
	})
	return generatedFormResponse(jsonStr), nil
}

// --- Tool: cherry_theme_pack ---

type CherryThemePackArgs struct {
	Mode         string `json:"mode" jsonschema:"required,description=Either 'pack' or 'unpack'."`
	Title        string `json:"title" jsonschema:"description=Title text (pack mode)."`
	Header       string `json:"header" jsonschema:"description=Header texture path (pack mode)."`
	Border       string `json:"border" jsonschema:"description=Border texture path (pack mode)."`
	HButton      string `json:"h_button" jsonschema:"description=Header button texture (pack mode)."`
	DeButton     string `json:"de_button" jsonschema:"description=Default button texture (pack mode)."`
	DiButton     string `json:"di_button" jsonschema:"description=Disabled button texture (pack mode)."`
	HoButton     string `json:"ho_button" jsonschema:"description=Hover button texture (pack mode)."`
	Paperdoll    string `json:"paperdoll" jsonschema:"description=Paperdoll texture path (pack mode)."`
	PackedString string `json:"packed_string" jsonschema:"description=Packed string to decode (unpack mode)."`
}

func cherryThemePack(args CherryThemePackArgs) (*mcp.ToolResponse, error) {
	result, err := bedrockui.CherryThemePack(bedrockui.CherryThemePackParams{
		Mode:         args.Mode,
		Title:        args.Title,
		Header:       args.Header,
		Border:       args.Border,
		HButton:      args.HButton,
		DeButton:     args.DeButton,
		DiButton:     args.DiButton,
		HoButton:     args.HoButton,
		Paperdoll:    args.Paperdoll,
		PackedString: args.PackedString,
	})
	if err != nil {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Error: %v", err))), nil
	}
	return mcp.NewToolResponse(mcp.NewTextContent(result)), nil
}

// --- Tool: analyze_bindings ---

type AnalyzeBindingsArgs struct {
	UIJSON string `json:"ui_json" jsonschema:"required,description=The JSON UI document to analyze for bindings."`
}

func analyzeBindings(args AnalyzeBindingsArgs) (*mcp.ToolResponse, error) {
	result, err := bedrockui.AnalyzeBindings(args.UIJSON)
	if err != nil {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Error: %v", err))), nil
	}
	b, _ := json.MarshalIndent(result, "", "    ")
	return mcp.NewToolResponse(mcp.NewTextContent(string(b))), nil
}

// --- Tool: expand_inheritance ---

type ExpandInheritanceArgs struct {
	UIJSON string `json:"ui_json" jsonschema:"required,description=The JSON UI document to expand @ inheritance for."`
}

func expandInheritance(args ExpandInheritanceArgs) (*mcp.ToolResponse, error) {
	result, err := bedrockui.ExpandInheritance(args.UIJSON)
	if err != nil {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Error: %v", err))), nil
	}
	b, _ := json.MarshalIndent(result, "", "    ")
	return mcp.NewToolResponse(mcp.NewTextContent(string(b))), nil
}

// generatedFormResponse keeps the layout JSON separate from installation guidance.
func generatedFormResponse(raw string) *mcp.ToolResponse {
	var doc struct {
		Namespace string `json:"namespace"`
	}
	_ = json.Unmarshal([]byte(raw), &doc)
	patch, _ := json.MarshalIndent(map[string]any{
		"namespace":                 "server_form",
		"third_party_server_screen": map[string]any{"$screen_content": doc.Namespace + ".main_content"},
	}, "", "  ")
	return mcp.NewToolResponse(mcp.NewTextContent(raw), mcp.NewTextContent(
		"Installation: save the layout in resource_pack/ui and register its path in ui/_ui_defs.json. "+
			"Merge this companion override into resource_pack/ui/server_form.json (namespace server_form), and register that file too:\n"+string(patch)+
			"\nThis replaces the layout for all forms of the generated type. Add title-based routing to scope it. "+
			"Keep a distinct custom namespace; do not use server_form as the generated namespace. "+
			"Labels, values, options, defaults and form response handling must be defined in the behavior pack using server-ui. "+
			"Tabbed output is a layout sketch requiring script-driven page selection; markers alone do not switch tabs. "+
			"Verify inherited templates and behavior against your target game's vanilla files."))
}
