package library

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
