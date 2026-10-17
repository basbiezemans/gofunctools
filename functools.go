// Package functools provides generic higher-order functions.
package functools

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"sync"
)

// Any, applied to a predicate and a slice, determines whether any element of
// the slice satisfies the predicate.
func Any[A any](fn func(A) bool, xs []A) bool {
	return slices.ContainsFunc(xs, fn)
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
	for i := range xs {
		acc = fn(acc, xs[i])
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

// Reduce applies a reducer function left-to-right to a non-empty slice, starting
// with the first element as the initial accumulator value. It returns the reduced
// value, or an error when the slice is empty or if the reducer panics.
func Reduce[A any](fn func(A, A) A, xs []A) (result A, err error) {
	if len(xs) == 0 {
		return result, errors.New("empty slice")
	}
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("reducer panicked: %v", rec)
		}
	}()
	var acc = xs[0]
	var n = len(xs)
	for i := 1; i < n; i++ {
		acc = fn(acc, xs[i])
	}
	return acc, nil
}

// ReduceLeft, applied to a reducer function and a non-empty slice, reduces the
// slice to a single value. The accumulated value must be of the same type as
// the slice elements. This function is non-total and will panic if the slice
// happens to be empty.
//
// Deprecated: Use Reduce or FoldLeft instead.
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
//
// Deprecated: Use FoldRight instead.
func ReduceRight[A any](fn func(A, A) A, xs []A) A {
	if len(xs) == 0 {
		panic("empty slice")
	}
	var n = len(xs) - 1
	return FoldRight(fn, xs[n], xs[:n])
}

// ParallelReduce, applied to a reducer function and a slice, reduces the slice
// to a single value in parallel. The result is equivalent to Reduce only when
// the reducer is associative.
//   - The reducer function must be associative because the parallel
//     implementation may change the grouping of values.
//   - The reducer function must not depend on shared mutable state.
//   - An error is returned when the slice is empty or if the reducer panics.
func ParallelReduce[A any](fn func(A, A) A, xs []A) (result A, err error) {
	if len(xs) == 0 {
		return result, errors.New("empty slice")
	}
	var maxProcs = min(runtime.GOMAXPROCS(0)*4, len(xs))
	var chunkSize = ceildiv(len(xs), maxProcs)
	var chunkCount = ceildiv(len(xs), chunkSize)
	var errCh = make(chan error, 1)
	var slice = make([]A, chunkCount)
	var index = 0
	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for chunk := range slices.Chunk(xs, chunkSize) {
		go func(chunk []A, chunkIndex int) {
			defer wg.Done()
			if val, err := Reduce(fn, chunk); err == nil {
				slice[chunkIndex] = val
			} else {
				select {
				case errCh <- fmt.Errorf("reducer panicked: %v", err):
				default: // don’t block on a full channel
				}
			}
		}(chunk, index)
		index++
	}
	wg.Wait()
	select {
	case err := <-errCh:
		return result, err
	default: // no panic to report
	}
	return Reduce(fn, slice)
}

// Map applies a unary function to each element of a slice.
func Map[A, B any](fn func(A) B, xs []A) []B {
	var ys = make([]B, len(xs))
	for i, x := range xs {
		ys[i] = fn(x)
	}
	return ys
}

// MapMaybe applies a unary function, which can fail, to each element of
// a slice and discards results for which the function returns an error.
func MapMaybe[A, B any](fn func(A) (B, error), xs []A) []B {
	var ys = make([]B, 0)
	for _, x := range xs {
		if y, err := fn(x); err == nil {
			ys = append(ys, y)
		}
	}
	return ys
}

// TryMap applies a unary function, which can fail, to each element of a
// slice. If the function succeeds for all elements then a slice with the
// new values is returned. If the function fails for any of the elements
// then the error is returned immediately. None of the slice elements are
// processed after one returns an error.
func TryMap[A, B any](fn func(A) (B, error), xs []A) ([]B, error) {
	var ys = make([]B, 0, len(xs))
	for _, x := range xs {
		if y, err := fn(x); err != nil {
			return nil, err
		} else {
			ys = append(ys, y)
		}
	}
	return ys, nil
}

// ParallelMap applies a unary function to each element of a slice in parallel.
//   - ParallelMap is a drop-in replacement for Map. It incurs slightly more
//     overhead than Map, but is faster with larger inputs and a CPU-intensive
//     mapper function.
//   - If a mapper function triggers a panic, the output slice will have the
//     default (zero) value for that specific element.
func ParallelMap[A, B any](fn func(A) B, xs []A) []B {
	if len(xs) == 0 {
		return []B{}
	}
	var maxProcs = min(runtime.GOMAXPROCS(0)*4, len(xs))
	var chunkSize = ceildiv(len(xs), maxProcs)
	var chunkCount = ceildiv(len(xs), chunkSize)
	var slice = make([]B, len(xs))
	var index = 0
	var recovery = func() { recover() }
	var wg sync.WaitGroup
	wg.Add(chunkCount)
	for chunk := range slices.Chunk(xs, chunkSize) {
		go func(chunk []A, chunkIndex int) {
			defer wg.Done()
			defer recovery()
			offset := chunkIndex * chunkSize
			for i, x := range chunk {
				slice[offset+i] = fn(x)
			}
		}(chunk, index)
		index++
	}
	wg.Wait()
	return slice
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
			return slices.Clone(xs[i:])
		}
	}
	return []A{}
}

// TakeWhile, applied to a predicate and a slice, takes elements from the slice
// as long as the predicate is true.
func TakeWhile[A any](fn func(A) bool, xs []A) []A {
	for i, x := range xs {
		if !fn(x) {
			return slices.Clone(xs[:i])
		}
	}
	return slices.Clone(xs)
}

// ZipWith, applied to a combiner function and two slices, combines elements
// from the two slices using the combiner function. If one input slice is
// shorter than the other, excess elements of the longer slice are discarded.
func ZipWith[A, B, C any](fn func(A, B) C, xs []A, ys []B) []C {
	var n = min(len(xs), len(ys))
	var zs = make([]C, n)
	for i := range n {
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

// Curry1 converts an uncurried, unary function to a curried function.
func Curry1[A, B any](fn func(A) B) func(A) B {
	return func(x A) B {
		return fn(x)
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

// Uncurry converts a curried binary function to a regular binary function.
func Uncurry[A, B, C any](fn func(A) func(B) C) func(A, B) C {
	return func(x A, y B) C {
		return fn(x)(y)
	}
}

// Partial1 takes a binary function and one value, and returns a unary
// function as its result.
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

// Partial1Pair takes a binary function — which returns two values — and a
// single value as an argument, and returns a unary function as its result.
// It is a modified version of Partial1 that allows for the partial
// application of Partition and UnzipWith.
func Partial1Pair[A, B, C, D any](fn func(A, B) (C, D), x A) func(B) (C, D) {
	return func(y B) (C, D) {
		return fn(x, y)
	}
}

// Flip converts a binary function to a function with the order of arguments
// flipped.
func Flip[A, B, C any](fn func(A, B) C) func(B, A) C {
	return func(x B, y A) C {
		return fn(y, x)
	}
}

// Rotate takes a ternary function and rotates the arguments to the left by
// one position.
func Rotate[A, B, C, D any](fn func(A, B, C) D) func(B, C, A) D {
	return func(x B, y C, z A) D {
		return fn(z, x, y)
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

// Unfold, builds a slice from a seed value. The function takes the current
// input value as seed and returns (output, next-seed, ok). If ok is false,
// the slice is complete; if ok is true, the output value is appended to the
// slice and the seed value is used next.
func Unfold[A, B any](fn func(B) (A, B, bool), seedValue B) []A {
	var xs = make([]A, 0)
	var x, next, ok = fn(seedValue)
	for ok {
		xs = append(xs, x)
		x, next, ok = fn(next)
	}
	return xs
}

// Find, takes a predicate and a slice and returns the first element in the
// slice matching the predicate (a,true), or (zero,false) if there is no such
// element. The element type must be comparable.
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
// there is no such element. The element type must be comparable.
func FindIndex[A comparable](fn func(A) bool, xs []A) (int, bool) {
	for i, x := range xs {
		if fn(x) {
			return i, true
		}
	}
	return -1, false
}

// FindIndices, extends FindIndex by returning the indices of all elements
// satisfying the predicate, in ascending order. The element type must be
// comparable.
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
	for i, x := range slices.Backward(xs) {
		ys[i] = fn(x, ys[i+1])
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

// GroupBy groups consecutive elements of a slice into subslices, using the
// specified equivalence relation to determine equality. Equality is checked
// against the first element in the group, not against the nearest neighbour.
// The element type must be comparable.
func GroupBy[A comparable](eq func(A, A) bool, xs []A) [][]A {
	var ys = make([][]A, 0)
	if len(xs) == 0 {
		return ys
	}
	var first = xs[0]
	var group = []A{first}
	for _, x := range xs[1:] {
		if eq(first, x) {
			group = append(group, x)
		} else {
			ys = append(ys, group)
			first = x
			group = []A{first}
		}
	}
	return append(ys, group)
}

// DeleteBy removes the first occurrence of an element from a slice using an
// equality predicate. The element type must be comparable.
func DeleteBy[A comparable](eq func(A, A) bool, x A, xs []A) []A {
	var ys = slices.Clone(xs)
	for i, y := range ys {
		if eq(x, y) {
			return append(ys[:i], ys[i+1:]...)
		}
	}
	return ys
}

// Returns the smallest integer greater than or equal to x/y
// Example: 7/3 = (7 + 3 − 1) / 3 = 9/3 = 3
func ceildiv(x, y int) int {
	return (x + y - 1) / y
}
