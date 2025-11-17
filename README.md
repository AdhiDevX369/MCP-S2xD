# MCP Try: AI vs Human Work Signal Server

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge) ![Status](https://img.shields.io/badge/Status-Prototype-FF6F61?style=for-the-badge) ![License](https://img.shields.io/badge/License-MIT-4B8BBE?style=for-the-badge) ![Transport](https://img.shields.io/badge/JSON--RPC-2.0-222222?style=for-the-badge) ![MCP Ready](https://img.shields.io/badge/MCP-Server-green?style=for-the-badge) [![Site](https://img.shields.io/badge/adithyabandara.com-contact-1E90FF?style=for-the-badge)](https://adithyabandara.com)

## ✦ Vision

Crafting a lightweight Model Context Protocol (MCP) server that helps tooling pipelines identify and justify whether work was produced by AI systems or humans. This repo is the starting foundation for richer provenance signals.

## ✦ Stack Highlights

- Go HTTP server hosted from `cmd/server` delegating to internal packages
- JSON-RPC 2.0 layer under `internal/mcp` to stay MCP-compatible
- Pluggable tool registry (`internal/tools`) for provenance heuristics
- HTTP mux at `/` for health checks and `/mcp` for contract traffic

## ✦ Quick Start

```bash
go mod tidy
go run .
```

Server boots at `http://localhost:3000`.

### Run While Developing

```bash
go run ./cmd/server
```

The server logs inbound method calls and will exit on fatal errors to keep the feedback loop tight.

## ✦ MCP Contract

| Method       | Description                               | Result Skeleton                              |
|--------------|-------------------------------------------|----------------------------------------------|
| `tools/list` | Lists available provenance tools          | `{ tools: [{ name: "echo", ... }] }`         |
| `tools/call` | Echoes submitted text for testing hookups | `{ result: { output: { text: "<input>" } } }` |

### Example JSON-RPC Exchange

```json
{
  "jsonrpc": "2.0",
  "id": "req-42",
  "method": "tools/call",
  "params": {
    "tool": "analyze_origin",
    "input": {
      "text": "Draft prepared by ChatGPT 4o",
      "source_hint": "slack export"
    }
  }
}
```

Response (truncated):

```json
{
  "result": {
    "output": {
      "origin": "ai_system",
      "confidence": 0.9,
      "explanation": "Detected explicit references to AI model tokens.",
      "timestamp": "..."
    }
  }
}
```

## ✦ Directory Map

| Path | Purpose |
|------|---------|
| `cmd/server` | Entry point that registers tools and starts the HTTP listener |
| `internal/httpserver` | HTTP mux + server bootstrap, wires `/` and `/mcp` routes |
| `internal/mcp` | JSON-RPC parsing, method dispatch, error responses |
| `internal/tools` | Registry, default echo tool, and provenance heuristics |

## ✦ Tool Catalog

| Tool | File | Description |
|------|------|-------------|
| `dummy_echo` | `internal/tools/dummy.go` | Loops input payloads back to clients for pipeline smoke tests |
| `analyze_origin` | `internal/tools/analyze_origin.go` | Simple string-heuristic origin detector with confidence + explanation |

Add your own tools via `tools.Register(...)` and they will flow automatically through `tools/list` and `tools/call`.

## ✦ Roadmap Highlights

1. Expand heuristics beyond keyword detection (style, cadence, metadata)
2. Persist justification metadata with signed attestations
3. Expose provenance summaries via dashboards or downstream APIs

## ✦ Contributing

1. Fork & branch: `git checkout -b feature/<name>`
2. Keep functions small, observable, and well-typed
3. Add unit coverage for new tools or handlers when logic grows
4. Open PR with justification notes for AI/Human heuristics

## ✦ Learning Resources

- [Go.dev Learn](https://go.dev/learn/)
- [Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Model Context Protocol Spec](https://modelcontextprotocol.io/)

---
Built with clarity-first Go patterns for trustworthy MCP ecosystems.
