package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"mcp-try/internal/mcp"
)

func Start(addr string) error{
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "S2xD MCP server OK")
	})

	mux.Handle("/mcp", mcp.Handler())
	
	srv := &http.Server{
		Addr:	addr,
		Handler: mux,
	}

	log.Printf("HTTP Server listining on %s\n", addr)
	return srv.ListenAndServe()
}