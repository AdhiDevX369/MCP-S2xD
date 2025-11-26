package metrics

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	mu              sync.RWMutex
	requestCount    int64
	errorCount      int64
	toolCalls       map[string]int64
	requestDuration []time.Duration
	startTime       time.Time
}

var global = &Metrics{
	toolCalls: make(map[string]int64),
	startTime: time.Now(),
}

func IncrementRequest() {
	atomic.AddInt64(&global.requestCount, 1)
}

func IncrementError() {
	atomic.AddInt64(&global.errorCount, 1)
}

func IncrementToolCall(toolName string) {
	global.mu.Lock()
	global.toolCalls[toolName]++
	global.mu.Unlock()
}

func RecordDuration(d time.Duration) {
	global.mu.Lock()
	global.requestDuration = append(global.requestDuration, d)
	if len(global.requestDuration) > 1000 {
		global.requestDuration = global.requestDuration[500:]
	}
	global.mu.Unlock()
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		global.mu.RLock()
		defer global.mu.RUnlock()

		var avgDuration float64
		if len(global.requestDuration) > 0 {
			var total time.Duration
			for _, d := range global.requestDuration {
				total += d
			}
			avgDuration = float64(total.Milliseconds()) / float64(len(global.requestDuration))
		}

		toolCalls := make(map[string]int64)
		for k, v := range global.toolCalls {
			toolCalls[k] = v
		}

		stats := map[string]any{
			"uptime_seconds":  time.Since(global.startTime).Seconds(),
			"total_requests":  atomic.LoadInt64(&global.requestCount),
			"total_errors":    atomic.LoadInt64(&global.errorCount),
			"avg_duration_ms": avgDuration,
			"tool_calls":      toolCalls,
			"sample_count":    len(global.requestDuration),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
