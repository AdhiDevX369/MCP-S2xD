package main

import (
	"log"

	"mcp-try/internal/httpserver"
	"mcp-try/internal/tools"
)

func main() {
	tools.RegisterDefaults()
	tools.RegisterProvenanceTools() 
	if err := httpserver.Start(":3000"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}