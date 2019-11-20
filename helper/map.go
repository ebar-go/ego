package helper

// MergeMaps 合并
func MergeMaps(items ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, item := range items {
		for key, value := range item {
			result[key] = value
		}
	}

	return result
}
