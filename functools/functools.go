package functools

func Any[A any](fn func(A) bool, xs []A) bool {
	for _, x := range xs {
		if fn(x) {
			return true
		}
	}
	return false
}

func All[A any](fn func(A) bool, xs []A) bool {
	for _, x := range xs {
		if !fn(x) {
			return false
		}
	}
	return true
}

func Reduce[A, B any](fn func(B, A) B, xs []A, initValue B) B {
	acc := initValue
	for _, x := range xs {
		acc = fn(acc, x)
	}
	return acc
}

func Apply[A, B any](fn func(A) B, xs []A) []B {
	reducer := func(ys []B, x A) []B {
		return append(ys, fn(x))
	}
	return Reduce(reducer, xs, make([]B, 0))
}

func Filter[A any](fn func(A) bool, xs []A) []A {
	reducer := func(ys []A, x A) []A {
		if fn(x) {
			return append(ys, x)
		}
		return ys
	}
	return Reduce(reducer, xs, make([]A, 0))
}
