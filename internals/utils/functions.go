package utils

func IsItemInCollection[T comparable](item T, collection []T) bool {
	for _, s := range collection {
		if s == item {
			return true
		}
	}
	return false
}
