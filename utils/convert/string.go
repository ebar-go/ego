package convert

import "strconv"

func ToString(a any) string {
	switch a.(type) {
	case string:
		return a.(string)
	case int:
		return strconv.Itoa(a.(int))
	}

	return ""
}
