package app

import "mcp-server/internal/mcpserver"

// Run starts the MCP server.
func Run() error {
	server, err := mcpserver.New()
	if err != nil {
		return err
	}

	return server.Serve()
}
