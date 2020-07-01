package strings

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

// UUID return unique id
func UUID() string {
	return uuid.NewV4().String()
}


// Default return defaultV if v is empty
func Default(v, defaultV string) string {
	if v == "" {
		return defaultV
	}

	return v
}

// ToBool check bool string
func ToBool(b string) bool {
	if b == "1" || "true" == strings.ToLower(b) {
		return true
	}

	return false
}

// ToInt return int of number string
func ToInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}
