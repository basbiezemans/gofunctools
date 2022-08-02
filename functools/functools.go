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

// Pipe chains multiple unary functions together. All functions have to accept
// and return a value of the same type in order to pipe it from one function
// to the next. Functions are evaluated from left to right.
func Pipe[A any](funcs ...func(A) A) func(A) A {
	return func(x A) A {
		var acc = x
		for _, fn := range funcs {
			acc = fn(acc)
		}
		return acc
	}
}

// Compose combines two simple, unary functions into a more complicated one.
// Functions are evaluated from right to left (f1 after f2).
func Compose[A, B, C any](f1 func(B) C, f2 func(A) B) func(A) C {
	return func(x A) C {
		return f1(f2(x))
	}
}

// Curry2 converts an uncurried, binary function to a curried function.
func Curry2[A, B, C any](fn func(A, B) C) func(A) func(B) C {
	return func(x A) func(B) C {
		return func(y B) C {
			return fn(x, y)
		}
	}
}

// Curry3 converts an uncurried, ternary function to a curried function.
func Curry3[A, B, C, D any](fn func(A, B, C) D) func(A) func(B) func(C) D {
	return func(x A) func(B) func(C) D {
		return func(y B) func(C) D {
			return func(z C) D {
				return fn(x, y, z)
			}
		}
	}
}

// Partial1 takes a binary function and a value, and returns a unary function
// as its result.
func Partial1[A, B, C any](fn func(A, B) C, x A) func(B) C {
	return func(y B) C {
		return fn(x, y)
	}
}

// Partial2 takes a ternary function and two values, and returns a unary
// function as its result.
func Partial2[A, B, C, D any](fn func(A, B, C) D, x A, y B) func(C) D {
	return func(z C) D {
		return fn(x, y, z)
	}
}

// Flip converts a binary function to a function with the order of arguments
// flipped.
func Flip[A, B, C any](fn func(A, B) C) func(B, A) C {
	return func(x B, y A) C {
		return fn(y, x)
	}
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

// Returns the smaller of its two arguments.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
