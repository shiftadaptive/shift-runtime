// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package models

type Request struct {
	Target string                 `json:"target"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
	Body   interface{}            `json:"body"`
}