package utils

// Find Finds element in array.
// Return: index of found element, -1 if not found.
func Find[T comparable](a []T, elem T) int {
	for i, val := range a {
		if val == elem {
			return i
		}
	}
	return -1
}
