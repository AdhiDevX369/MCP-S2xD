package rpcerror

import "fmt"

type Code int

const (
	ParseError     Code = -32700
	InvalidRequest Code = -32600
	MethodNotFound Code = -32601
	InvalidParams  Code = -32602
	InternalError  Code = -32603
	ToolNotFound   Code = -32000
	ToolExecError  Code = -32001
	Unauthorized   Code = -32002
	RateLimited    Code = -32003
)

type RPCError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error %d: %s", e.Code, e.Message)
}

func New(code Code, message string) *RPCError {
	return &RPCError{Code: code, Message: message}
}

func NewWithData(code Code, message string, data any) *RPCError {
	return &RPCError{Code: code, Message: message, Data: data}
}

func Parse(msg string) *RPCError {
	return New(ParseError, msg)
}

func InvalidReq(msg string) *RPCError {
	return New(InvalidRequest, msg)
}

func NotFound(method string) *RPCError {
	return New(MethodNotFound, fmt.Sprintf("method not found: %s", method))
}

func InvalidParam(msg string) *RPCError {
	return New(InvalidParams, msg)
}

func Internal(msg string) *RPCError {
	return New(InternalError, msg)
}

func ToolExec(toolName, msg string) *RPCError {
	return NewWithData(ToolExecError, msg, map[string]string{"tool": toolName})
}

func (e *RPCError) ToMap() map[string]any {
	m := map[string]any{
		"code":    e.Code,
		"message": e.Message,
	}
	if e.Data != nil {
		m["data"] = e.Data
	}
	return m
}
