package main

import (
	"encoding/json"	
	"fmt"
	"net/http"
)

type RPCRequest struct {
	JSONRPC string	  `json:"jsonrpc"`
	ID 	any		  `json:"id"`
	Method  string	  `json:"method"`
	Params  map[string]any `json:"params"`
}

func main() {
	fmt.Println("Server Running on http://localhost:3000")

	// Defualt route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MCP Base Server is OK")
	})

	// MCP route
	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only allowed POST allowed", http.StatusMethodNotAllowed)
			return
		}
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		fmt.Println("Incoming MCP request:", req)

		switch req.Method {
			case "tools/list":
				resp := map[string]any{
					"jsonrpc": "2.0",
					"id": req.ID,
					"result": map[string]any{
						"tools": []map[string]any{
							{
                            "name":        "echo",
                            "description": "Echo input text back",
                        },
						},
					},
				}
				json.NewEncoder(w).Encode(resp)

				case "tools/call":
					input := req.Params["input"].(map[string]any)
					text := input["text"].(string)

					resp := map[string]any{
						"jsonrpc": "2.0",
						"id": req.ID,
						"result": map[string]any{
							"output": map[string]any{
								"text": text,
							},
						},
					}
					json.NewEncoder(w).Encode(resp)
			default:
				json.NewEncoder(w).Encode(map[string]any{
					"jsonrpc": "2.0",
					"id": req.ID,
					"error": map[string]any{
						"code": -32601,
						"message": "Method not found",
					},
				})
		}
	})
	http.ListenAndServe(":3000", nil)
}