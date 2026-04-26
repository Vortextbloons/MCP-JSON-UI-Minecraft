package main

import (
	"log"

	"bedrock-json-ui-mcp/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
