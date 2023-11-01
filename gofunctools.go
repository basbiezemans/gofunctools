// Package gofunctools defines a number of functions that are useful when
// programming in a functional style.
package gofunctools

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
