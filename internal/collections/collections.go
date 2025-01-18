// Package collections provides generic function utilities for working with collections.
package collections

import (
	"sort"
)

// Filter returns a slice of all elements of slice that pass the predicate.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	evenNumbers := Filter(numbers, func(n int) bool { return n%2 == 0 })
//	// evenNumbers will be []int{2, 4}
func Filter[T any](slice []T, pred func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if pred(v) {
			result = append(result, v)
		}
	}
	return result
}

// SortBy sorts a slice by the [less] function.
func SortBy[E any](slice []E, less func(a, b E) bool) {
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
}
