// SHIFT ::: Runtime
// Lightweight adaptive middleware for API compatibility
// (c) 2026 ShiftAdaptive

package cache

var paramMapping = make(map[string]string)

func StoreMapping(original string, corrected string) {
	paramMapping[original] = corrected
}

func GetMapping(original string) (string, bool) {
	val, exists := paramMapping[original]
	return val, exists
}
