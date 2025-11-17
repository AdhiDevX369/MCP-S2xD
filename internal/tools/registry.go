package tools

import "fmt"

type ToolHandler func(params map[string]any) (map[string]any, error)

type Tool struct {
    Name        string
    Description string
    Handler     ToolHandler
}

var registry = map[string]Tool{}

func Register(t Tool) {
    registry[t.Name] = t
}
func List() []Tool {
    out := make([]Tool, 0, len(registry))
    for _, t := range registry {
        out = append(out, t)
    }
    return out
}
func Call(name string, params map[string]any) (map[string]any, error) {
    t, ok := registry[name]
    if !ok {
        return nil, fmt.Errorf("unknown tool: %s", name)
    }
    return t.Handler(params)
}
