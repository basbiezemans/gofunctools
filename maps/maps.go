// Package maps defines various functions useful with maps of any type.
package maps

// HashMapToSlice, applied to a hash map and a combiner function, combines
// key-value pairs as elements of a new slice.
func HashMapToSlice[K comparable, V, T any](fn func(K, V) T, hm map[K]V) []T {
	var xs = make([]T, 0, len(hm))
	for k, v := range hm {
		xs = append(xs, fn(k, v))
	}
	return xs
}
