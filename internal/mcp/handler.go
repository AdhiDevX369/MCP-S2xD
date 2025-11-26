package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"mcp-try/internal/logger"
	"mcp-try/internal/metrics"
	"mcp-try/internal/rpcerror"
	"mcp-try/internal/storage"
	"mcp-try/internal/tools"
)

type RPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type RPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics.IncrementRequest()

	var req RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode request", "error", err)
		metrics.IncrementError()
		h.writeError(w, nil, rpcerror.Parse("invalid JSON"))
		return
	}

	logger.Info(ctx, "incoming request", "method", req.Method, "id", req.ID)

	var result any
	var rpcErr *rpcerror.RPCError

	switch req.Method {
	case "tools/list":
		result = h.handleToolsList()
	case "tools/call":
		result, rpcErr = h.handleToolsCall(ctx, req)
	default:
		rpcErr = rpcerror.NotFound(req.Method)
	}

	metrics.RecordDuration(time.Since(start))

	if rpcErr != nil {
		metrics.IncrementError()
		h.writeError(w, req.ID, rpcErr)
		return
	}

	h.writeSuccess(w, req.ID, result)
}

func (h *Handler) handleToolsList() map[string]any {
	listed := tools.List()

	out := []map[string]any{}
	for _, t := range listed {
		toolInfo := map[string]any{
			"name":        t.Name,
			"description": t.Description,
		}
		if t.InputSchema != nil {
			toolInfo["inputSchema"] = t.InputSchema
		}
		out = append(out, toolInfo)
	}

	return map[string]any{"tools": out}
}

func (h *Handler) handleToolsCall(ctx context.Context, req RPCRequest) (map[string]any, *rpcerror.RPCError) {
	if req.Params == nil {
		return nil, rpcerror.InvalidReq("missing params")
	}

	toolName, ok := req.Params["tool"].(string)
	if !ok || toolName == "" {
		return nil, rpcerror.InvalidParam("'tool' must be a non-empty string")
	}

	args, _ := req.Params["arguments"].(map[string]any)
	if args == nil {
		args = make(map[string]any)
	}

	metrics.IncrementToolCall(toolName)

	output, err := tools.Call(toolName, args)
	if err != nil {
		return nil, rpcerror.ToolExec(toolName, err.Error())
	}

	if h.store != nil && toolName == "analyze_origin" {
		h.saveAttestation(ctx, toolName, args, output)
	}

	return map[string]any{"output": output}, nil
}

func (h *Handler) saveAttestation(ctx context.Context, toolName string, params, output map[string]any) {
	origin, _ := output["origin"].(string)
	confidence, _ := output["confidence"].(float64)
	explanation, _ := output["explanation"].(string)

	a := &storage.Attestation{
		ID:          storage.GenerateID(),
		ToolName:    toolName,
		Input:       params,
		Output:      output,
		Origin:      origin,
		Confidence:  confidence,
		Explanation: explanation,
		RequestID:   logger.GetRequestID(ctx),
		CreatedAt:   time.Now().UTC(),
	}

	if err := h.store.SaveAttestation(ctx, a); err != nil {
		logger.Error(ctx, "failed to save attestation", "error", err)
	} else {
		logger.Debug(ctx, "attestation saved", "id", a.ID)
	}
}

func (h *Handler) writeSuccess(w http.ResponseWriter, id any, result any) {
	h.writeJSON(w, RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

func (h *Handler) writeError(w http.ResponseWriter, id any, err *rpcerror.RPCError) {
	h.writeJSON(w, RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err.ToMap(),
	})
}

func (h *Handler) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
