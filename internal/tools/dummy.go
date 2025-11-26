package tools

import "fmt"

func RegisterDefaults() {
	Register(Tool{
		Name:        "dummy_echo",
		Description: "Echo back input text for testing.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "object",
					"description": "Any JSON object to echo back",
				},
			},
			"required": []string{"input"},
		},
		Handler: func(params map[string]any) (map[string]any, error) {
			in, ok := params["input"].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("missing or invalid 'input' parameter")
			}
			return map[string]any{
				"echo": in,
			}, nil
		},
	})
}
