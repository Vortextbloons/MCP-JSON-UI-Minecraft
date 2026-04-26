package mcpserver

import (
	"encoding/json"
	"fmt"
	"strings"

	mcp "github.com/metoro-io/mcp-golang"
	"bedrock-json-ui-mcp/internal/bedrockui"
)

const validatePreface = `📚 Resources: bedrock://docs/ui_reference (quick) | bedrock://docs/json_ui_full_reference (full) | bedrock://docs/custom_json_ui_design (implementation guide)

💡 Validation Hints:
  - Root namespace is required
  - Container types (panel, stack_panel, grid, collection_panel) should include a controls array
  - Bindings require binding_name; use binding_name_override to target element properties
  - Animations belong in $definitions or referenced via anims array, not as root elements
  - Factory props only work on factory elements, not grids or panels
  - Use percentage-based sizing on root panels for cross-platform compatibility
  - Interactive elements need higher layer than backgrounds for click registration`

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
		response.WriteString("✅ Success: UI code passes Bedrock schema and logic constraints.")
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
	return mcp.NewToolResponse(mcp.NewTextContent(jsonStr)), nil
}

// --- Tool: generate_modal_form ---

type GenerateModalFormArgs struct {
	Namespace string                      `json:"namespace" jsonschema:"description=UI namespace to use (default: custom_ui)"`
	Controls  []bedrockui.ModalFormControl `json:"controls" jsonschema:"required,description=Array of modal form control definitions."`
}

func generateModalForm(args GenerateModalFormArgs) (*mcp.ToolResponse, error) {
	jsonStr := bedrockui.GenerateModalForm(bedrockui.ModalFormParams{
		Namespace: args.Namespace,
		Controls:  args.Controls,
	})
	return mcp.NewToolResponse(mcp.NewTextContent(jsonStr)), nil
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
	return mcp.NewToolResponse(mcp.NewTextContent(jsonStr)), nil
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
	return mcp.NewToolResponse(mcp.NewTextContent(jsonStr)), nil
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
