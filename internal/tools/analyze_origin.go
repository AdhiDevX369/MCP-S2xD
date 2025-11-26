package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func RegisterProvenanceTools() {
	Register(Tool{
		Name:        "analyze_origin",
		Description: "Deep analysis to determine if content was written by AI or human using multiple heuristics.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "The content to analyze for origin detection",
						},
						"source_hint": map[string]any{
							"type":        "string",
							"description": "Optional hint about the content source (e.g., 'email', 'code', 'document')",
						},
						"context": map[string]any{
							"type":        "string",
							"description": "Additional context about where the content came from",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: analyzeOriginHandler,
	})

	Register(Tool{
		Name:        "batch_analyze",
		Description: "Analyze multiple text samples and return aggregate origin statistics.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"samples": map[string]any{
							"type":        "array",
							"description": "Array of text samples to analyze",
							"items":       map[string]any{"type": "string"},
						},
					},
					"required": []string{"samples"},
				},
			},
			"required": []string{"input"},
		},
		Handler: batchAnalyzeHandler,
	})

	Register(Tool{
		Name:        "content_fingerprint",
		Description: "Generate a unique fingerprint for content to track provenance across systems.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Content to fingerprint",
						},
						"metadata": map[string]any{
							"type":        "object",
							"description": "Optional metadata to include in fingerprint",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: contentFingerprintHandler,
	})

	Register(Tool{
		Name:        "compare_texts",
		Description: "Compare two texts to determine if they share the same origin or authorship style.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text_a": map[string]any{
							"type":        "string",
							"description": "First text sample",
						},
						"text_b": map[string]any{
							"type":        "string",
							"description": "Second text sample",
						},
					},
					"required": []string{"text_a", "text_b"},
				},
			},
			"required": []string{"input"},
		},
		Handler: compareTextsHandler,
	})

	Register(Tool{
		Name:        "code_origin",
		Description: "Specialized analysis for determining if source code was AI-generated or human-written.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"code": map[string]any{
							"type":        "string",
							"description": "Source code to analyze",
						},
						"language": map[string]any{
							"type":        "string",
							"description": "Programming language (e.g., 'python', 'javascript', 'go')",
						},
					},
					"required": []string{"code"},
				},
			},
			"required": []string{"input"},
		},
		Handler: codeOriginHandler,
	})

	Register(Tool{
		Name:        "style_profile",
		Description: "Generate a writing style profile that can be used to identify authorship patterns.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to profile",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: styleProfileHandler,
	})
}

type analysisResult struct {
	origin      string
	confidence  float64
	explanation string
	signals     []string
}

func analyzeOriginHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	sourceHint, _ := input["source_hint"].(string)

	result := deepAnalyze(text, sourceHint)

	return map[string]any{
		"origin":      result.origin,
		"confidence":  result.confidence,
		"explanation": result.explanation,
		"signals":     result.signals,
		"text_length": len(text),
		"word_count":  len(strings.Fields(text)),
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func deepAnalyze(text, sourceHint string) analysisResult {
	var signals []string
	var aiScore, humanScore float64

	lower := strings.ToLower(text + " " + sourceHint)

	if containsAny(lower, []string{"chatgpt", "gpt-4", "gpt-3", "claude", "gemini", "copilot", "ai generated", "generated by ai"}) {
		signals = append(signals, "explicit_ai_mention")
		aiScore += 0.4
	}

	if containsAny(lower, []string{"handwritten", "scanned", "typed manually", "my personal", "i think", "in my opinion"}) {
		signals = append(signals, "personal_language")
		humanScore += 0.25
	}

	sentences := splitSentences(text)
	if len(sentences) > 3 {
		lengths := make([]float64, len(sentences))
		for i, s := range sentences {
			lengths[i] = float64(len(strings.Fields(s)))
		}
		variance := calculateVariance(lengths)

		if variance < 15 {
			signals = append(signals, "uniform_sentence_length")
			aiScore += 0.15
		} else if variance > 50 {
			signals = append(signals, "varied_sentence_length")
			humanScore += 0.1
		}
	}

	transitionWords := []string{"furthermore", "moreover", "additionally", "consequently", "nevertheless", "in conclusion", "to summarize"}
	transitionCount := 0
	for _, tw := range transitionWords {
		transitionCount += strings.Count(lower, tw)
	}
	if len(sentences) > 0 && float64(transitionCount)/float64(len(sentences)) > 0.3 {
		signals = append(signals, "high_transition_density")
		aiScore += 0.1
	}

	hedging := []string{"it's important to note", "it should be noted", "generally speaking", "in most cases", "typically"}
	for _, h := range hedging {
		if strings.Contains(lower, h) {
			signals = append(signals, "hedging_language")
			aiScore += 0.1
			break
		}
	}

	typoPatterns := regexp.MustCompile(`\b(teh|recieve|occured|seperate|definately)\b`)
	if typoPatterns.MatchString(lower) {
		signals = append(signals, "common_typos")
		humanScore += 0.15
	}

	contractions := []string{"i'm", "don't", "won't", "can't", "shouldn't", "wouldn't", "i've", "we're", "they're"}
	contractionCount := 0
	for _, c := range contractions {
		contractionCount += strings.Count(lower, c)
	}
	words := strings.Fields(text)
	if len(words) > 0 && float64(contractionCount)/float64(len(words)) > 0.02 {
		signals = append(signals, "informal_contractions")
		humanScore += 0.1
	}

	if strings.Contains(text, "!") || strings.Contains(text, "...") || strings.Contains(text, "??") {
		signals = append(signals, "expressive_punctuation")
		humanScore += 0.05
	}

	listPatterns := regexp.MustCompile(`(?m)^[\s]*[-*•]\s|^\s*\d+\.\s`)
	if listPatterns.MatchString(text) {
		matches := listPatterns.FindAllString(text, -1)
		if len(matches) > 3 {
			signals = append(signals, "structured_lists")
			aiScore += 0.1
		}
	}

	origin := "unknown"
	confidence := 0.3
	explanation := "No strong indicators detected."

	totalScore := aiScore + humanScore
	if totalScore > 0 {
		if aiScore > humanScore {
			origin = "ai_system"
			confidence = math.Min(0.95, 0.5+aiScore)
			explanation = fmt.Sprintf("AI indicators detected: %v", signals)
		} else if humanScore > aiScore {
			origin = "human"
			confidence = math.Min(0.95, 0.5+humanScore)
			explanation = fmt.Sprintf("Human indicators detected: %v", signals)
		} else {
			origin = "mixed"
			confidence = 0.5
			explanation = "Mixed signals suggest possible AI-assisted human writing or human-edited AI content."
		}
	}

	return analysisResult{
		origin:      origin,
		confidence:  confidence,
		explanation: explanation,
		signals:     signals,
	}
}

func batchAnalyzeHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	samplesRaw, ok := input["samples"].([]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'samples' field")
	}

	var results []map[string]any
	aiCount, humanCount, unknownCount := 0, 0, 0
	var totalConfidence float64

	for i, s := range samplesRaw {
		text, ok := s.(string)
		if !ok {
			continue
		}

		result := deepAnalyze(text, "")
		results = append(results, map[string]any{
			"index":      i,
			"origin":     result.origin,
			"confidence": result.confidence,
			"signals":    result.signals,
		})

		totalConfidence += result.confidence
		switch result.origin {
		case "ai_system":
			aiCount++
		case "human":
			humanCount++
		default:
			unknownCount++
		}
	}

	total := len(results)
	avgConfidence := 0.0
	if total > 0 {
		avgConfidence = totalConfidence / float64(total)
	}

	return map[string]any{
		"total_samples":      total,
		"ai_count":           aiCount,
		"human_count":        humanCount,
		"unknown_count":      unknownCount,
		"ai_percentage":      float64(aiCount) / float64(maxInt(total, 1)) * 100,
		"human_percentage":   float64(humanCount) / float64(maxInt(total, 1)) * 100,
		"average_confidence": avgConfidence,
		"results":            results,
		"timestamp":          time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func contentFingerprintHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	metadata, _ := input["metadata"].(map[string]any)

	normalized := normalizeText(text)
	contentHash := sha256.Sum256([]byte(normalized))

	words := strings.Fields(normalized)
	var structureHash string
	if len(words) > 0 {
		structure := fmt.Sprintf("%d-%d-%d", len(words), len(splitSentences(text)), countParagraphs(text))
		h := sha256.Sum256([]byte(structure))
		structureHash = hex.EncodeToString(h[:8])
	}

	return map[string]any{
		"content_hash":   hex.EncodeToString(contentHash[:]),
		"structure_hash": structureHash,
		"word_count":     len(words),
		"char_count":     len(text),
		"metadata":       metadata,
		"generated_at":   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func compareTextsHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	textA, ok := input["text_a"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text_a' field")
	}

	textB, ok := input["text_b"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text_b' field")
	}

	profileA := generateStyleMetrics(textA)
	profileB := generateStyleMetrics(textB)

	similarity := calculateStyleSimilarity(profileA, profileB)

	resultA := deepAnalyze(textA, "")
	resultB := deepAnalyze(textB, "")

	sameOrigin := resultA.origin == resultB.origin

	return map[string]any{
		"style_similarity":   similarity,
		"same_origin_likely": sameOrigin && similarity > 0.7,
		"text_a_origin":      resultA.origin,
		"text_b_origin":      resultB.origin,
		"text_a_confidence":  resultA.confidence,
		"text_b_confidence":  resultB.confidence,
		"comparison_notes":   generateComparisonNotes(profileA, profileB, similarity),
		"timestamp":          time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func codeOriginHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	code, ok := input["code"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'code' field")
	}

	language, _ := input["language"].(string)

	var signals []string
	var aiScore, humanScore float64

	lines := strings.Split(code, "\n")
	commentLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
			commentLines++
		}
	}
	commentRatio := float64(commentLines) / float64(maxInt(len(lines), 1))
	if commentRatio > 0.3 {
		signals = append(signals, "high_comment_ratio")
		aiScore += 0.15
	}

	todoPattern := regexp.MustCompile(`(?i)(TODO|FIXME|XXX|HACK|BUG)`)
	if todoPattern.MatchString(code) {
		signals = append(signals, "todo_markers")
		humanScore += 0.1
	}

	docPatterns := regexp.MustCompile(`(?m)^\s*("""|'''|/\*\*|\* @)`)
	if docPatterns.MatchString(code) {
		signals = append(signals, "formal_documentation")
		aiScore += 0.1
	}

	magicNumbers := regexp.MustCompile(`\b\d{3,}\b`)
	if len(magicNumbers.FindAllString(code, -1)) > 3 {
		signals = append(signals, "magic_numbers")
		humanScore += 0.1
	}

	errorHandling := regexp.MustCompile(`(?i)(try|catch|except|error|err\s*!=|if err)`)
	if errorHandling.MatchString(code) {
		signals = append(signals, "error_handling")
	}

	origin := "unknown"
	confidence := 0.4
	explanation := "Code analysis inconclusive."

	if aiScore > humanScore {
		origin = "ai_generated"
		confidence = math.Min(0.9, 0.5+aiScore)
		explanation = fmt.Sprintf("AI code patterns detected: %v", signals)
	} else if humanScore > aiScore {
		origin = "human_written"
		confidence = math.Min(0.9, 0.5+humanScore)
		explanation = fmt.Sprintf("Human code patterns detected: %v", signals)
	}

	return map[string]any{
		"origin":        origin,
		"confidence":    confidence,
		"explanation":   explanation,
		"signals":       signals,
		"language":      language,
		"line_count":    len(lines),
		"comment_ratio": commentRatio,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func styleProfileHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	metrics := generateStyleMetrics(text)

	return map[string]any{
		"metrics":   metrics,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func generateStyleMetrics(text string) map[string]float64 {
	words := strings.Fields(text)
	sentences := splitSentences(text)

	avgWordLen := 0.0
	if len(words) > 0 {
		totalChars := 0
		for _, w := range words {
			totalChars += len(w)
		}
		avgWordLen = float64(totalChars) / float64(len(words))
	}

	avgSentLen := 0.0
	if len(sentences) > 0 {
		avgSentLen = float64(len(words)) / float64(len(sentences))
	}

	punctCount := 0
	for _, r := range text {
		if unicode.IsPunct(r) {
			punctCount++
		}
	}
	punctDensity := float64(punctCount) / float64(maxInt(len(text), 1))

	upperCount := 0
	for _, r := range text {
		if unicode.IsUpper(r) {
			upperCount++
		}
	}
	upperRatio := float64(upperCount) / float64(maxInt(len(text), 1))

	return map[string]float64{
		"avg_word_length":     avgWordLen,
		"avg_sentence_length": avgSentLen,
		"punctuation_density": punctDensity,
		"uppercase_ratio":     upperRatio,
		"word_count":          float64(len(words)),
		"sentence_count":      float64(len(sentences)),
	}
}

func calculateStyleSimilarity(a, b map[string]float64) float64 {
	keys := []string{"avg_word_length", "avg_sentence_length", "punctuation_density", "uppercase_ratio"}
	var totalDiff float64

	for _, k := range keys {
		valA, okA := a[k]
		valB, okB := b[k]
		if okA && okB {
			maxVal := math.Max(valA, valB)
			if maxVal > 0 {
				diff := math.Abs(valA-valB) / maxVal
				totalDiff += diff
			}
		}
	}

	similarity := 1 - (totalDiff / float64(len(keys)))
	return math.Max(0, math.Min(1, similarity))
}

func generateComparisonNotes(a, b map[string]float64, similarity float64) string {
	if similarity > 0.85 {
		return "Very similar writing styles suggest same author or same generation source."
	} else if similarity > 0.7 {
		return "Similar styles with minor variations."
	} else if similarity > 0.5 {
		return "Moderate style differences detected."
	}
	return "Significant style differences suggest different authors or sources."
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func splitSentences(text string) []string {
	re := regexp.MustCompile(`[.!?]+\s+`)
	parts := re.Split(text, -1)
	var sentences []string
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			sentences = append(sentences, p)
		}
	}
	return sentences
}

func calculateVariance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	return variance / float64(len(values))
}

func normalizeText(text string) string {
	text = strings.ToLower(text)
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func countParagraphs(text string) int {
	paragraphs := strings.Split(text, "\n\n")
	count := 0
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			count++
		}
	}
	return count
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
