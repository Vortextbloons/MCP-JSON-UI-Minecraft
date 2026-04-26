package bedrockui

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BindingInfo describes a single binding found in the UI.
type BindingInfo struct {
	Path                   string `json:"path"`
	BindingType            string `json:"binding_type,omitempty"`
	BindingName            string `json:"binding_name,omitempty"`
	BindingNameOverride    string `json:"binding_name_override,omitempty"`
	BindingCollectionName  string `json:"binding_collection_name,omitempty"`
	SourcePropertyName     string `json:"source_property_name,omitempty"`
	TargetPropertyName     string `json:"target_property_name,omitempty"`
	BindingCondition       string `json:"binding_condition,omitempty"`
}

// AnalyzeBindingsResult holds the output of binding analysis.
type AnalyzeBindingsResult struct {
	TotalBindings    int           `json:"total_bindings"`
	ByType           map[string]int `json:"by_type"`
	ByCollection     map[string]int `json:"by_collection"`
	Bindings         []BindingInfo `json:"bindings"`
	Warnings         []string      `json:"warnings,omitempty"`
}

// AnalyzeBindings parses raw JSON UI and extracts all binding definitions.
func AnalyzeBindings(raw string) (*AnalyzeBindingsResult, error) {
	var doc any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	result := &AnalyzeBindingsResult{
		ByType:       make(map[string]int),
		ByCollection: make(map[string]int),
		Bindings:     []BindingInfo{},
		Warnings:     []string{},
	}

	walkBindings("", doc, result)
	result.TotalBindings = len(result.Bindings)
	return result, nil
}

func walkBindings(path string, node any, result *AnalyzeBindingsResult) {
	switch v := node.(type) {
	case map[string]any:
		// Check if this object has a bindings array
		if bindingsRaw, ok := v["bindings"]; ok {
			if bindingsArr, ok := bindingsRaw.([]any); ok {
				for i, b := range bindingsArr {
					if bMap, ok := b.(map[string]any); ok {
						info := BindingInfo{Path: fmt.Sprintf("%s/bindings[%d]", path, i)}
						if bt, ok := bMap["binding_type"].(string); ok {
							info.BindingType = bt
							result.ByType[bt]++
						} else {
							result.Warnings = append(result.Warnings, fmt.Sprintf("%s: binding missing binding_type", info.Path))
						}
						if bn, ok := bMap["binding_name"].(string); ok {
							info.BindingName = bn
						} else if info.BindingType != "collection_details" {
							result.Warnings = append(result.Warnings, fmt.Sprintf("%s: binding missing binding_name", info.Path))
						}
						if bno, ok := bMap["binding_name_override"].(string); ok {
							info.BindingNameOverride = bno
						}
						if bcn, ok := bMap["binding_collection_name"].(string); ok {
							info.BindingCollectionName = bcn
							result.ByCollection[bcn]++
						}
						if spn, ok := bMap["source_property_name"].(string); ok {
							info.SourcePropertyName = spn
						}
						if tpn, ok := bMap["target_property_name"].(string); ok {
							info.TargetPropertyName = tpn
						}
						if bc, ok := bMap["binding_condition"].(string); ok {
							info.BindingCondition = bc
						}
						result.Bindings = append(result.Bindings, info)
					}
				}
			}
		}

		// Recurse into controls
		if controlsRaw, ok := v["controls"]; ok {
			if controlsArr, ok := controlsRaw.([]any); ok {
				for i, child := range controlsArr {
					if childMap, ok := child.(map[string]any); ok {
						for childKey, childVal := range childMap {
							childPath := fmt.Sprintf("%s/controls[%d]/%s", path, i, childKey)
							walkBindings(childPath, childVal, result)
						}
					}
				}
			}
		}

		// Recurse into factory control_ids
		if factoryRaw, ok := v["factory"]; ok {
			if factoryMap, ok := factoryRaw.(map[string]any); ok {
				if controlIDsRaw, ok := factoryMap["control_ids"]; ok {
					if controlIDsMap, ok := controlIDsRaw.(map[string]any); ok {
						for cidKey, cidVal := range controlIDsMap {
							// control_ids values are strings like @namespace.control
							_ = cidVal
							_ = cidKey
						}
					}
				}
			}
		}

		// Recurse into other map values
		for k, val := range v {
			if k == "bindings" || k == "controls" {
				continue
			}
			walkBindings(fmt.Sprintf("%s/%s", path, k), val, result)
		}
	case []any:
		for i, item := range v {
			walkBindings(fmt.Sprintf("%s[%d]", path, i), item, result)
		}
	}
}

// InheritanceNode represents an expanded control with resolved inheritance.
type InheritanceNode struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type,omitempty"`
	InheritedFrom string              `json:"inherited_from,omitempty"`
	Properties map[string]any         `json:"properties"`
	Children   []InheritanceNode      `json:"children,omitempty"`
}

// ExpandInheritanceResult holds the output of inheritance expansion.
type ExpandInheritanceResult struct {
	Namespace string            `json:"namespace"`
	Controls  []InheritanceNode `json:"controls"`
	Warnings  []string          `json:"warnings,omitempty"`
}

// ExpandInheritance resolves @namespace.control syntax and returns flattened definitions.
func ExpandInheritance(raw string) (*ExpandInheritanceResult, error) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	ns, _ := doc["namespace"].(string)
	result := &ExpandInheritanceResult{
		Namespace: ns,
		Controls:  []InheritanceNode{},
		Warnings:  []string{},
	}

	// Build a registry of all root-level controls keyed by name (without @)
	registry := make(map[string]map[string]any)
	for key, val := range doc {
		if key == "namespace" {
			continue
		}
		if obj, ok := val.(map[string]any); ok {
			// Strip @prefix if present on key names in registry lookup
			cleanKey := key
			if idx := strings.Index(cleanKey, "@"); idx != -1 {
				cleanKey = cleanKey[idx+1:]
			}
			registry[cleanKey] = obj
		}
	}

	for key, val := range doc {
		if key == "namespace" {
			continue
		}
		if obj, ok := val.(map[string]any); ok {
			node := resolveNode(key, obj, registry, "/"+key, &result.Warnings)
			result.Controls = append(result.Controls, node)
		}
	}

	return result, nil
}

func resolveNode(name string, obj map[string]any, registry map[string]map[string]any, path string, warnings *[]string) InheritanceNode {
	node := InheritanceNode{
		Name:       name,
		Properties: make(map[string]any),
		Children:   []InheritanceNode{},
	}

	// Detect inheritance syntax: "my_btn@common.button"
	inheritedFrom := ""
	actualName := name
	if idx := strings.Index(name, "@"); idx != -1 {
		actualName = name[:idx]
		inheritedFrom = name[idx+1:]
		node.Name = actualName
		node.InheritedFrom = inheritedFrom

		// Look up parent in registry (try full key, then just the name after last dot)
		parent := registry[inheritedFrom]
		if parent == nil {
			// Try extracting just the control name after the last dot
			if dotIdx := strings.LastIndex(inheritedFrom, "."); dotIdx != -1 {
				parent = registry[inheritedFrom[dotIdx+1:]]
			}
		}
		if parent != nil {
			resolvedParent := resolveNode(inheritedFrom, parent, registry, path, warnings)
			// Copy parent properties
			for k, v := range resolvedParent.Properties {
				node.Properties[k] = v
			}
			// Note: we intentionally do NOT copy children from parent to keep output focused on this node
		} else {
			*warnings = append(*warnings, fmt.Sprintf("%s: could not resolve inheritance %s", path, inheritedFrom))
		}
	}

	// Apply own properties (override parent)
	for k, v := range obj {
		switch k {
		case "type":
			if t, ok := v.(string); ok {
				node.Type = t
			}
		case "controls":
			// Handled separately for children
		default:
			node.Properties[k] = v
		}
	}

	if t, ok := obj["type"].(string); ok && node.Type == "" {
		node.Type = t
	}

	// Process children
	if controlsRaw, ok := obj["controls"]; ok {
		if controlsArr, ok := controlsRaw.([]any); ok {
			for i, child := range controlsArr {
				if childMap, ok := child.(map[string]any); ok {
					for childKey, childVal := range childMap {
						if childObj, ok := childVal.(map[string]any); ok {
							childPath := fmt.Sprintf("%s/controls[%d]/%s", path, i, childKey)
							childNode := resolveNode(childKey, childObj, registry, childPath, warnings)
							node.Children = append(node.Children, childNode)
						}
					}
				}
			}
		}
	}

	return node
}
