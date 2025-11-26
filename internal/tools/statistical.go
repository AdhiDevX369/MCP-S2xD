package tools

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

func RegisterStatisticalTools() {
	Register(Tool{
		Name:        "entropy_analyze",
		Description: "Calculate text entropy to detect AI-generated content. AI text typically has lower entropy.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to analyze for entropy patterns",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: entropyAnalyzeHandler,
	})

	Register(Tool{
		Name:        "ngram_analyze",
		Description: "Detect overused AI phrases and n-gram patterns that indicate AI generation.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to analyze for n-gram patterns",
						},
						"n": map[string]any{
							"type":        "integer",
							"description": "N-gram size (2-5, default 3)",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: ngramAnalyzeHandler,
	})

	Register(Tool{
		Name:        "burstiness_analyze",
		Description: "Measure token distribution uniformity. AI text has unnaturally uniform patterns.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to analyze for burstiness patterns",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: burstinessAnalyzeHandler,
	})

	Register(Tool{
		Name:        "zipf_analyze",
		Description: "Check adherence to Zipf's law. Natural language follows Zipf distribution; AI often deviates.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to analyze for Zipf's law compliance",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: zipfAnalyzeHandler,
	})

	Register(Tool{
		Name:        "statistical_detect",
		Description: "Comprehensive statistical AI detection combining entropy, n-gram, burstiness, and Zipf analysis.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to analyze using all statistical methods",
						},
					},
					"required": []string{"text"},
				},
			},
			"required": []string{"input"},
		},
		Handler: statisticalDetectHandler,
	})
}

func entropyAnalyzeHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	charEntropy := calcCharEntropy(text)
	wordEntropy := calcWordEntropy(text)
	bigramEntropy := calcBigramEntropy(text)
	avgEntropy := (charEntropy + wordEntropy + bigramEntropy) / 3

	origin := "unknown"
	confidence := 0.5
	explanation := "Entropy levels are inconclusive."

	if avgEntropy < 3.5 {
		origin = "ai_system"
		confidence = 0.7 + (3.5-avgEntropy)*0.1
		explanation = "Low entropy suggests repetitive, predictable AI-generated patterns."
	} else if avgEntropy > 4.5 {
		origin = "human"
		confidence = 0.6 + (avgEntropy-4.5)*0.08
		explanation = "High entropy indicates natural human variation."
	}

	return map[string]any{
		"origin":          origin,
		"confidence":      math.Min(0.95, confidence),
		"explanation":     explanation,
		"char_entropy":    charEntropy,
		"word_entropy":    wordEntropy,
		"bigram_entropy":  bigramEntropy,
		"average_entropy": avgEntropy,
		"word_count":      len(strings.Fields(text)),
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func calcCharEntropy(text string) float64 {
	if len(text) == 0 {
		return 0
	}
	freq := make(map[rune]int)
	total := 0
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			freq[r]++
			total++
		}
	}
	if total == 0 {
		return 0
	}
	entropy := 0.0
	for _, count := range freq {
		p := float64(count) / float64(total)
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func calcWordEntropy(text string) float64 {
	words := statTokenize(text)
	if len(words) == 0 {
		return 0
	}
	freq := make(map[string]int)
	for _, w := range words {
		freq[strings.ToLower(w)]++
	}
	entropy := 0.0
	total := float64(len(words))
	for _, count := range freq {
		p := float64(count) / total
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

func calcBigramEntropy(text string) float64 {
	words := statTokenize(text)
	if len(words) < 2 {
		return 0
	}
	freq := make(map[string]int)
	for i := 0; i < len(words)-1; i++ {
		bigram := strings.ToLower(words[i] + " " + words[i+1])
		freq[bigram]++
	}
	entropy := 0.0
	total := float64(len(words) - 1)
	for _, count := range freq {
		p := float64(count) / total
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}

var aiTypicalPhrases = []string{
	"it is important to note", "it's worth noting", "in conclusion", "to summarize",
	"first and foremost", "in today's world", "it is essential", "it is crucial",
	"plays a vital role", "a wide range of", "in order to", "due to the fact",
	"at the end of the day", "when it comes to", "on the other hand", "having said that",
	"it goes without saying", "needless to say", "as a matter of fact", "in light of",
	"with that being said", "it can be argued", "from this perspective", "taking into account",
	"in this regard", "for the purpose of", "in the context of", "it should be noted",
	"generally speaking", "as previously mentioned", "delve into", "dive deep",
	"leverage", "utilize", "facilitate", "implement", "comprehensive", "robust",
	"streamline", "optimize", "furthermore", "moreover", "additionally", "consequently",
}

func ngramAnalyzeHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	n := 3
	if nVal, ok := input["n"].(float64); ok {
		n = int(nVal)
		if n < 2 {
			n = 2
		}
		if n > 5 {
			n = 5
		}
	}

	words := statTokenize(text)
	ngrams := extractNgramsSlice(words, n)
	ngramFreq := make(map[string]int)
	for _, ng := range ngrams {
		ngramFreq[ng]++
	}

	var repeatedNgrams []map[string]any
	for ng, count := range ngramFreq {
		if count > 1 {
			repeatedNgrams = append(repeatedNgrams, map[string]any{"ngram": ng, "count": count})
		}
	}
	sort.Slice(repeatedNgrams, func(i, j int) bool {
		return repeatedNgrams[i]["count"].(int) > repeatedNgrams[j]["count"].(int)
	})
	if len(repeatedNgrams) > 10 {
		repeatedNgrams = repeatedNgrams[:10]
	}

	lower := strings.ToLower(text)
	var detectedPhrases []string
	aiPhraseScore := 0.0
	for _, phrase := range aiTypicalPhrases {
		count := strings.Count(lower, phrase)
		if count > 0 {
			detectedPhrases = append(detectedPhrases, phrase)
			aiPhraseScore += float64(count) * 0.1
		}
	}

	repetitionRate := 0.0
	if len(ngrams) > 0 {
		repeated := 0
		for _, count := range ngramFreq {
			if count > 1 {
				repeated += count - 1
			}
		}
		repetitionRate = float64(repeated) / float64(len(ngrams))
	}

	origin := "unknown"
	confidence := 0.5
	explanation := "N-gram patterns are inconclusive."
	combinedScore := aiPhraseScore + repetitionRate*2

	if combinedScore > 0.5 || len(detectedPhrases) >= 3 {
		origin = "ai_system"
		confidence = math.Min(0.95, 0.6+combinedScore*0.2)
		explanation = fmt.Sprintf("Detected %d AI-typical phrases and %.1f%% n-gram repetition.", len(detectedPhrases), repetitionRate*100)
	} else if repetitionRate < 0.05 && len(detectedPhrases) == 0 {
		origin = "human"
		confidence = 0.65
		explanation = "Low repetition and no AI-typical phrases detected."
	}

	return map[string]any{
		"origin":           origin,
		"confidence":       confidence,
		"explanation":      explanation,
		"n":                n,
		"total_ngrams":     len(ngrams),
		"unique_ngrams":    len(ngramFreq),
		"repetition_rate":  repetitionRate,
		"ai_phrases_found": detectedPhrases,
		"ai_phrase_score":  aiPhraseScore,
		"top_repeated":     repeatedNgrams,
		"timestamp":        time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func extractNgramsSlice(words []string, n int) []string {
	if len(words) < n {
		return nil
	}
	var ngrams []string
	for i := 0; i <= len(words)-n; i++ {
		ng := strings.Join(words[i:i+n], " ")
		ngrams = append(ngrams, strings.ToLower(ng))
	}
	return ngrams
}

func burstinessAnalyzeHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	sentences := statSplitSentences(text)
	if len(sentences) < 3 {
		return map[string]any{
			"origin":      "unknown",
			"confidence":  0.3,
			"explanation": "Insufficient text for burstiness analysis (need 3+ sentences).",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}, nil
	}

	sentenceLengths := make([]float64, len(sentences))
	for i, s := range sentences {
		sentenceLengths[i] = float64(len(strings.Fields(s)))
	}
	sentenceCV := coefficientOfVariation(sentenceLengths)

	words := statTokenize(text)
	wordLengths := make([]float64, len(words))
	for i, w := range words {
		wordLengths[i] = float64(len(w))
	}
	wordCV := coefficientOfVariation(wordLengths)

	paragraphs := strings.Split(text, "\n\n")
	var paragraphLengths []float64
	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			paragraphLengths = append(paragraphLengths, float64(len(strings.Fields(p))))
		}
	}
	paragraphCV := 0.0
	if len(paragraphLengths) > 1 {
		paragraphCV = coefficientOfVariation(paragraphLengths)
	}

	burstinessScore := (sentenceCV + wordCV + paragraphCV) / 3

	origin := "unknown"
	confidence := 0.5
	explanation := "Burstiness patterns are inconclusive."

	if burstinessScore < 0.3 {
		origin = "ai_system"
		confidence = 0.7 + (0.3-burstinessScore)*0.5
		explanation = fmt.Sprintf("Very uniform text structure (burstiness=%.2f) suggests AI.", burstinessScore)
	} else if burstinessScore > 0.5 {
		origin = "human"
		confidence = 0.6 + (burstinessScore-0.5)*0.4
		explanation = fmt.Sprintf("High variation (burstiness=%.2f) indicates human writing.", burstinessScore)
	}

	return map[string]any{
		"origin":              origin,
		"confidence":          math.Min(0.95, confidence),
		"explanation":         explanation,
		"burstiness_score":    burstinessScore,
		"sentence_cv":         sentenceCV,
		"word_length_cv":      wordCV,
		"paragraph_cv":        paragraphCV,
		"sentence_count":      len(sentences),
		"avg_sentence_length": statMean(sentenceLengths),
		"timestamp":           time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func zipfAnalyzeHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	words := statTokenize(text)
	if len(words) < 20 {
		return map[string]any{
			"origin":      "unknown",
			"confidence":  0.3,
			"explanation": "Insufficient text for Zipf analysis (need 20+ words).",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}, nil
	}

	freq := make(map[string]int)
	for _, w := range words {
		freq[strings.ToLower(w)]++
	}

	type wf struct {
		word  string
		count int
	}
	var sorted []wf
	for w, c := range freq {
		sorted = append(sorted, wf{w, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	topN := statMinInt(len(sorted), 20)
	var deviations []float64
	for i := 0; i < topN; i++ {
		rank := float64(i + 1)
		expected := float64(sorted[0].count) / rank
		actual := float64(sorted[i].count)
		deviation := math.Abs(actual-expected) / expected
		deviations = append(deviations, deviation)
	}
	avgDeviation := statMean(deviations)

	var logRanks, logFreqs []float64
	for i := 0; i < topN; i++ {
		logRanks = append(logRanks, math.Log(float64(i+1)))
		logFreqs = append(logFreqs, math.Log(float64(sorted[i].count)))
	}
	slope, rSquared := statLinearRegression(logRanks, logFreqs)

	origin := "unknown"
	confidence := 0.5
	explanation := "Zipf's law analysis is inconclusive."

	if avgDeviation > 0.4 || rSquared < 0.8 {
		origin = "ai_system"
		confidence = 0.6 + avgDeviation*0.3
		explanation = fmt.Sprintf("Word frequency deviates from Zipf's law (R²=%.2f).", rSquared)
	} else if rSquared > 0.9 && avgDeviation < 0.25 {
		origin = "human"
		confidence = 0.65 + rSquared*0.2
		explanation = fmt.Sprintf("Word frequency follows natural Zipf distribution (R²=%.2f).", rSquared)
	}

	topWords := make([]map[string]any, statMinInt(len(sorted), 10))
	for i := 0; i < len(topWords); i++ {
		topWords[i] = map[string]any{"word": sorted[i].word, "count": sorted[i].count, "rank": i + 1}
	}

	return map[string]any{
		"origin":        origin,
		"confidence":    math.Min(0.95, confidence),
		"explanation":   explanation,
		"zipf_slope":    slope,
		"r_squared":     rSquared,
		"avg_deviation": avgDeviation,
		"unique_words":  len(freq),
		"total_words":   len(words),
		"top_words":     topWords,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func statisticalDetectHandler(params map[string]any) (map[string]any, error) {
	input, ok := params["input"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'input' parameter")
	}

	text, ok := input["text"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'text' field")
	}

	entropyResult, _ := entropyAnalyzeHandler(params)
	ngramResult, _ := ngramAnalyzeHandler(params)
	burstinessResult, _ := burstinessAnalyzeHandler(params)
	zipfResult, _ := zipfAnalyzeHandler(params)

	weights := map[string]float64{"entropy": 0.25, "ngram": 0.30, "burstiness": 0.25, "zipf": 0.20}
	scores := make(map[string]float64)

	extractScore := func(result map[string]any, method string) {
		if result == nil {
			return
		}
		origin, _ := result["origin"].(string)
		conf, _ := result["confidence"].(float64)
		if origin == "ai_system" {
			scores[method] = conf
		} else if origin == "human" {
			scores[method] = -conf
		}
	}

	extractScore(entropyResult, "entropy")
	extractScore(ngramResult, "ngram")
	extractScore(burstinessResult, "burstiness")
	extractScore(zipfResult, "zipf")

	weightedScore := 0.0
	for method, score := range scores {
		weightedScore += score * weights[method]
	}

	var signals []string
	if entropyResult != nil && entropyResult["origin"] == "ai_system" {
		signals = append(signals, "low_entropy")
	}
	if ngramResult != nil && ngramResult["origin"] == "ai_system" {
		signals = append(signals, "ai_phrases")
	}
	if burstinessResult != nil && burstinessResult["origin"] == "ai_system" {
		signals = append(signals, "uniform_structure")
	}
	if zipfResult != nil && zipfResult["origin"] == "ai_system" {
		signals = append(signals, "zipf_deviation")
	}

	origin := "unknown"
	confidence := 0.5
	explanation := "Statistical analysis is inconclusive."

	if weightedScore > 0.2 {
		origin = "ai_system"
		confidence = math.Min(0.95, 0.5+weightedScore)
		explanation = fmt.Sprintf("Statistical analysis indicates AI generation. Signals: %v", signals)
	} else if weightedScore < -0.2 {
		origin = "human"
		confidence = math.Min(0.95, 0.5-weightedScore)
		explanation = "Statistical analysis indicates human authorship."
	}

	return map[string]any{
		"origin":            origin,
		"confidence":        confidence,
		"explanation":       explanation,
		"weighted_score":    weightedScore,
		"signals":           signals,
		"method_scores":     scores,
		"entropy_result":    entropyResult,
		"ngram_result":      ngramResult,
		"burstiness_result": burstinessResult,
		"zipf_result":       zipfResult,
		"word_count":        len(strings.Fields(text)),
		"timestamp":         time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func statTokenize(text string) []string {
	re := regexp.MustCompile(`[a-zA-Z]+`)
	return re.FindAllString(text, -1)
}

func statSplitSentences(text string) []string {
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

func statMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func statVariance(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := statMean(values)
	variance := 0.0
	for _, v := range values {
		variance += (v - m) * (v - m)
	}
	return variance / float64(len(values))
}

func coefficientOfVariation(values []float64) float64 {
	m := statMean(values)
	if m == 0 {
		return 0
	}
	return math.Sqrt(statVariance(values)) / m
}

func statMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func statLinearRegression(x, y []float64) (slope, rSquared float64) {
	n := float64(len(x))
	if n < 2 {
		return 0, 0
	}
	sumX, sumY, sumXY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}
	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, 0
	}
	slope = (n*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / n
	ssTotal := sumY2 - (sumY*sumY)/n
	ssRes := 0.0
	for i := range x {
		pred := slope*x[i] + intercept
		ssRes += (y[i] - pred) * (y[i] - pred)
	}
	if ssTotal == 0 {
		rSquared = 1
	} else {
		rSquared = 1 - ssRes/ssTotal
	}
	return slope, rSquared
}
