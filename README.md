# MCP Try: AI vs Human Work Signal Server

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge) ![Status](https://img.shields.io/badge/Status-Prototype-FF6F61?style=for-the-badge) ![License](https://img.shields.io/badge/License-MIT-4B8BBE?style=for-the-badge) ![MCP Ready](https://img.shields.io/badge/MCP-Server-green?style=for-the-badge) [![Site](https://img.shields.io/badge/adithyabandara.com-contact-1E90FF?style=for-the-badge)](https://adithyabandara.com)

## ✦ Vision

Crafting a lightweight Model Context Protocol (MCP) server that helps tooling pipelines identify and justify whether work was produced by AI systems or humans. This repo is the starting foundation for richer provenance signals.

## ✦ Features

- Structured JSON-RPC endpoint at `/mcp` with `tools/list` and `tools/call`
- Extensible tool registry for future AI/Human justification heuristics
- Minimal logging to trace ingestion paths for downstream attribution models

## ✦ Quick Start

```
go mod tidy
go run .
```

Server boots at `http://localhost:3000`.

## ✦ MCP Contract

| Method       | Description                               | Result Skeleton                              |
|--------------|-------------------------------------------|----------------------------------------------|
| `tools/list` | Lists available provenance tools          | `{ tools: [{ name: "echo", ... }] }`         |
| `tools/call` | Echoes submitted text for testing hookups | `{ result: { output: { text: "<input>" } } }` |

## ✦ Roadmap Highlights

1. Add scoring tool for AI-likelihood vs human craftsmanship
2. Persist justification metadata with signed attestations
3. Expose provenance summaries via rich UI dashboards

## ✦ Contributing

1. Fork & branch: `git checkout -b feature/<name>`
2. Keep functions small and observable
3. Open PR with justification notes for AI/Human heuristics

## ✦ Learning Resources

- [Go.dev Learn](https://go.dev/learn/)
- [Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Model Context Protocol Spec](https://modelcontextprotocol.io/)

---
Built with clarity-first Go patterns for trustworthy MCP ecosystems.
