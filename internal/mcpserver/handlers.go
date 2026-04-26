package mcpserver

import (
	"fmt"
	"strings"

	mcp "github.com/metoro-io/mcp-golang"
	"mcp-server/internal/bedrockui"
)

// ValidateUIContractArgs defines the tool input payload.
type ValidateUIContractArgs struct {
	UIJSON string `json:"ui_json" jsonschema:"required,description=The fully stringified Minecraft Bedrock JSON UI document to validate."`
}

func validateUIContract(args ValidateUIContractArgs) (*mcp.ToolResponse, error) {
	errors, err := bedrockui.ValidateDocument(args.UIJSON)
	if err != nil {
		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Internal Validation Error: %v", err))), nil
	}

	if len(errors) == 0 {
		return mcp.NewToolResponse(mcp.NewTextContent("✅ Success: UI code passes Bedrock schema and logic constraints.")), nil
	}

	var b strings.Builder
	b.WriteString("❌ Validation Failed with the following errors:\n")
	for i, msg := range errors {
		fmt.Fprintf(&b, "%d. %s\n", i+1, msg)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(b.String())), nil
}

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
