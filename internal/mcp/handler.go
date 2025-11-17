package mcp

import (
	"encoding/json"
	"log"
	"net/http"
	"mcp-try/internal/tools"
)

type RPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[ERROR] decode: %v", err)
			writeJSON(w, map[string]any{
				"jsonrpc": "2.0",
				"id":      nil,
				"error": map[string]any{
					"code":    -32700,
					"message": "Parse error",
				},
			})
			return
		}
		log.Printf("[IN] method=%s id=%v", req.Method, req.ID)

		switch req.Method {

		case "tools/list":
			handleToolsList(w, req)

		case "tools/call":
			handleToolsCall(w, req)

		default:
			writeJSON(w, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"error": map[string]any{
					"code":    -32601,
					"message": "Method not found",
				},
			})
		}

	}
}
func handleToolsList(w http.ResponseWriter, req RPCRequest) {
    listed := tools.List()

    out := []map[string]any{}
    for _, t := range listed {
        out = append(out, map[string]any{
            "name":        t.Name,
            "description": t.Description,
        })
    }

    resp := map[string]any{
        "jsonrpc": "2.0",
        "id":      req.ID,
        "result": map[string]any{
            "tools": out,
        },
    }

    writeJSON(w, resp)
}

func handleToolsCall(w http.ResponseWriter, req RPCRequest) {
    toolName := req.Params["tool"].(string)

    output, err := tools.Call(toolName, req.Params)
    if err != nil {
        writeJSON(w, map[string]any{
            "jsonrpc": "2.0",
            "id":      req.ID,
            "error": map[string]any{
                "code":    -32000,
                "message": err.Error(),
            },
        })
        return
    }

    resp := map[string]any{
        "jsonrpc": "2.0",
        "id":      req.ID,
        "result": map[string]any{
            "output": output,
        },
    }
    writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[ERROR] encode: %v", err)
	}
}
