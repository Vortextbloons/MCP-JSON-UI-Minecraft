package bedrockui

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// BindingInfo describes a single binding found in the UI.
type BindingInfo struct {
	Path                  string `json:"path"`
	BindingType           string `json:"binding_type,omitempty"`
	BindingName           string `json:"binding_name,omitempty"`
	BindingNameOverride   string `json:"binding_name_override,omitempty"`
	BindingCollectionName string `json:"binding_collection_name,omitempty"`
	SourcePropertyName    string `json:"source_property_name,omitempty"`
	TargetPropertyName    string `json:"target_property_name,omitempty"`
	BindingCondition      string `json:"binding_condition,omitempty"`
}

// AnalyzeBindingsResult holds the output of binding analysis.
type AnalyzeBindingsResult struct {
	TotalBindings int            `json:"total_bindings"`
	ByType        map[string]int `json:"by_type"`
	ByCollection  map[string]int `json:"by_collection"`
	Bindings      []BindingInfo  `json:"bindings"`
	Warnings      []string       `json:"warnings,omitempty"`
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
							info.BindingType = "global"
							result.ByType["global"]++
						}
						if bn, ok := bMap["binding_name"].(string); ok {
							info.BindingName = bn
						} else if info.BindingType != "collection_details" && info.BindingType != "view" && info.BindingType != "none" {
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
	Name          string            `json:"name"`
	Type          string            `json:"type,omitempty"`
	InheritedFrom string            `json:"inherited_from,omitempty"`
	Properties    map[string]any    `json:"properties"`
	Children      []InheritanceNode `json:"children,omitempty"`
}

// ExpandInheritanceResult holds the output of inheritance expansion.
type ExpandInheritanceResult struct {
	Namespace string            `json:"namespace"`
	Controls  []InheritanceNode `json:"controls"`
	Warnings  []string          `json:"warnings,omitempty"`
}

// localControl preserves both the instance name and its inheritance declaration.
type localControl struct {
	name       string
	properties map[string]any
}

// ExpandInheritance resolves static local inheritance only. External files,
// variable evaluation and engine modifications require the game runtime.
func ExpandInheritance(raw string) (*ExpandInheritanceResult, error) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}
	ns, _ := doc["namespace"].(string)
	result := &ExpandInheritanceResult{Namespace: ns, Controls: []InheritanceNode{}, Warnings: []string{}}
	registry := map[string]localControl{}
	for name, value := range doc {
		if name == "namespace" {
			continue
		}
		if obj, ok := value.(map[string]any); ok {
			base := strings.SplitN(name, "@", 2)[0]
			entry := localControl{name, obj}
			registry[base] = entry
			registry[ns+"."+base] = entry
		}
	}
	names := make([]string, 0, len(doc))
	for name := range doc {
		if name != "namespace" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		if obj, ok := doc[name].(map[string]any); ok {
			result.Controls = append(result.Controls, resolveNode(name, obj, registry, "/"+name, &result.Warnings, map[string]bool{}))
		}
	}
	return result, nil
}

func resolveNode(name string, obj map[string]any, registry map[string]localControl, path string, warnings *[]string, active map[string]bool) InheritanceNode {
	parts := strings.SplitN(name, "@", 2)
	node := InheritanceNode{Name: parts[0], Properties: map[string]any{}, Children: []InheritanceNode{}}
	identity := fmt.Sprintf("%p", obj)
	if active[identity] {
		*warnings = append(*warnings, fmt.Sprintf("%s: inheritance cycle at %s", path, name))
		return node
	}
	active[identity] = true
	defer delete(active, identity)
	if len(parts) == 2 {
		node.InheritedFrom = parts[1]
		if parent, ok := registry[parts[1]]; ok {
			resolved := resolveNode(parent.name, parent.properties, registry, path, warnings, active)
			node.Type = resolved.Type
			for k, v := range resolved.Properties {
				node.Properties[k] = v
			}
			node.Children = append(node.Children, resolved.Children...)
		} else {
			*warnings = append(*warnings, fmt.Sprintf("%s: could not resolve inheritance %s (external namespaces are not loaded)", path, parts[1]))
		}
	}
	for k, v := range obj {
		if k == "type" {
			node.Type, _ = v.(string)
		} else if k != "controls" {
			node.Properties[k] = v
		}
	}
	if raw, exists := obj["controls"]; exists {
		node.Children = []InheritanceNode{}
		if children, ok := raw.([]any); ok {
			for i, child := range children {
				if wrapper, ok := child.(map[string]any); ok {
					for childName, value := range wrapper {
						if childObj, ok := value.(map[string]any); ok {
							node.Children = append(node.Children, resolveNode(childName, childObj, registry, fmt.Sprintf("%s/controls[%d]/%s", path, i, childName), warnings, active))
						}
					}
				}
			}
		}
	}
	if _, ok := obj["modifications"]; ok {
		*warnings = append(*warnings, path+": modifications are preserved but not evaluated")
	}
	return node
}
