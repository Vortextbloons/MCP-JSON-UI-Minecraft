package mcpserver

import (
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"mcp-server/internal/bedrockui"
)

const (
	serverName        = "minecraft-bedrock-json-ui"
	serverVersion     = "0.1.0"
	serverBlurb       = "Read bedrock://docs/ui_reference (quick) or bedrock://docs/json_ui_full_reference (full) before authoring JSON UI. Validate each draft with validate_ui_contract and iterate until it passes."
	toolName          = "validate_ui_contract"
	toolDescription   = "Validates Minecraft Bedrock JSON UI strings against the scaffolded Bedrock contract."
	resourceName      = "ui_reference"
	resourceBlurb     = "Curated Bedrock JSON UI reference derived from the design and technical docs."
	fullResourceName  = "json_ui_full_reference"
	fullResourceBlurb = "Restructured, comprehensive Bedrock JSON UI documentation covering all element types, properties, hardcoded values, and more."
)

// New constructs the MCP server and registers the tool/resource pair.
func New() (*mcp.Server, error) {
	server := mcp.NewServer(
		stdio.NewStdioServerTransport(),
		mcp.WithName(serverName),
		mcp.WithVersion(serverVersion),
		mcp.WithInstructions(serverBlurb),
	)

	if err := server.RegisterResource(bedrockui.ReferenceURI, resourceName, resourceBlurb, bedrockui.ReferenceMIMEType, uiReferenceResource); err != nil {
		return nil, err
	}

	if err := server.RegisterResource(bedrockui.FullReferenceURI, fullResourceName, fullResourceBlurb, bedrockui.ReferenceMIMEType, fullReferenceResource); err != nil {
		return nil, err
	}

	if err := server.RegisterTool(toolName, toolDescription, validateUIContract); err != nil {
		return nil, err
	}

	return server, nil
}
