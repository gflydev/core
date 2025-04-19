package utils

// Contains checks if a slice contains a specific element.
//
// Parameters:
//   - slice []T: The slice to search in.
//   - element T: The element to search for.
//
// Returns:
//   - bool: True if the element is found in the slice, false otherwise.
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// IndexOf finds the index of an element in a slice.
//
// Parameters:
//   - slice []T: The slice to search in.
//   - element T: The element to search for.
//
// Returns:
//   - int: The index of the element in the slice, or -1 if not found.
func IndexOf[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

// Filter creates a new slice containing only the elements that satisfy the predicate function.
//
// Parameters:
//   - slice []T: The input slice to filter.
//   - predicate func(T) bool: A function that returns true for elements to include.
//
// Returns:
//   - []T: A new slice containing only the elements for which the predicate returns true.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Map creates a new slice by applying a function to each element of the input slice.
//
// Parameters:
//   - slice []T: The input slice to transform.
//   - mapper func(T) U: A function that transforms elements of type T to type U.
//
// Returns:
//   - []U: A new slice containing the transformed elements.
func Map[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = mapper(v)
	}
	return result
}

// Reduce applies a function to each element in the slice, accumulating a single result.
//
// Parameters:
//   - slice []T: The input slice to reduce.
//   - initialValue U: The initial value for the accumulator.
//   - reducer func(U, T) U: A function that combines the accumulator with each element.
//
// Returns:
//   - U: The final accumulated value.
func Reduce[T any, U any](slice []T, initialValue U, reducer func(U, T) U) U {
	result := initialValue
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// ForEach applies a function to each element in the slice.
//
// Parameters:
//   - slice []T: The input slice to iterate over.
//   - action func(T): A function to apply to each element.
func ForEach[T any](slice []T, action func(T)) {
	for _, v := range slice {
		action(v)
	}
}

// Any returns true if any element in the slice satisfies the predicate.
//
// Parameters:
//   - slice []T: The input slice to check.
//   - predicate func(T) bool: A function that returns true for matching elements.
//
// Returns:
//   - bool: True if any element satisfies the predicate, false otherwise.
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// All returns true if all elements in the slice satisfy the predicate.
//
// Parameters:
//   - slice []T: The input slice to check.
//   - predicate func(T) bool: A function that returns true for matching elements.
//
// Returns:
//   - bool: True if all elements satisfy the predicate, false otherwise.
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Unique returns a new slice with duplicate elements removed.
//
// Parameters:
//   - slice []T: The input slice that may contain duplicates.
//
// Returns:
//   - []T: A new slice with duplicate elements removed.
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	var result []T

	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Chunk splits a slice into chunks of the specified size.
//
// Parameters:
//   - slice []T: The input slice to split.
//   - size int: The size of each chunk.
//
// Returns:
//   - [][]T: A slice of slices, where each inner slice has at most 'size' elements.
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Reverse returns a new slice with the elements in reverse order.
//
// Parameters:
//   - slice []T: The input slice to reverse.
//
// Returns:
//   - []T: A new slice with the elements in reverse order.
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Shuffle returns a new slice with the elements randomly shuffled.
//
// Parameters:
//   - slice []T: The input slice to shuffle.
//
// Returns:
//   - []T: A new slice with the elements randomly shuffled.
func Shuffle[T any](slice []T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for i := len(result) - 1; i > 0; i-- {
		j := int(RandInt64(int64(i + 1)))
		result[i], result[j] = result[j], result[i]
	}

	return result
}
