package lib

func CreateCopy[T any](items []T) []T {
	newCopy := make([]T, len(items))
	copy(newCopy, items)

	return newCopy
}
