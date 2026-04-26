package app

import (
	"os"
	"os/signal"
	"syscall"

	"mcp-server/internal/mcpserver"
)

// Run starts the MCP server and blocks until an interrupt signal is received.
func Run() error {
	server, err := mcpserver.New()
	if err != nil {
		return err
	}

	if err := server.Serve(); err != nil {
		return err
	}

	// Serve() returns immediately after connecting the transport;
	// block here so the stdio readLoop goroutine stays alive.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	return nil
}
