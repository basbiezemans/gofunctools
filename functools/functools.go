// Package functools provides generic higher-order functions.
package functools

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

// Reduce, applied to a reducer function and a slice, reduces the slice to a
// single value. It is also known as a left fold, because it "folds" a data
// structure from left to right.
func Reduce[A, B any](fn func(B, A) B, xs []A, initValue B) B {
	acc := initValue
	for _, x := range xs {
		acc = fn(acc, x)
	}
	return acc
}

// Apply (a.k.a. "map"), applies a unary function to each element of a slice.
func Apply[A, B any](fn func(A) B, xs []A) []B {
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

// ZipWith, applied to a binary function and two slices, combines elements from
// the two slices in a pairwise fashion. If one input slice is shorter than the
// other, excess elements of the longer slice are discarded.
func ZipWith[A, B, C any](fn func(A, B) C, xs []A, ys []B) []C {
	var n = min(len(xs), len(ys))
	var zs = make([]C, n)
	for i := 0; i < n; i++ {
		zs[i] = fn(xs[i], ys[i])
	}
	return zs
}

// Returns the smaller of its two arguments.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
