package utils

func ArrayGet[T any](array []T, i int) T {
	if i < 0 || i >= len(array) {
		var zero T
		return zero
	}
	return array[i]
}
