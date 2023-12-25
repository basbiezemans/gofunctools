// This package defines common higher-order functions that are useful with
// slices of any type.
package slices

// Any, applied to a predicate and a slice, determines whether any element of
// the slice satisfies the predicate.
func Any[A any](fn func(A) bool, xs []A) bool {
	for _, x := range xs {
		if fn(x) {
			return true
		}
	}
	return false
}

// All, applied to a predicate and a slice, determines whether all elements of
// the slice satisfy the predicate.
func All[A any](fn func(A) bool, xs []A) bool {
	for _, x := range xs {
		if !fn(x) {
			return false
		}
	}
	return true
}

// FoldLeft, applied to a reducer function, an initialization value and a slice,
// reduces the slice to a single value. This function "folds" a slice from left
// to right, starting with the initialization value.
func FoldLeft[A, B any](fn func(B, A) B, initValue B, xs []A) B {
	var acc = initValue
	for _, x := range xs {
		acc = fn(acc, x)
	}
	return acc
}

// FoldRight, applied to a reducer function, an initialization value and a slice,
// reduces the slice to a single value. This function "folds" a slice from right
// to left, starting with the initialization value.
func FoldRight[A, B any](fn func(A, B) B, initValue B, xs []A) B {
	var acc = initValue
	var n = len(xs) - 1
	for i := range xs {
		acc = fn(xs[n-i], acc)
	}
	return acc
}

// ReduceLeft, applied to a reducer function and a non-empty slice, reduces the
// slice to a single value. The accumulated value must be of the same type as
// the slice elements. This function is non-total and will panic if the slice
// happens to be empty.
func ReduceLeft[A any](fn func(A, A) A, xs []A) A {
	if len(xs) == 0 {
		panic("empty slice")
	}
	return FoldLeft(fn, xs[0], xs[1:])
}

// ReduceRight, applied to a reducer function and a non-empty slice, reduces
// the slice to a single value. The accumulated value must be of the same type
// as the slice elements. This function is non-total and will panic if the
// slice happens to be empty.
func ReduceRight[A any](fn func(A, A) A, xs []A) A {
	if len(xs) == 0 {
		panic("empty slice")
	}
	var n = len(xs) - 1
	return FoldRight(fn, xs[n], xs[:n])
}

// Map applies a unary function to each element of a slice.
func Map[A, B any](fn func(A) B, xs []A) []B {
	var ys = make([]B, len(xs))
	for i, x := range xs {
		ys[i] = fn(x)
	}
	return ys
}

// Filter, applied to a predicate and a slice, filters the slice of those
// elements that satisfy the predicate.
func Filter[A any](fn func(A) bool, xs []A) []A {
	var ys = make([]A, 0, len(xs)/2)
	for _, x := range xs {
		if fn(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

// DropWhile, applied to a predicate and a slice, drops elements from the slice
// as long as the predicate is true.
func DropWhile[A any](fn func(A) bool, xs []A) []A {
	for i, x := range xs {
		if !fn(x) {
			return xs[i:]
		}
	}
	return []A{}
}

// TakeWhile, applied to a predicate and a slice, takes elements from the slice
// as long as the predicate is true.
func TakeWhile[A any](fn func(A) bool, xs []A) []A {
	for i, x := range xs {
		if !fn(x) {
			return xs[:i]
		}
	}
	return xs
}

// ZipWith, applied to a combiner function and two slices, combines elements
// from the two slices using the combiner function. If one input slice is
// shorter than the other, excess elements of the longer slice are discarded.
func ZipWith[A, B, C any](fn func(A, B) C, xs []A, ys []B) []C {
	var n = min(len(xs), len(ys))
	var zs = make([]C, n)
	for i := 0; i < n; i++ {
		zs[i] = fn(xs[i], ys[i])
	}
	return zs
}

// UnzipWith, applied to a splitter function and a slice, splits elements in
// the slice into two parts using the splitter function. Each part is directed
// to a different slice.
func UnzipWith[A, B, C any](fn func(A) (B, C), xs []A) ([]B, []C) {
	var ys = make([]B, len(xs))
	var zs = make([]C, len(xs))
	for i, x := range xs {
		ys[i], zs[i] = fn(x)
	}
	return ys, zs
}

// Partition takes a predicate and a slice, and splits the elements into two
// slices which do and do not satisfy the predicate.
func Partition[A any](fn func(A) bool, xs []A) ([]A, []A) {
	var ys = make([]A, 0, len(xs)/2)
	var zs = make([]A, 0, len(xs)/2)
	for _, x := range xs {
		if fn(x) {
			ys = append(ys, x)
		} else {
			zs = append(zs, x)
		}
	}
	return ys, zs
}

// Count, applied to a predicate and a slice, counts the number of elements
// that satisfy the predicate.
func Count[A any](fn func(A) bool, xs []A) int {
	var k = 0
	for _, x := range xs {
		if fn(x) {
			k += 1
		}
	}
	return k
}

// ToHashMap, applied to a slice and a splitter function, creates a
// hash map with elements from the slice splitted as key-value pairs.
func ToHashMap[A comparable, B, C any](fn func(B) (A, C), xs []B) map[A]C {
	var hm = make(map[A]C, len(xs))
	for _, x := range xs {
		k, v := fn(x)
		hm[k] = v
	}
	return hm
}

// Deprecated: use function ToHashMap.
func SliceToHashMap[A comparable, B, C any](fn func(B) (A, C), xs []B) map[A]C {
	return ToHashMap(fn, xs)
}

// Unfold, builds a slice from a seed value. The function takes the element and
// returns (,,false) if it is done producing the slice or returns (a,b,true),
// in which case, a is appended to the slice and b is used as the next element.
func Unfold[A, B any](fn func(B) (A, B, bool), initValue B) []A {
	var xs = make([]A, 0)
	var x, next, ok = fn(initValue)
	for ok {
		xs = append(xs, x)
		x, next, ok = fn(next)
	}
	return xs
}

// Find, takes a predicate and a slice and returns the first element in the slice
// matching the predicate (a,true), or (zero,false) if there is no such element.
func Find[A comparable](fn func(A) bool, xs []A) (A, bool) {
	var zero A
	for _, x := range xs {
		if fn(x) {
			return x, true
		}
	}
	return zero, false
}

// FindIndex, takes a predicate and a slice and returns the index of the first
// element in the slice satisfying the predicate (index,true), or (-1,false) if
// there is no such element.
func FindIndex[A comparable](fn func(A) bool, xs []A) (int, bool) {
	for i, x := range xs {
		if fn(x) {
			return i, true
		}
	}
	return -1, false
}

// FindIndices, extends FindIndex, by returning the indices of all elements
// satisfying the predicate, in ascending order.
func FindIndices[A comparable](fn func(A) bool, xs []A) []int {
	var ys = make([]int, 0)
	for i, x := range xs {
		if fn(x) {
			ys = append(ys, i)
		}
	}
	return ys
}

// ScanLeft is similar to FoldLeft, but returns a slice of successive reduced
// values from the left.
func ScanLeft[A, B any](fn func(B, A) B, initValue B, xs []A) []B {
	var ys = make([]B, 1+len(xs))
	ys[0] = initValue
	for i, x := range xs {
		ys[i+1] = fn(ys[i], x)
	}
	return ys
}

// ScanRight is the right-to-left dual of ScanLeft. Note that the order of
// parameters on the accumulating function are reversed compared to ScanLeft.
func ScanRight[A, B any](fn func(A, B) B, initValue B, xs []A) []B {
	var n = len(xs)
	var ys = make([]B, 1+n)
	ys[n] = initValue
	for i := n; i > 0; i-- {
		ys[i-1] = fn(xs[i-1], ys[i])
	}
	return ys
}

// ConcatMap applies a function, returning a slice, over a slice and concatenates
// the results.
func ConcatMap[A, B any](fn func(A) []B, xs []A) []B {
	var ys = make([]B, 0, len(xs))
	for _, x := range xs {
		ys = append(ys, fn(x)...)
	}
	return ys
}

// Returns the smaller of its two arguments.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
