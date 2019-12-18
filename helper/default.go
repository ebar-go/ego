package helper

// DefaultInt if v is zero, return defaultV
func DefaultInt(v int, defaultV int) int {
	if v == 0 {
		return defaultV
	}

	return v
}
