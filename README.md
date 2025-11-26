<div align="center">

# 🔍 MCP-S2xD

### AI vs Human Content Provenance Server

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-green?style=flat-square&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiPjxwYXRoIGQ9Ik0xMiAyTDIgN2wxMCA1IDEwLTV6Ii8+PHBhdGggZD0iTTIgMTdsMTAgNSAxMC01Ii8+PHBhdGggZD0iTTIgMTJsMTAgNSAxMC01Ii8+PC9zdmc+)](https://modelcontextprotocol.io/)
[![JSON-RPC](https://img.shields.io/badge/JSON--RPC-2.0-222222?style=flat-square)](https://www.jsonrpc.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat-square)](CONTRIBUTING.md)

<p align="center">
  <strong>Production-ready MCP server for detecting AI-generated vs human-written content</strong>
</p>

[Features](#-features) •
[Quick Start](#-quick-start) •
[API Reference](#-api-reference) •
[Tools](#-tool-catalog) •
[Configuration](#%EF%B8%8F-configuration)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🛡️ Security First
- Rate limiting (token bucket)
- API key authentication
- Request size limits
- Security headers (CSP, XSS, CORS)
- Graceful shutdown

</td>
<td width="50%">

### 📊 Observability
- Structured JSON logging (slog)
- Request ID tracking
- Metrics endpoint
- Request duration tracking
- Tool call statistics

</td>
</tr>
<tr>
<td width="50%">

### 🔧 Production Ready
- Environment-based config
- Health & readiness probes
- Panic recovery middleware
- Request timeouts
- SQLite persistence (optional)

</td>
<td width="50%">

### 🤖 AI Detection
- Multi-signal heuristic analysis
- **Statistical detection** (entropy, n-gram, burstiness, Zipf)
- Style profiling
- Code origin detection
- Batch analysis
- Content fingerprinting

</td>
</tr>
</table>

---

## 🚀 Quick Start

```bash
# Clone and run
git clone https://github.com/AdhiDevX369/MCP-S2xD.git
cd MCP-S2xD
go mod tidy
go run ./cmd/server
```

Server starts at `http://localhost:3000`

### Verify Installation

```bash
# Health check
curl http://localhost:3000/health

# List available tools
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

---

## 📡 API Reference

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/ready` | GET | Readiness probe |
| `/metrics` | GET | Server metrics (requests, errors, tool calls) |
| `/mcp` | POST | JSON-RPC 2.0 MCP endpoint |

### MCP Methods

| Method | Description |
|--------|-------------|
| `tools/list` | List all available tools with schemas |
| `tools/call` | Execute a tool with parameters |

### Example: Analyze Content Origin

```json
// Request
{
  "jsonrpc": "2.0",
  "id": "req-42",
  "method": "tools/call",
  "params": {
    "tool": "analyze_origin",
    "input": {
      "text": "Furthermore, it is important to note that the implementation follows best practices.",
      "source_hint": "documentation"
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": "req-42",
  "result": {
    "output": {
      "origin": "ai_system",
      "confidence": 0.85,
      "explanation": "AI indicators detected: [hedging_language, uniform_sentence_length]",
      "signals": ["hedging_language", "uniform_sentence_length"],
      "word_count": 12,
      "timestamp": "2025-11-26T09:00:00Z"
    }
  }
}
```

### Example: Statistical Detection

```json
// Request - Comprehensive statistical analysis
{
  "jsonrpc": "2.0",
  "id": "stat-1",
  "method": "tools/call",
  "params": {
    "tool": "statistical_detect",
    "arguments": {
      "input": {
        "text": "It is important to note that in the evolving landscape of AI..."
      }
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": "stat-1",
  "result": {
    "output": {
      "origin": "ai_system",
      "confidence": 0.78,
      "weighted_score": 0.42,
      "signals": ["ai_phrases", "uniform_structure", "zipf_deviation"],
      "entropy_result": { "average_entropy": 3.2, "origin": "ai_system" },
      "ngram_result": { "ai_phrases_found": ["it is important to note"], "origin": "ai_system" },
      "burstiness_result": { "burstiness_score": 0.18, "origin": "ai_system" },
      "zipf_result": { "r_squared": 0.72, "origin": "ai_system" }
    }
  }
}
```

### Example: Entropy Analysis

```json
// Request
{
  "jsonrpc": "2.0",
  "id": "ent-1",
  "method": "tools/call",
  "params": {
    "tool": "entropy_analyze",
    "arguments": {
      "input": { "text": "Your text to analyze..." }
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": "ent-1",
  "result": {
    "output": {
      "origin": "human",
      "confidence": 0.72,
      "char_entropy": 4.23,
      "word_entropy": 4.56,
      "bigram_entropy": 4.89,
      "average_entropy": 4.56,
      "explanation": "High entropy indicates natural human variation."
    }
  }
}
```

---

## 🧰 Tool Catalog

### Core Detection Tools

| Tool | Description | Use Case |
|------|-------------|----------|
| `analyze_origin` | Deep AI/human detection with multiple heuristics | Content moderation, authorship verification |
| `batch_analyze` | Analyze multiple samples, get aggregate stats | Dataset analysis, bulk processing |
| `content_fingerprint` | SHA256 hash + structure hash | Provenance tracking, deduplication |
| `compare_texts` | Style similarity & origin comparison | Plagiarism detection, authorship matching |
| `code_origin` | Specialized code analysis | Code review, AI-assisted code detection |
| `style_profile` | Generate authorship metrics | Writer profiling, style analysis |

### Statistical Detection Tools

| Tool | Description | Key Metrics |
|------|-------------|-------------|
| `entropy_analyze` | Calculate text entropy patterns | Character, word, bigram entropy |
| `ngram_analyze` | Detect AI phrases & repetition patterns | AI phrase detection, n-gram frequency |
| `burstiness_analyze` | Measure text uniformity | Sentence/word length variance |
| `zipf_analyze` | Check natural language distribution | Zipf's law compliance, R² score |
| `statistical_detect` | Ensemble of all statistical methods | Weighted multi-signal analysis |

### Utility Tools

| Tool | Description |
|------|-------------|
| `dummy_echo` | Echo input for testing |

### Detection Signals

```text
Statistical Signals              Heuristic Signals
─────────────────────────────    ─────────────────────────────
• Low entropy (repetitive)       • Uniform sentence length
• AI-typical phrases             • Transition word density
• Uniform burstiness             • Hedging language
• Zipf's law deviation           • Structured lists
• N-gram repetition              • Explicit AI mentions

Human Indicators
─────────────────────────────
• High entropy (varied)
• Natural burstiness
• Zipf-compliant distribution
• Common typos
• Informal contractions
```

### Statistical Methods Deep Dive

#### Entropy Analysis

Measures information density using Shannon entropy at character, word, and bigram levels.

- **Low entropy (<3.5)** → Repetitive, predictable patterns → AI likely
- **High entropy (>4.5)** → Natural variation → Human likely

#### N-gram Analysis

Detects overused AI phrases and repetition patterns.

- Scans for 40+ known AI-typical phrases ("it is important to note", "in conclusion", etc.)
- Measures n-gram repetition rate
- Combined scoring for detection

#### Burstiness Analysis

Measures variance in sentence/word lengths using coefficient of variation.

- **Low burstiness (<0.3)** → Uniform structure → AI likely
- **High burstiness (>0.5)** → Natural variation → Human likely

#### Zipf's Law Analysis

Natural language follows Zipf's distribution (word frequency ∝ 1/rank).

- Calculates linear regression on log-log frequency plot
- **R² < 0.8** → Deviation from natural distribution → AI likely
- **R² > 0.9** → Natural Zipf compliance → Human likely

---

## ⚙️ Configuration

All settings via environment variables:

```bash
# Server
SERVER_PORT=:3000
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=15s

# Security
SECURITY_API_KEY_ENABLED=false
SECURITY_API_KEY=your-secret-key
SECURITY_RATE_LIMIT_RPS=100
SECURITY_MAX_REQUEST_SIZE=1048576

# Logging
LOG_LEVEL=info          # debug, info, warn, error
LOG_FORMAT=json         # json, text

# Database (optional)
DATABASE_PATH=data/mcp.db
```

See [`.env.example`](.env.example) for full reference.

---

## 📁 Project Structure

```
mcp-s2xd/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Environment configuration
│   ├── httpserver/      # HTTP server with middleware
│   ├── logger/          # Structured logging (slog)
│   ├── mcp/             # JSON-RPC handler
│   ├── metrics/         # Request & tool metrics
│   ├── middleware/      # Security middleware stack
│   ├── rpcerror/        # Centralized error types
│   ├── storage/         # Attestation persistence
│   └── tools/           # MCP tool implementations
├── .env.example
├── go.mod
└── README.md
```

---

## 🔒 Security Middleware Stack

```
Request → Rate Limit → API Key Auth → Max Body Size → CORS → Security Headers → Handler
```

| Middleware | Purpose |
|------------|---------|
| `Recovery` | Panic recovery, prevents crashes |
| `RequestID` | X-Request-ID header tracking |
| `Logging` | Structured request logging |
| `RateLimit` | Token bucket rate limiting |
| `APIKeyAuth` | Optional API key validation |
| `MaxBodySize` | Request size limits |
| `CORS` | Cross-origin resource sharing |
| `SecurityHeaders` | CSP, XSS, Frame options |
| `Timeout` | Request context timeouts |

---

## 📈 Metrics

```bash
curl http://localhost:3000/metrics
```

```json
{
  "uptime_seconds": 3600,
  "total_requests": 1250,
  "total_errors": 3,
  "avg_duration_ms": 2.5,
  "tool_calls": {
    "analyze_origin": 800,
    "code_origin": 300,
    "compare_texts": 150
  }
}
```

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

---

## 📚 Resources

### Protocol Specifications

| Resource | Description |
|----------|-------------|
| [Model Context Protocol](https://modelcontextprotocol.io/) | Official MCP specification |
| [MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk) | Reference implementation |
| [JSON-RPC 2.0](https://www.jsonrpc.org/specification) | Transport protocol spec |

### AI Detection Research

| Paper/Resource | Topic |
|----------------|-------|
| [GPTZero Research](https://gptzero.me/technology) | Perplexity & burstiness detection |
| [DetectGPT](https://arxiv.org/abs/2301.11305) | Zero-shot machine-generated text detection |
| [Zipf's Law in NLP](https://en.wikipedia.org/wiki/Zipf%27s_law) | Natural language frequency distribution |
| [Shannon Entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) | Information theory fundamentals |

### Go Documentation

| Resource | Description |
|----------|-------------|
| [Go Official Docs](https://go.dev/doc/) | Language reference |
| [slog Package](https://pkg.go.dev/log/slog) | Structured logging |
| [net/http](https://pkg.go.dev/net/http) | HTTP server package |

---

<div align="center">

**Built with ❤️ for trustworthy AI ecosystems**

[![GitHub](https://img.shields.io/badge/GitHub-AdhiDevX369-181717?style=flat-square&logo=github)](https://github.com/AdhiDevX369)
[![Website](https://img.shields.io/badge/Web-adithyabandara.com-1E90FF?style=flat-square&logo=safari&logoColor=white)](https://adithyabandara.com)

</div>
