package mcpserver

import (
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"mcp-server/internal/bedrockui"
)

const (
	serverName    = "minecraft-bedrock-json-ui"
	serverVersion = "2.0.0"
	serverBlurb   = "Read bedrock://docs/ui_reference (quick), bedrock://docs/json_ui_full_reference (full), bedrock://docs/custom_json_ui_design (implementation guide), bedrock://docs/leaf_button_templates, bedrock://docs/leaf_dialog_layouts, bedrock://docs/binding_patterns, or bedrock://docs/custom_form_controls before authoring JSON UI. Validate each draft with validate_ui_contract, generate forms with generate_action_form/generate_modal_form/generate_tabbed_form/generate_grid_form, analyze bindings with analyze_bindings, expand inheritance with expand_inheritance, and pack Cherry UI themes with cherry_theme_pack."
)

// Tool and resource names/descriptions
var (
	// Existing
	validateToolName        = "validate_ui_contract"
	validateToolDescription = "Validates Minecraft Bedrock JSON UI strings AND returns an implementation guide for building custom JSON UI via the Server Form Hijack pattern. Always read the guide before authoring."

	resourceName        = "ui_reference"
	resourceBlurb       = "Curated Bedrock JSON UI reference derived from the design and technical docs."
	fullResourceName    = "json_ui_full_reference"
	fullResourceBlurb   = "Restructured, comprehensive Bedrock JSON UI documentation covering all element types, properties, hardcoded values, and more."
	designResourceName  = "custom_json_ui_design"
	designResourceBlurb = "Complete implementation guide for custom JSON UI via the Server Form Hijack pattern, including architecture, components, steps, and best practices."

	// New resources
	leafButtonTemplatesResourceName   = "leaf_button_templates"
	leafButtonTemplatesResourceBlurb  = "Leaf Essentials button template patterns including light_text_button, light_content_button, light_glyph_button, and StarLib2 styles."
	leafDialogLayoutsResourceName     = "leaf_dialog_layouts"
	leafDialogLayoutsResourceBlurb    = "Leaf Essentials dialog layout patterns including main_panel_no_buttons, main_panel_one_button, main_panel_two_buttons, and background types."
	bindingPatternsResourceName       = "binding_patterns"
	bindingPatternsResourceBlurb      = "Complete binding pattern reference: collection, view, global bindings with string operation syntax and factory wiring."
	customFormControlsResourceName    = "custom_form_controls"
	customFormControlsResourceBlurb   = "Modal form control definitions: label, toggle, slider, step_slider, dropdown, input with their binding names."

	// New tools
	generateActionFormToolName        = "generate_action_form"
	generateActionFormToolDescription = "Generates a complete Minecraft Bedrock action form JSON UI with factory-based dynamic buttons, optional icon support, scrolling, and layout variants (default, fullscreen, grid, tabs, light)."
	generateModalFormToolName         = "generate_modal_form"
	generateModalFormToolDescription  = "Generates a complete Minecraft Bedrock modal form JSON UI with control_ids factory mapping for label, toggle, slider, step_slider, dropdown, and input controls."
	generateTabbedFormToolName        = "generate_tabbed_form"
	generateTabbedFormToolDescription = "Generates a tabbed action form JSON UI with tab-switching logic using §t§a§b markers and visibility bindings."
	generateGridFormToolName          = "generate_grid_form"
	generateGridFormToolDescription   = "Generates a grid-based action form JSON UI with configurable columns, rows, and optional icon support."
	cherryThemePackToolName           = "cherry_theme_pack"
	cherryThemePackToolDescription    = "Packs or unpacks Cherry UI / Nut UI theme strings using the Leaf Essentials string-packing format (§f§0§0 prefix with fixed-width segments)."
	analyzeBindingsToolName           = "analyze_bindings"
	analyzeBindingsToolDescription    = "Analyzes a JSON UI document and extracts all bindings with their types, collections, paths, and warnings."
	expandInheritanceToolName         = "expand_inheritance"
	expandInheritanceToolDescription  = "Resolves @namespace.control inheritance syntax in a JSON UI document and returns flattened control definitions with merged properties."
)

// New constructs the MCP server and registers all tools and resources.
func New() (*mcp.Server, error) {
	server := mcp.NewServer(
		stdio.NewStdioServerTransport(),
		mcp.WithName(serverName),
		mcp.WithVersion(serverVersion),
		mcp.WithInstructions(serverBlurb),
	)

	// Resources
	if err := server.RegisterResource(bedrockui.ReferenceURI, resourceName, resourceBlurb, bedrockui.ReferenceMIMEType, uiReferenceResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.FullReferenceURI, fullResourceName, fullResourceBlurb, bedrockui.ReferenceMIMEType, fullReferenceResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.DesignDocURI, designResourceName, designResourceBlurb, bedrockui.ReferenceMIMEType, designDocResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.LeafButtonTemplatesURI, leafButtonTemplatesResourceName, leafButtonTemplatesResourceBlurb, bedrockui.ReferenceMIMEType, leafButtonTemplatesResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.LeafDialogLayoutsURI, leafDialogLayoutsResourceName, leafDialogLayoutsResourceBlurb, bedrockui.ReferenceMIMEType, leafDialogLayoutsResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.BindingPatternsURI, bindingPatternsResourceName, bindingPatternsResourceBlurb, bedrockui.ReferenceMIMEType, bindingPatternsResource); err != nil {
		return nil, err
	}
	if err := server.RegisterResource(bedrockui.CustomFormControlsURI, customFormControlsResourceName, customFormControlsResourceBlurb, bedrockui.ReferenceMIMEType, customFormControlsResource); err != nil {
		return nil, err
	}

	// Tools
	if err := server.RegisterTool(validateToolName, validateToolDescription, validateUIContract); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(generateActionFormToolName, generateActionFormToolDescription, generateActionForm); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(generateModalFormToolName, generateModalFormToolDescription, generateModalForm); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(generateTabbedFormToolName, generateTabbedFormToolDescription, generateTabbedForm); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(generateGridFormToolName, generateGridFormToolDescription, generateGridForm); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(cherryThemePackToolName, cherryThemePackToolDescription, cherryThemePack); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(analyzeBindingsToolName, analyzeBindingsToolDescription, analyzeBindings); err != nil {
		return nil, err
	}
	if err := server.RegisterTool(expandInheritanceToolName, expandInheritanceToolDescription, expandInheritance); err != nil {
		return nil, err
	}

	return server, nil
}
