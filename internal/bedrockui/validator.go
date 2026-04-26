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
