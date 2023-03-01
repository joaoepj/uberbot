package internal

// global utility functions for the bot

func Filter[T any](items []T, Func func(item T) bool) []T {
	var filteredItems []T
	if len(items) == 0 {
		return filteredItems
	}
	for _, value := range items {
		if Func(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

// ToPtr makes a given value into it's pointer form
func ToPtr[T any](v T) *T {
	return &v
}

//
//func Find[T any](items []T, fn func(item T) bool) T {
//	foundItem := T{}
//	if len(items) == 0 {
//		return foundItem
//	}
//	for _, value := range items {
//		if fn(value) {
//			foundItem = value
//		}
//	}
//	return foundItem
//}
