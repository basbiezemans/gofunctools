// Package maps defines various functions useful with maps of any type.
package maps

// ToSlice, applied to a hash map and a combiner function, combines
// key-value pairs as elements of a new slice.
func ToSlice[A comparable, B, C any](fn func(A, B) C, hm map[A]B) []C {
	var xs = make([]C, 0, len(hm))
	for k, v := range hm {
		xs = append(xs, fn(k, v))
	}
	return xs
}

// Deprecated: use function ToSlice.
func HashMapToSlice[A comparable, B, C any](fn func(A, B) C, hm map[A]B) []C {
	return ToSlice(fn, hm)
}

// Map applies a binary function to each key-value pair of a hash map.
func Map[A, B comparable, C, D any](fn func(A, C) (B, D), hm map[A]C) map[B]D {
	var nm = make(map[B]D, len(hm))
	for k, v := range hm {
		n, m := fn(k, v)
		nm[n] = m
	}
	return nm
}
