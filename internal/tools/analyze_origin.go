package tools

import (
    "strings"
    "time"
)

func RegisterProvenanceTools() {
    Register(Tool{
        Name:        "analyze_origin",
        Description: "Heuristically assess whether content was written by an AI system or a human.",
        Handler:     analyzeOriginHandler,
    })
}

func analyzeOriginHandler(params map[string]any) (map[string]any, error) {
    input := params["input"].(map[string]any)

    text, _ := input["text"].(string)
    sourceHint, _ := input["source_hint"].(string)

    origin := "unknown"
    confidence := 0.3
    explanation := "Baseline heuristic. No strong indicators detected."

    lower := strings.ToLower(text + " " + sourceHint)

    switch {
    case strings.Contains(lower, "chatgpt") || strings.Contains(lower, "gpt-") || strings.Contains(lower, "ai generated"):
        origin = "ai_system"
        confidence = 0.9
        explanation = "Detected explicit references to AI model tokens."

    case strings.Contains(lower, "handwritten") || strings.Contains(lower, "scanned") || strings.Contains(lower, "typed manually"):
        origin = "human"
        confidence = 0.7
        explanation = "Source hint suggests human-origin physical content."
    }

    return map[string]any{
        "origin":      origin,
        "confidence":  confidence,
        "explanation": explanation,
        "timestamp":   time.Now().UTC().Format(time.RFC3339),
    }, nil
}
