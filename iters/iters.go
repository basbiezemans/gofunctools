// This package defines common higher-order functions that work on iterators.
package iters

import (
	"iter"
)

// Map applies a unary function to each element of an iterator.
func Map[A, B any](fn func(A) B, seq iter.Seq[A]) iter.Seq[B] {
	return func(yield func(B) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Filter, applied to a predicate and an iterator, filters the elements that
// satisfy the predicate.
func Filter[A any](fn func(A) bool, seq iter.Seq[A]) iter.Seq[A] {
	return func(yield func(A) bool) {
		for v := range seq {
			if fn(v) && !yield(v) {
				return
			}
		}
	}
}

// DropWhile, applied to a predicate and an iterator, drops elements as long
// as the predicate is true.
func DropWhile[A any](fn func(A) bool, seq iter.Seq[A]) iter.Seq[A] {
	return func(yield func(A) bool) {
		predicate_holds := true
		for v := range seq {
			if predicate_holds && fn(v) {
				continue // drop element
			} else {
				predicate_holds = false
				if !yield(v) {
					return
				}
			}
		}
	}
}

// TakeWhile, applied to a predicate and an iterator, takes elements as long
// as the predicate is true.
func TakeWhile[A any](fn func(A) bool, seq iter.Seq[A]) iter.Seq[A] {
	return func(yield func(A) bool) {
		for v := range seq {
			if fn(v) {
				if !yield(v) {
					return
				}
			} else {
				break
			}
		}
	}
}

// ZipWith, applied to a combiner function and two iterators, combines their
// elements using the combiner function.
func ZipWith[A, B, C any](fn func(A, B) C, seq1 iter.Seq[A], seq2 iter.Seq[B]) iter.Seq[C] {
	return func(yield func(C) bool) {
		next, stop := iter.Pull(seq2)
		defer stop()
		for v1 := range seq1 {
			v2, ok := next()
			if !ok {
				return
			}
			if !yield(fn(v1, v2)) {
				return
			}
		}
	}
}

// UnzipWith, applied to a splitter function and an iterator, splits elements
// into two parts using the splitter function.
func UnzipWith[A, B, C any](fn func(A) (B, C), seq iter.Seq[A]) iter.Seq2[B, C] {
	return func(yield func(B, C) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Unfold, applied to a build function and a seed value, produces an iterator
// from the seed value until the build function returns false.
func Unfold[A, B any](fn func(B) (A, B, bool), initValue B) iter.Seq[A] {
	return func(yield func(A) bool) {
		var v, next, ok = fn(initValue)
		for ok {
			if !yield(v) {
				return
			}
			v, next, ok = fn(next)
		}
	}
}

// Scan, applied to a reducer function and an iterator, produces an iterator
// of successive reduced values.
func Scan[A, B any](fn func(B, A) B, initValue B, seq iter.Seq[A]) iter.Seq[B] {
	return func(yield func(B) bool) {
		if !yield(initValue) {
			return
		}
		var acc = initValue
		for v := range seq {
			acc = fn(acc, v)
			if !yield(acc) {
				return
			}
		}
	}
}
