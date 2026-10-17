package bedrockui

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ActionFormParams configures action form generation.
type ActionFormParams struct {
	Namespace     string `json:"namespace"`
	Variant       string `json:"variant"` // default, fullscreen, grid, tabs, light
	IconSupport   bool   `json:"icon_support"`
	IncludeScroll bool   `json:"include_scroll"`
}

// ModalFormControl defines a single modal form control.
type ModalFormControl struct {
	Type        string   `json:"type"` // label, toggle, slider, step_slider, dropdown, input
	Label       string   `json:"label"`
	Placeholder string   `json:"placeholder"` // for input
	Default     any      `json:"default"`     // default value
	Min         int      `json:"min"`         // for slider
	Max         int      `json:"max"`         // for slider
	Step        int      `json:"step"`        // for step_slider
	Options     []string `json:"options"`     // for dropdown
}

// ModalFormParams configures modal form generation.
type ModalFormParams struct {
	Namespace string             `json:"namespace"`
	Controls  []ModalFormControl `json:"controls"`
}

// TabbedFormParams configures tabbed form generation.
type TabbedFormParams struct {
	Namespace string   `json:"namespace"`
	TabNames  []string `json:"tab_names"`
}

// GridFormParams configures grid form generation.
type GridFormParams struct {
	Namespace   string `json:"namespace"`
	Columns     int    `json:"columns"`
	Rows        int    `json:"rows"`
	IconSupport bool   `json:"icon_support"`
}

// CherryThemePackParams configures cherry theme string packing/unpacking.
type CherryThemePackParams struct {
	Mode         string `json:"mode"` // "pack" or "unpack"
	Title        string `json:"title"`
	Header       string `json:"header"`
	Border       string `json:"border"`
	HButton      string `json:"h_button"`  // header button
	DeButton     string `json:"de_button"` // default button
	DiButton     string `json:"di_button"` // disabled button
	HoButton     string `json:"ho_button"` // hover button
	Paperdoll    string `json:"paperdoll"`
	PackedString string `json:"packed_string"` // for unpack mode
}

// GenerateActionForm creates action form templates requiring server_form integration.
func GenerateActionForm(p ActionFormParams) string {
	if strings.TrimSpace(p.Namespace) == "" {
		p.Namespace = "custom_ui"
	}
	variant := p.Variant
	if variant == "" {
		variant = "default"
	}

	var root any
	switch variant {
	case "fullscreen":
		root = buildFullscreenForm(p)
	case "grid":
		return GenerateGridForm(GridFormParams{Namespace: p.Namespace, Columns: 3, Rows: 4, IconSupport: p.IconSupport})
	case "tabs":
		return GenerateTabbedForm(TabbedFormParams{Namespace: p.Namespace, TabNames: []string{"Tab 1", "Tab 2", "Tab 3"}})
	case "light":
		root = buildLightForm(p)
	default:
		root = buildDefaultForm(p)
	}

	return marshalGeneratedForm(root)
}

// GenerateModalForm creates modal widget templates; field data comes from server-ui.
func GenerateModalForm(p ModalFormParams) string {
	if strings.TrimSpace(p.Namespace) == "" {
		p.Namespace = "custom_ui"
	}

	// Inherit the vanilla input widgets so their visuals, focus, and response wiring
	// stay intact. Field data is supplied by ModalFormData, not these templates.
	root := map[string]any{
		"namespace": p.Namespace,
		"custom_form_screen@server_form.custom_form": map[string]any{
			"$child_control": p.Namespace + ".custom_form_panel",
		},
		"custom_form_panel@server_form.custom_form_panel": map[string]any{
			"$scrolling_content": p.Namespace + ".custom_form_scrolling_content",
		},
		"custom_form_scrolling_content@server_form.custom_form_scrolling_content": map[string]any{
			"modifications": []any{map[string]any{
				"array_name": "controls", "operation": "replace",
				"control_name": "generated_form",
				"value":        []any{map[string]any{"generated_form@" + p.Namespace + ".generated_contents": map[string]any{}}},
			}},
		},
	}
	ids := map[string]any{}
	for _, kind := range []string{"label", "toggle", "slider", "step_slider", "dropdown", "input", "header", "divider", "multiselect"} {
		parent := "server_form.custom_" + kind
		if kind == "divider" {
			parent = "settings_common.option_group_section_divider"
		}
		ids[kind] = "@" + parent
	}
	for _, c := range p.Controls {
		if _, ok := ids[c.Type]; !ok {
			continue
		}
		parent := strings.TrimPrefix(ids[c.Type].(string), "@")
		root["custom_"+c.Type+"@"+parent] = map[string]any{}
		ids[c.Type] = "@" + p.Namespace + ".custom_" + c.Type
	}
	root["generated_contents@server_form.generated_contents"] = map[string]any{
		"factory": map[string]any{"name": "buttons", "control_ids": ids},
	}
	root["main_content"] = map[string]any{
		"type": "panel",
		"factory": map[string]any{
			"name":        "server_form_factory",
			"control_ids": map[string]any{"long_form": "@server_form.long_form", "custom_form": "@" + p.Namespace + ".custom_form_screen"},
		},
	}
	return marshalGeneratedForm(root)
}

// GenerateTabbedForm creates a tabbed layout sketch; page selection requires script logic.
func GenerateTabbedForm(p TabbedFormParams) string {
	if strings.TrimSpace(p.Namespace) == "" {
		p.Namespace = "custom_ui"
	}

	tabButtons := make([]map[string]any, 0, len(p.TabNames))
	for i, name := range p.TabNames {
		tabButtons = append(tabButtons, map[string]any{
			fmt.Sprintf("tab_btn_%d@common_buttons.light_text_button", i): map[string]any{
				"size":                 []any{"fill", 24},
				"$pressed_button_name": "button.form_button_click",
				"$button_text":         fmt.Sprintf("§t§a§b%s", name),
				"collection_index":     i,
				"bindings": []map[string]any{
					{"binding_type": "collection_details", "binding_collection_name": "form_buttons"},
				},
			},
		})
	}

	root := map[string]any{
		"namespace": p.Namespace,
		"third_party_server_screen@common.base_screen": map[string]any{
			"$screen_content": fmt.Sprintf("%s.main_content", p.Namespace),
			"button_mappings": []map[string]any{
				{"from_button_id": "button.menu_cancel", "to_button_id": "button.menu_exit", "mapping_type": "global"},
			},
		},
		"main_content": map[string]any{
			"type": "factory",
			"control_ids": map[string]any{
				"long_form": fmt.Sprintf("@%s.long_form", p.Namespace),
			},
		},
		"long_form": map[string]any{
			"type": "panel",
			"size": []any{"100%", "100%"},
			"controls": []map[string]any{
				{
					"form@common_dialogs.main_panel_no_buttons": map[string]any{
						"size":                     []any{370, 240},
						"$text_name":               "#title_text",
						"$title_text_binding_type": "global",
						"$child_control":           fmt.Sprintf("%s.content_area", p.Namespace),
					},
				},
			},
		},
		"content_area": map[string]any{
			"type":        "stack_panel",
			"size":        []any{"100%", "100%"},
			"orientation": "vertical",
			"controls": []map[string]any{
				{
					"tab_bar": map[string]any{
						"type":        "stack_panel",
						"orientation": "horizontal",
						"size":        []any{"100%", 28},
						"controls":    tabButtons,
					},
				},
				{
					"scroll@common.scrolling_panel": map[string]any{
						"size":                   []any{"100%", "100% - 28px"},
						"$show_background":       false,
						"$scrolling_content":     fmt.Sprintf("%s.button_list", p.Namespace),
						"$scroll_size":           []any{5, "100% - 4px"},
						"$scrolling_pane_size":   []any{"100% - 4px", "100% - 2px"},
						"$scrolling_pane_offset": []any{2, 0},
					},
				},
			},
		},
		"button_list": map[string]any{
			"type":        "stack_panel",
			"size":        []any{"100%", "100%c"},
			"orientation": "vertical",
			"factory": map[string]any{
				"name":         "buttons",
				"control_name": fmt.Sprintf("%s.dynamic_button", p.Namespace),
			},
			"collection_name": "form_buttons",
			"bindings": []map[string]any{
				{"binding_name": "#form_button_length", "binding_name_override": "#collection_length"},
			},
		},
		"dynamic_button": map[string]any{
			"type":        "stack_panel",
			"size":        []any{"100%", 32},
			"orientation": "horizontal",
			"controls": []map[string]any{
				{
					"form_button@common_buttons.light_text_button": map[string]any{
						"size":                              []any{"fill", 32},
						"$pressed_button_name":              "button.form_button_click",
						"$button_text":                      "#form_button_text",
						"$button_text_binding_type":         "collection",
						"$button_text_grid_collection_name": "form_buttons",
						"bindings": []map[string]any{
							{"binding_type": "collection_details", "binding_collection_name": "form_buttons"},
						},
					},
				},
			},
		},
	}

	return marshalGeneratedForm(root)
}

// GenerateGridForm creates a grid-based action form JSON UI.
func GenerateGridForm(p GridFormParams) string {
	if strings.TrimSpace(p.Namespace) == "" {
		p.Namespace = "custom_ui"
	}
	if p.Columns <= 0 {
		p.Columns = 3
	}
	if p.Rows <= 0 {
		p.Rows = 4
	}

	root := map[string]any{
		"namespace": p.Namespace,
		"third_party_server_screen@common.base_screen": map[string]any{
			"$screen_content": fmt.Sprintf("%s.main_content", p.Namespace),
			"button_mappings": []map[string]any{
				{"from_button_id": "button.menu_cancel", "to_button_id": "button.menu_exit", "mapping_type": "global"},
			},
		},
		"main_content": map[string]any{
			"type": "factory",
			"control_ids": map[string]any{
				"long_form": fmt.Sprintf("@%s.long_form", p.Namespace),
			},
		},
		"long_form": map[string]any{
			"type": "panel",
			"size": []any{"100%", "100%"},
			"controls": []map[string]any{
				{
					"form@common_dialogs.main_panel_no_buttons": map[string]any{
						"size":                     []any{370, 240},
						"$text_name":               "#title_text",
						"$title_text_binding_type": "global",
						"$child_control":           fmt.Sprintf("%s.grid_area", p.Namespace),
					},
				},
			},
		},
		"grid_area": map[string]any{
			"type":                "grid",
			"grid_dimensions":     []any{p.Columns, p.Rows},
			"grid_item_template":  fmt.Sprintf("%s.dynamic_button", p.Namespace),
			"grid_rescaling_type": "horizontal",
			"factory": map[string]any{
				"name":         "buttons",
				"control_name": fmt.Sprintf("%s.dynamic_button", p.Namespace),
			},
			"collection_name": "form_buttons",
			"bindings": []map[string]any{
				{"binding_name": "#form_button_length", "binding_name_override": "#maximum_grid_items"},
			},
		},
		"dynamic_button": buildDynamicButton(p.Namespace, p.IconSupport),
	}

	return marshalGeneratedForm(root)
}

// CherryThemePack packs or unpacks a Cherry UI theme string.
func CherryThemePack(p CherryThemePackParams) (string, error) {
	if p.Mode == "unpack" {
		return cherryUnpack(p.PackedString)
	}
	if p.Mode == "pack" {
		return cherryPack(p)
	}
	return "", fmt.Errorf("invalid mode %q; use \"pack\" or \"unpack\"", p.Mode)
}

func cherryPack(p CherryThemePackParams) (string, error) {
	title := adjustTextLength(p.Title, 100)
	header := adjustTextLength(p.Header, 50)
	border := adjustTextLength(p.Border, 50)
	hButton := adjustTextLength(p.HButton, 50)
	deButton := adjustTextLength(p.DeButton, 50)
	diButton := adjustTextLength(p.DiButton, 50)
	hoButton := adjustTextLength(p.HoButton, 50)
	paperdoll := adjustTextLength(p.Paperdoll, 50)

	result := "\u00a7f\u00a70\u00a70" + // §f§0§0 prefix
		title + header + border + hButton + deButton + diButton + hoButton + paperdoll
	return result, nil
}

func cherryUnpack(s string) (string, error) {
	prefix := "\u00a7f\u00a70\u00a70"
	if len(s) < len(prefix) || s[:len(prefix)] != prefix {
		return "", fmt.Errorf("packed string missing §f§0§0 prefix")
	}
	data := s[len(prefix):]

	title := safeSlice(data, 0, 100)
	header := safeSlice(data, 100, 150)
	border := safeSlice(data, 150, 200)
	hButton := safeSlice(data, 200, 250)
	deButton := safeSlice(data, 250, 300)
	diButton := safeSlice(data, 300, 350)
	hoButton := safeSlice(data, 350, 400)
	paperdoll := safeSlice(data, 400, 450)

	result := map[string]string{
		"title":     strings.TrimSpace(title),
		"header":    strings.TrimSpace(header),
		"border":    strings.TrimSpace(border),
		"h_button":  strings.TrimSpace(hButton),
		"de_button": strings.TrimSpace(deButton),
		"di_button": strings.TrimSpace(diButton),
		"ho_button": strings.TrimSpace(hoButton),
		"paperdoll": strings.TrimSpace(paperdoll),
	}
	b, _ := json.MarshalIndent(result, "", "    ")
	return string(b), nil
}

func adjustTextLength(s string, length int) string {
	if len(s) > length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func safeSlice(s string, start, end int) string {
	if start > len(s) {
		start = len(s)
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		start = end
	}
	return s[start:end]
}

// --- Builders ---

func buildDefaultForm(p ActionFormParams) any {
	return buildStandardForm(p, "default")
}

func buildFullscreenForm(p ActionFormParams) any {
	return buildStandardForm(p, "fullscreen")
}

func buildLightForm(p ActionFormParams) any {
	return buildStandardForm(p, "light")
}

func buildStandardForm(p ActionFormParams, variant string) any {
	root := map[string]any{
		"namespace": p.Namespace,
		"third_party_server_screen@common.base_screen": map[string]any{
			"$screen_content": fmt.Sprintf("%s.main_content", p.Namespace),
			"button_mappings": []map[string]any{
				{"from_button_id": "button.menu_cancel", "to_button_id": "button.menu_exit", "mapping_type": "global"},
			},
		},
		"main_content": map[string]any{
			"type": "factory",
			"control_ids": map[string]any{
				"long_form": fmt.Sprintf("@%s.long_form", p.Namespace),
			},
		},
		"long_form": map[string]any{
			"type": "panel",
			"size": []any{"100%", "100%"},
			"controls": []map[string]any{
				{
					"form@common_dialogs.main_panel_no_buttons": map[string]any{
						"size":                     []any{370, 200},
						"$text_name":               "#title_text",
						"$title_text_binding_type": "global",
						"$child_control":           fmt.Sprintf("%s.button_panel", p.Namespace),
					},
				},
			},
		},
		"button_panel": map[string]any{
			"type":     "panel",
			"size":     []any{"100%", "100%"},
			"controls": buildButtonPanelControls(p),
		},
		"button_list": map[string]any{
			"type":        "stack_panel",
			"size":        []any{"100%", "100%c"},
			"orientation": "vertical",
			"factory": map[string]any{
				"name":         "buttons",
				"control_name": fmt.Sprintf("%s.dynamic_button", p.Namespace),
			},
			"collection_name": "form_buttons",
			"bindings": []map[string]any{
				{"binding_name": "#form_button_length", "binding_name_override": "#collection_length"},
			},
		},
		"dynamic_button": buildDynamicButton(p.Namespace, p.IconSupport),
	}

	if variant == "fullscreen" {
		root["third_party_server_screen@common.base_screen"].(map[string]any)["render_game_behind"] = false
		root["long_form"].(map[string]any)["size"] = []any{"100%", "100%"}
		root["button_panel"].(map[string]any)["size"] = []any{"100%", "100%"}
	}

	return root
}

func buildButtonPanelControls(p ActionFormParams) []map[string]any {
	if p.IncludeScroll {
		return []map[string]any{
			{
				"scroll@common.scrolling_panel": map[string]any{
					"size":                   []any{"100%", "100%"},
					"$show_background":       false,
					"$scrolling_content":     fmt.Sprintf("%s.button_list", p.Namespace),
					"$scroll_size":           []any{5, "100% - 4px"},
					"$scrolling_pane_size":   []any{"100% - 4px", "100% - 2px"},
					"$scrolling_pane_offset": []any{2, 0},
				},
			},
		}
	}
	return []map[string]any{
		{
			"list@" + p.Namespace + ".button_list": map[string]any{},
		},
	}
}

func buildDynamicButton(namespace string, iconSupport bool) any {
	if !iconSupport {
		return map[string]any{
			"type":        "stack_panel",
			"size":        []any{"100%", 32},
			"orientation": "horizontal",
			"controls": []map[string]any{
				{
					"form_button@common_buttons.light_text_button": map[string]any{
						"size":                              []any{"fill", 32},
						"$pressed_button_name":              "button.form_button_click",
						"$button_text":                      "#form_button_text",
						"$button_text_binding_type":         "collection",
						"$button_text_grid_collection_name": "form_buttons",
						"bindings": []map[string]any{
							{"binding_type": "collection_details", "binding_collection_name": "form_buttons"},
						},
					},
				},
			},
		}
	}

	return map[string]any{
		"type":        "stack_panel",
		"size":        []any{"100%", 32},
		"orientation": "horizontal",
		"controls": []map[string]any{
			{
				"icon_panel": map[string]any{
					"type": "panel",
					"size": []any{34, "100%c"},
					"controls": []map[string]any{
						{
							"image": map[string]any{
								"type": "image",
								"size": []any{32, 32},
								"bindings": []map[string]any{
									{"binding_name": "#form_button_texture", "binding_name_override": "#texture", "binding_type": "collection", "binding_collection_name": "form_buttons"},
									{"binding_type": "view", "source_property_name": "(not ((#texture = '') or (#texture = 'loading')))", "target_property_name": "#visible"},
								},
							},
						},
					},
				},
			},
			{
				"form_button@common_buttons.light_text_button": map[string]any{
					"size":                              []any{"fill", 32},
					"$pressed_button_name":              "button.form_button_click",
					"$button_text":                      "#form_button_text",
					"$button_text_binding_type":         "collection",
					"$button_text_grid_collection_name": "form_buttons",
					"bindings": []map[string]any{
						{"binding_type": "collection_details", "binding_collection_name": "form_buttons"},
					},
				},
			},
		},
	}
}

// marshalGeneratedForm emits a document of templates. The caller must wire
// server_form.third_party_server_screen to this namespace's main_content.
func marshalGeneratedForm(root any) string {
	if doc, ok := root.(map[string]any); ok {
		if content, ok := doc["main_content"].(map[string]any); ok && content["type"] == "factory" {
			ids := content["control_ids"].(map[string]any)
			if _, ok := ids["long_form"]; !ok {
				ids["long_form"] = "@server_form.long_form"
			}
			if _, ok := ids["custom_form"]; !ok {
				ids["custom_form"] = "@server_form.custom_form"
			}
			doc["main_content"] = map[string]any{"type": "panel", "factory": map[string]any{
				"name": "server_form_factory", "control_ids": ids,
			}}
		}
	}
	b, _ := json.MarshalIndent(root, "", "    ")
	return string(b)
}
