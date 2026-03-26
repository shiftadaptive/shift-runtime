// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package utils

func ConvertParams(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for k, v := range params {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}

	return result
}