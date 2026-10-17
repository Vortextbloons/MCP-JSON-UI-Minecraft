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
	// DesignDocURI is the MCP resource URI for the custom JSON UI implementation design document.
	DesignDocURI = "bedrock://docs/custom_json_ui_design"
	// LeafButtonTemplatesURI is the MCP resource URI for Leaf button template patterns.
	LeafButtonTemplatesURI = "bedrock://docs/leaf_button_templates"
	// LeafDialogLayoutsURI is the MCP resource URI for Leaf dialog layout patterns.
	LeafDialogLayoutsURI = "bedrock://docs/leaf_dialog_layouts"
	// BindingPatternsURI is the MCP resource URI for binding pattern examples.
	BindingPatternsURI = "bedrock://docs/binding_patterns"
	// CustomFormControlsURI is the MCP resource URI for modal form control definitions.
	CustomFormControlsURI = "bedrock://docs/custom_form_controls"
	// ReferenceMIMEType is the MIME type for the reference resources.
	ReferenceMIMEType = "text/plain"
)

//go:embed full_reference.txt
var fullReferenceContent string

//go:embed leaf_button_templates.txt
var leafButtonTemplatesContent string

//go:embed leaf_dialog_layouts.txt
var leafDialogLayoutsContent string

//go:embed binding_patterns.txt
var bindingPatternsContent string

//go:embed custom_form_controls.txt
var customFormControlsContent string

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

// DesignDoc returns the custom JSON UI implementation design document.
func DesignDoc() string {
	return `=== CUSTOM JSON UI IMPLEMENTATION GUIDE ===

Architecture: "Server Form Hijack"
- Behavior Pack: Trigger a Server Form (ActionFormData) via @minecraft/server-ui API
- Resource Pack: Override server_form.json to inject your custom layout

Components:
  1. ui/_ui_defs.json          — Register your .json files for loading
  2. ui/my_custom_ui.json       — Your custom layout with a namespace
  3. ui/server_form.json        — Override of the vanilla form screen
  4. scripts/main.js            — Trigger script using @minecraft/server-ui

Implementation Steps:

Step 1 — Create Custom Layout (ui/my_custom_ui.json):
  {
    "namespace": "custom_ui_namespace",
    "my_main_screen": {
      "type": "panel",
      "controls": [ /* buttons, images, text */ ]
    }
  }

Step 2 — Register UI Files (ui/_ui_defs.json):
  { "ui_defs": ["ui/my_custom_ui.json", "ui/server_form.json"] }

Step 3 — Override server_form.json:
  For a minimal global replacement, override long_form (this affects all action forms):
  {
    "namespace": "server_form",
    "long_form": {
      "type": "panel",
      "controls": [
        { "custom_injection@custom_ui_namespace.my_main_screen": {} }
      ]
    }
  }
  Use a named #title_text binding plus a view expression targeting #visible
  to show your panel only when form title matches a secret string.

Step 4 — Trigger Script (scripts/main.js):
  - Subscribe to world.afterEvents.itemUse
  - Check event.itemStack.typeId === "minecraft:compass"
  - Call new ActionFormData().title("Trigger Custom UI").body("Loading...").button("Close").show(player)

Best Practices:
  - Use percentage-based sizing on root panels (e.g. "size": ["100%", "100%"])
  - Set higher layer on interactive elements than background panels
  - List @minecraft/server and @minecraft/server-ui in BP manifest dependencies
  - Named child definitions go inside controls: [{"child_name": {"type": "label"}}].
  - Empty containers and containers populated by inheritance or factories are valid.
  - View bindings use source_property_name and target_property_name; collection_details supplies item context without binding_name. Named global/collection bindings use binding_name.
  - Define animations at the document root; reference @namespace.animation through anims or an animated property. $definitions is not an animation container.
  - Collection containers can have a factory property; a factory element is a separate mechanism.
  - binding_condition controls when data is refreshed. Use a view expression targeting #visible to select a form by title, and hide the vanilla layout for the same title.
  - A new namespace does not register a screen. Inject it into server_form.main_screen_content through a factory named server_form_factory, or override the existing server_form screen.
  - JSON UI belongs in the resource pack. Script form fields, labels, defaults and options come from the behavior pack's server-ui call.
  - Passing the local schema does not establish engine compatibility. Test on the target Bedrock version with content logs enabled.

Sources:
  https://github.com/Mojang/bedrock-samples/blob/main/resource_pack/ui/server_form.json
  https://wiki.bedrock.dev/json-ui/json-ui-intro
  https://wiki.bedrock.dev/json-ui/json-ui-documentation
  https://wiki.bedrock.dev/json-ui/modifying-server-forms
  Leaf/StarLib and Cherry conventions are pack-specific dependencies, not engine APIs.`
}

// LeafButtonTemplates returns the Leaf button templates reference document.
func LeafButtonTemplates() string {
	return leafButtonTemplatesContent
}

// LeafDialogLayouts returns the Leaf dialog layouts reference document.
func LeafDialogLayouts() string {
	return leafDialogLayoutsContent
}

// BindingPatterns returns the binding patterns reference document.
func BindingPatterns() string {
	return bindingPatternsContent
}

// CustomFormControls returns the custom form controls reference document.
func CustomFormControls() string {
	return customFormControlsContent
}
