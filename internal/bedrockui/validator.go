package bedrockui

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	compiledSchema *jsonschema.Schema
	schemaErr      error
	schemaOnce     sync.Once
)

// Schema returns the compiled Bedrock UI schema.
func Schema() (*jsonschema.Schema, error) {
	schemaOnce.Do(func() {
		compiler := jsonschema.NewCompiler()
		compiler.Draft = jsonschema.Draft2020
		compiler.AssertFormat = true

		if err := compiler.AddResource(SchemaURI, strings.NewReader(SchemaSource())); err != nil {
			schemaErr = fmt.Errorf("load schema: %w", err)
			return
		}

		compiledSchema, schemaErr = compiler.Compile(SchemaURI)
		if schemaErr != nil {
			schemaErr = fmt.Errorf("compile schema: %w", schemaErr)
		}
	})

	return compiledSchema, schemaErr
}

// ValidateDocument validates raw JSON against the scaffolded Bedrock UI contract.
func ValidateDocument(raw string) ([]string, error) {
	schema, err := Schema()
	if err != nil {
		return nil, err
	}

	var document any
	if err := json.Unmarshal([]byte(raw), &document); err != nil {
		return []string{fmt.Sprintf("Syntax Error: Invalid JSON structure. %v", err)}, nil
	}

	if err := schema.Validate(document); err != nil {
		if validationErr, ok := err.(*jsonschema.ValidationError); ok {
			return formatValidationErrors(validationErr), nil
		}

		return []string{err.Error()}, nil
	}

	if errs := validateNamespace(document); len(errs) > 0 {
		return errs, nil
	}

	hints := validateHints(document)
	if len(hints) > 0 {
		return hints, nil
	}

	return nil, nil
}

func validateNamespace(document any) []string {
	root, ok := document.(map[string]any)
	if !ok {
		return []string{"Logic Error: Root document must be a JSON object."}
	}

	namespace, ok := root["namespace"].(string)
	if !ok || strings.TrimSpace(namespace) == "" {
		return []string{"Logic Error: Missing required 'namespace' declaration at root level."}
	}

	return nil
}

func validateHints(document any) []string {
	var hints []string
	root, ok := document.(map[string]any)
	if !ok {
		return hints
	}

	for key, val := range root {
		if key == "namespace" {
			continue
		}
		obj, ok := val.(map[string]any)
		if !ok {
			continue
		}

		// Check for root-level animation objects
		if _, hasAnimType := obj["anim_type"]; hasAnimType {
			hints = append(hints, fmt.Sprintf("Hint at /%s: Root-level animation objects are not loaded by the engine. Place animations in $definitions or reference them via the 'anims' array.", key))
			continue
		}

		// Check for container types missing controls
		containerTypes := map[string]bool{"panel": true, "stack_panel": true, "collection_panel": true, "grid": true}
		if t, ok := obj["type"].(string); ok && containerTypes[t] {
			if _, hasControls := obj["controls"]; !hasControls {
				hints = append(hints, fmt.Sprintf("Hint at /%s: Container type '%s' should include a 'controls' array for child elements.", key, t))
			}
		}

		// Check for factory props on non-factory elements
		if _, hasFactory := obj["factory"]; hasFactory {
			if t, ok := obj["type"].(string); ok && t != "factory" {
				hints = append(hints, fmt.Sprintf("Hint at /%s: 'factory' property is only valid on elements with type 'factory', not '%s'.", key, t))
			}
		}

		// Check for collection_name on non-collection elements
		if _, hasCollection := obj["collection_name"]; hasCollection {
			if t, ok := obj["type"].(string); ok && t != "grid" && t != "stack_panel" && t != "collection_panel" {
				hints = append(hints, fmt.Sprintf("Hint at /%s: 'collection_name' is typically only used by grid, stack_panel, or collection_panel, not '%s'.", key, t))
			}
		}

		// Recurse into controls
		if controls, ok := obj["controls"].([]any); ok {
			for _, child := range controls {
				childObj, ok := child.(map[string]any)
				if !ok {
					continue
				}
				for childKey, childVal := range childObj {
					childMap, ok := childVal.(map[string]any)
					if !ok {
						continue
					}
					// Check bindings in children
					if bindings, ok := childMap["bindings"].([]any); ok {
						for _, b := range bindings {
							bMap, ok := b.(map[string]any)
							if !ok {
								continue
							}
							if _, hasName := bMap["binding_name"]; !hasName {
								hints = append(hints, fmt.Sprintf("Hint at /%s/%s/bindings: binding is missing required 'binding_name'.", key, childKey))
							}
						}
					}
				}
			}
		}
	}

	return hints
}

func formatValidationErrors(err *jsonschema.ValidationError) []string {
	output := err.BasicOutput()
	messages := make([]string, 0, len(output.Errors))

	for _, item := range output.Errors {
		location := item.InstanceLocation
		if location == "" {
			location = "/"
		}

		kind := "Logic Error"
		if location == "/" {
			kind = "Schema Error"
		}

		messages = append(messages, fmt.Sprintf("%s at %s: %s", kind, location, item.Error))
	}

	return messages
}
