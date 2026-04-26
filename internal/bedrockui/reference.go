package bedrockui

import (
	_ "embed"
	"strings"
)

const (
	// ReferenceURI is the MCP resource URI exposed to clients (curated quick reference).
	ReferenceURI = "bedrock://docs/ui_reference"
	// FullReferenceURI is the MCP resource URI for the restructured full JSON UI docs.
	FullReferenceURI = "bedrock://docs/json_ui_full_reference"
	// ReferenceMIMEType is the MIME type for the reference resources.
	ReferenceMIMEType = "text/plain"
)

//go:embed full_reference.txt
var fullReferenceContent string

// UIReference returns a concise curated reference for the MCP resource.
func UIReference() string {
	sections := []string{
		"=== MINECRAFT BEDROCK JSON UI REFERENCE ===",
		"Core element types: panel, stack_panel, collection_panel, grid, label, image, input_panel, button, toggle, dropdown, slider, slider_box, edit_box, scroll_view, scrollbar_track, scrollbar_box, factory, screen, custom, selection_wheel.",
		"Legacy element types: tab, carousel_label, grid_item, scrollbar.",
		"Key control/layout props: visible, enabled, layer, alpha, propagate_alpha, clips_children, allow_clipping, clip_offset, controls, anims, ignored, variables, modifications, grid_position, collection_index.",
		"Key layout props: size, max_size, min_size, offset, anchor_from, anchor_to, inherit_max_sibling_width, inherit_max_sibling_height, draggable, follows_cursor.",
		"Key binding props: binding_type, binding_name, binding_name_override, binding_condition, source_control_name, source_property_name, target_property_name, binding_collection_name, binding_collection_prefix.",
		"Useful hardcoded bindings: #visible, #text, #hover_text, #item_id_aux, #item_custom_color, #toggle_state, #slider_value, #slider_steps, #property_field, #using_ellipses, #filtered_light_multiplier, #playername, #open, #progress_bar_current_amount, #progress_bar_total_amount.",
		"Workflow: read this resource, draft JSON, call validate_ui_contract, patch the JSON, and validate again before final output.",
	}

	return strings.Join(sections, "\n\n")
}

// FullReference returns the restructured, comprehensive JSON UI reference document.
func FullReference() string {
	return fullReferenceContent
}
