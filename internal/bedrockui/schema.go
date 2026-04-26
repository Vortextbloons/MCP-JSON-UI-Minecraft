package bedrockui

import (
	"strconv"
	"strings"
)

const SchemaURI = "https://bedrock.dev/schema/ui/bedrock_ui.json"

const schemaTemplate = `{
  "$id": "https://bedrock.dev/schema/ui/bedrock_ui.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Minecraft Bedrock JSON UI Document",
  "type": "object",
  "required": ["namespace"],
  "properties": {
    "namespace": {
      "type": "string",
      "minLength": 1,
      "description": "Required UI namespace for resource packs"
    }
  },
  "patternProperties": {
    "__ROOT_PATTERN__": {
      "$ref": "#/$defs/element"
    }
  },
  "additionalProperties": false,
  "$defs": {
    "vector2": {
      "type": "array",
      "minItems": 2,
      "maxItems": 2,
      "items": {
        "anyOf": [
          {"type": "string"},
          {"type": "number"}
        ]
      }
    },
    "vector4": {
      "type": "array",
      "minItems": 4,
      "maxItems": 4,
      "items": {
        "anyOf": [
          {"type": "string"},
          {"type": "number"}
        ]
      }
    },
    "anchor": {
      "type": "string",
      "enum": [
        "top_left",
        "top_middle",
        "top_right",
        "left_middle",
        "center",
        "right_middle",
        "bottom_left",
        "bottom_middle",
        "bottom_right"
      ]
    },
    "binding": {
      "type": "object",
      "properties": {
        "ignored": {"type": "boolean"},
        "binding_type": {
          "type": "string",
          "enum": ["global", "view", "collection", "collection_details", "none"]
        },
        "binding_name": {"type": "string"},
        "binding_name_override": {"type": "string"},
        "binding_collection_name": {"type": "string"},
        "binding_collection_prefix": {"type": "string"},
        "binding_condition": {
          "type": "string",
          "enum": ["always", "always_when_visible", "visible", "once", "none", "visibility_changed"]
        },
        "source_control_name": {"type": "string"},
        "source_property_name": {"type": "string"},
        "target_property_name": {"type": "string"},
        "resolve_sibling_scope": {"type": "boolean"}
      },
      "required": ["binding_name"],
      "additionalProperties": true
    },
    "element": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": [
            "panel",
            "stack_panel",
            "collection_panel",
            "grid",
            "label",
            "image",
            "input_panel",
            "button",
            "toggle",
            "dropdown",
            "slider",
            "slider_box",
            "edit_box",
            "scroll_view",
            "scrollbar_track",
            "scrollbar_box",
            "factory",
            "screen",
            "custom",
            "selection_wheel",
            "tab",
            "carousel_label",
            "grid_item",
            "scrollbar"
          ]
        },
        "visible": {"type": "boolean"},
        "enabled": {"type": "boolean"},
        "layer": {"type": "integer"},
        "alpha": {"type": "number"},
        "propagate_alpha": {"type": "boolean"},
        "clips_children": {"type": "boolean"},
        "allow_clipping": {"type": "boolean"},
        "clip_offset": {"$ref": "#/$defs/vector2"},
        "selected": {"type": "boolean"},
        "use_child_anchors": {"type": "boolean"},
        "controls": {
          "type": "array",
          "items": {"$ref": "#/$defs/element"}
        },
        "bindings": {
          "type": "array",
          "items": {"$ref": "#/$defs/binding"}
        },
        "size": {"$ref": "#/$defs/vector2"},
        "max_size": {"$ref": "#/$defs/vector2"},
        "min_size": {"$ref": "#/$defs/vector2"},
        "offset": {"$ref": "#/$defs/vector2"},
        "anchor_from": {"$ref": "#/$defs/anchor"},
        "anchor_to": {"$ref": "#/$defs/anchor"},
        "inherit_max_sibling_width": {"type": "boolean"},
        "inherit_max_sibling_height": {"type": "boolean"},
        "use_anchored_offset": {"type": "boolean"},
        "contained": {"type": "boolean"},
        "draggable": {
          "type": "string",
          "enum": ["vertical", "horizontal", "both"]
        },
        "follows_cursor": {"type": "boolean"},
        "orientation": {
          "type": "string",
          "enum": ["vertical", "horizontal"]
        },
        "grid_dimensions": {"$ref": "#/$defs/vector2"},
        "maximum_grid_items": {"type": "integer"},
        "grid_dimension_binding": {"type": "string"},
        "grid_rescaling_type": {
          "type": "string",
          "enum": ["vertical", "horizontal", "none"]
        },
        "grid_fill_direction": {
          "type": "string",
          "enum": ["vertical", "horizontal", "none"]
        },
        "grid_item_template": {"type": "string"},
        "text": {"type": "string"},
        "color": {"$ref": "#/$defs/vector4"},
        "locked_color": {"$ref": "#/$defs/vector4"},
        "shadow": {"type": "boolean"},
        "hide_hyphen": {"type": "boolean"},
        "notify_on_ellipses": {
          "type": "array",
          "items": {"type": "string"}
        },
        "enable_profanity_filter": {"type": "boolean"},
        "locked_alpha": {"type": "number"},
        "font_size": {
          "type": "string",
          "enum": ["small", "normal", "large", "extra_large"]
        },
        "font_scale_factor": {"type": "number"},
        "localize": {"type": "boolean"},
        "line_padding": {"type": "number"},
        "font_type": {"type": "string"},
        "backup_font_type": {"type": "string"},
        "text_alignment": {"type": "string"},
        "texture": {"type": "string"},
        "allow_debug_missing_texture": {"type": "boolean"},
        "uv": {"$ref": "#/$defs/vector2"},
        "uv_size": {"$ref": "#/$defs/vector2"},
        "texture_file_system": {"type": "string"},
        "nineslice_size": {
          "oneOf": [
            {"type": "integer"},
            {"$ref": "#/$defs/vector4"}
          ]
        },
        "tiled": {
          "oneOf": [
            {"type": "boolean"},
            {"type": "string", "enum": ["x", "y"]}
          ]
        },
        "tiled_scale": {"$ref": "#/$defs/vector2"},
        "clip_direction": {
          "type": "string",
          "enum": ["left", "right", "up", "down", "center"]
        },
        "clip_ratio": {"type": "number"},
        "clip_pixelperfect": {"type": "boolean"},
        "keep_ratio": {"type": "boolean"},
        "bilinear": {"type": "boolean"},
        "fill": {"type": "boolean"},
        "grayscale": {"type": "boolean"},
        "force_texture_reload": {"type": "boolean"},
        "base_size": {"$ref": "#/$defs/vector2"},
        "button_mappings": {
          "type": "array",
          "items": {"type": "object"}
        },
        "modal": {"type": "boolean"},
        "inline_modal": {"type": "boolean"},
        "always_listen_to_input": {"type": "boolean"},
        "always_handle_pointer": {"type": "boolean"},
        "always_handle_controller_direction": {"type": "boolean"},
        "hover_enabled": {"type": "boolean"},
        "prevent_touch_input": {"type": "boolean"},
        "consume_event": {"type": "boolean"},
        "consume_hover_events": {"type": "boolean"},
        "gesture_tracking_button": {"type": "string"},
        "default_focus_precedence": {"type": "integer"},
        "focus_enabled": {"type": "boolean"},
        "focus_wrap_enabled": {"type": "boolean"},
        "focus_magnet_enabled": {"type": "boolean"},
        "focus_identifier": {"type": "string"},
        "focus_change_down": {"type": "string"},
        "focus_change_up": {"type": "string"},
        "focus_change_left": {"type": "string"},
        "focus_change_right": {"type": "string"},
        "focus_container": {"type": "boolean"},
        "use_last_focus": {"type": "boolean"},
        "sound_name": {"type": "string"},
        "sound_volume": {"type": "number"},
        "sound_pitch": {"type": "number"},
        "sounds": {
          "type": "array",
          "items": {"type": "object"}
        },
        "text_box_name": {"type": "string"},
        "max_length": {"type": "integer"},
        "constrain_to_rect": {"type": "boolean"},
        "enabled_newline": {"type": "boolean"},
        "text_type": {
          "type": "string",
          "enum": ["ExtendedASCII", "IdentifierChars", "NumberChars"]
        },
        "text_control": {"type": "string"},
        "place_holder_control": {"type": "string"},
        "can_be_deselected": {"type": "boolean"},
        "always_listening": {"type": "boolean"},
        "slider_track_button": {"type": "string"},
        "slider_small_decrease_button": {"type": "string"},
        "slider_small_increase_button": {"type": "string"},
        "slider_steps": {"type": "integer"},
        "slider_direction": {
          "type": "string",
          "enum": ["vertical", "horizontal"]
        },
        "slider_timeout": {"type": "number"},
        "slider_collection_name": {"type": "string"},
        "slider_name": {"type": "string"},
        "slider_select_on_hover": {"type": "boolean"},
        "slider_selected_button": {"type": "string"},
        "slider_deselected_button": {"type": "string"},
        "slider_box_control": {"type": "string"},
        "background_control": {"type": "string"},
        "background_hover_control": {"type": "string"},
        "progress_control": {"type": "string"},
        "progress_hover_control": {"type": "string"},
        "dropdown_name": {"type": "string"},
        "dropdown_content_control": {"type": "string"},
        "dropdown_area": {"type": "string"},
        "toggle_name": {"type": "string"},
        "toggle_default_state": {"type": "boolean"},
        "radio_toggle_group": {"type": "boolean"},
        "toggle_group_forced_index": {"type": "integer"},
        "toggle_group_default_selected": {"type": "integer"},
        "reset_on_focus_lost": {"type": "boolean"},
        "toggle_on_hover": {"type": "string"},
        "toggle_on_button": {"type": "string"},
        "toggle_off_button": {"type": "string"},
        "enable_directional_toggling": {"type": "boolean"},
        "toggle_grid_collection_name": {"type": "string"},
        "checked_control": {"type": "string"},
        "unchecked_control": {"type": "string"},
        "checked_hover_control": {"type": "string"},
        "unchecked_hover_control": {"type": "string"},
        "checked_locked_control": {"type": "string"},
        "unchecked_locked_control": {"type": "string"},
        "checked_locked_hover_control": {"type": "string"},
        "unchecked_locked_hover_control": {"type": "string"},
        "property_bag": {"type": "object"},
        "anims": {
          "type": "array",
          "items": {"type": "string"}
        },
        "disable_anim_fast_forward": {"type": "boolean"},
        "animation_reset_name": {"type": "string"},
        "ignored": {"type": "boolean"},
        "variables": {
          "oneOf": [
            {"type": "object"},
            {"type": "array", "items": {"type": "object"}}
          ]
        },
        "modifications": {
          "type": "array",
          "items": {"type": "object"}
        },
        "grid_position": {"$ref": "#/$defs/vector2"},
        "collection_index": {"type": "integer"}
      },
      "additionalProperties": true
    }
  }
}`

// SchemaSource returns the JSON Schema used by the validator.
func SchemaSource() string {
	return strings.ReplaceAll(schemaTemplate, "__ROOT_PATTERN__", rootElementPattern())
}

func rootElementPattern() string {
	return identifierPatternExcept("namespace")
}

// identifierPatternExcept builds an identifier regex that excludes one exact key.
// It keeps the root schema permissive for control names while reserving namespace.
func identifierPatternExcept(exact string) string {
	runes := []rune(exact)
	parts := make([]string, 0, len(runes)+2)

	if len(runes) > 1 {
		parts = append(parts, "[A-Za-z0-9_]{1,"+strconv.Itoa(len(runes)-1)+"}")
	}
	parts = append(parts, "[A-Za-z0-9_]{"+strconv.Itoa(len(runes)+1)+",}")

	for i, r := range runes {
		prefix := string(runes[:i])
		suffixLen := len(runes) - i - 1
		if suffixLen == 0 {
			parts = append(parts, prefix+allowedClassExcept(r))
			continue
		}
		parts = append(parts, prefix+allowedClassExcept(r)+"[A-Za-z0-9_]{"+strconv.Itoa(suffixLen)+"}")
	}

	return "^(?:" + strings.Join(parts, "|") + ")$"
}

var allowedIdentifierRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_")

func allowedClassExcept(excluded rune) string {
	var b strings.Builder
	b.WriteByte('[')
	for _, r := range allowedIdentifierRunes {
		if r == excluded {
			continue
		}
		b.WriteRune(r)
	}
	b.WriteByte(']')
	return b.String()
}
