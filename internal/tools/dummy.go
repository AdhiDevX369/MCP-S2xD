package tools

func RegisterDefaults() {
    Register(Tool{
        Name:        "dummy_echo",
        Description: "Echo back input text for testing.",
        Handler: func(params map[string]any) (map[string]any, error) {
            in := params["input"].(map[string]any)
            return map[string]any{
                "echo": in,
            }, nil
        },
    })
}
