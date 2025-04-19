package utils

// Keys returns a slice containing all the keys in the map.
//
// Parameters:
//   - m map[K]V: The input map.
//
// Returns:
//   - []K: A slice containing all the keys in the map.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice containing all the values in the map.
//
// Parameters:
//   - m map[K]V: The input map.
//
// Returns:
//   - []V: A slice containing all the values in the map.
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Entries returns a slice of key-value pairs from the map.
//
// Parameters:
//   - m map[K]V: The input map.
//
// Returns:
//   - []struct{Key K; Value V}: A slice of key-value pairs.
func Entries[K comparable, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	entries := make([]struct {
		Key   K
		Value V
	}, 0, len(m))

	for k, v := range m {
		entries = append(entries, struct {
			Key   K
			Value V
		}{k, v})
	}

	return entries
}

// HasKey checks if a map contains a specific key.
//
// Parameters:
//   - m map[K]V: The input map.
//   - key K: The key to check for.
//
// Returns:
//   - bool: True if the map contains the key, false otherwise.
func HasKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Merge combines two maps into a new map.
// If a key exists in both maps, the value from the second map is used.
//
// Parameters:
//   - m1 map[K]V: The first input map.
//   - m2 map[K]V: The second input map.
//
// Returns:
//   - map[K]V: A new map containing all key-value pairs from both input maps.
func Merge[K comparable, V any](m1, m2 map[K]V) map[K]V {
	result := make(map[K]V, len(m1)+len(m2))

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}

// Pick creates a new map with only the specified keys from the input map.
//
// Parameters:
//   - m map[K]V: The input map.
//   - keys []K: The keys to include in the new map.
//
// Returns:
//   - map[K]V: A new map containing only the specified keys and their values.
func Pick[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}

	return result
}

// Omit creates a new map without the specified keys from the input map.
//
// Parameters:
//   - m map[K]V: The input map.
//   - keys []K: The keys to exclude from the new map.
//
// Returns:
//   - map[K]V: A new map containing all key-value pairs except those with the specified keys.
func Omit[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if !Contains(keys, k) {
			result[k] = v
		}
	}

	return result
}

// MapKeys creates a new map by applying a function to each key in the input map.
//
// Parameters:
//   - m map[K]V: The input map.
//   - mapper func(K) L: A function that transforms keys of type K to type L.
//
// Returns:
//   - map[L]V: A new map with transformed keys and the same values.
func MapKeys[K comparable, L comparable, V any](m map[K]V, mapper func(K) L) map[L]V {
	result := make(map[L]V, len(m))

	for k, v := range m {
		result[mapper(k)] = v
	}

	return result
}

// MapValues creates a new map by applying a function to each value in the input map.
//
// Parameters:
//   - m map[K]V: The input map.
//   - mapper func(V) W: A function that transforms values of type V to type W.
//
// Returns:
//   - map[K]W: A new map with the same keys and transformed values.
func MapValues[K comparable, V any, W any](m map[K]V, mapper func(V) W) map[K]W {
	result := make(map[K]W, len(m))

	for k, v := range m {
		result[k] = mapper(v)
	}

	return result
}

// FilterKeys creates a new map with only the key-value pairs whose keys satisfy the predicate.
//
// Parameters:
//   - m map[K]V: The input map.
//   - predicate func(K) bool: A function that returns true for keys to include.
//
// Returns:
//   - map[K]V: A new map containing only the key-value pairs whose keys satisfy the predicate.
func FilterKeys[K comparable, V any](m map[K]V, predicate func(K) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if predicate(k) {
			result[k] = v
		}
	}

	return result
}

// FilterValues creates a new map with only the key-value pairs whose values satisfy the predicate.
//
// Parameters:
//   - m map[K]V: The input map.
//   - predicate func(V) bool: A function that returns true for values to include.
//
// Returns:
//   - map[K]V: A new map containing only the key-value pairs whose values satisfy the predicate.
func FilterValues[K comparable, V any](m map[K]V, predicate func(V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if predicate(v) {
			result[k] = v
		}
	}

	return result
}
