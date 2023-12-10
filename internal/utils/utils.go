package utils

func Contains[T comparable](elements []T, target T) bool {
	for idx := range elements {
		if elements[idx] == target {
			return true
		}
	}

	return false
}
