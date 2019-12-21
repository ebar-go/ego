package helper

// DefaultInt if v is zero, return defaultV
func DefaultInt(v int, defaultV int) int {
	if v == 0 {
		return defaultV
	}

	return v
}

// DefaultString if v is empty, return defaultV
func DefaultString(v string, defaultV string) string  {
	if v == "" {
		return defaultV
	}

	return v
}