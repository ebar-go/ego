package helper

import "math"

//Min min
func Min(a, b int) int {
	if a < b {
		return a
	}
	if b < a {
		return b
	}
	return a
}

//Max max
func Max(a, b int) int {
	if a > b {
		return a
	}
	if b > a {
		return b
	}
	return a
}

// Div
func Div(a, b int) int {
	if b == 0 || a == 0 {
		return 0
	}

	return int(math.Ceil(float64(a) / float64(b)))
}

// Round 四舍五入取整
func Round(f float64) int {
	return int(math.Floor(f + 0.5))
}
