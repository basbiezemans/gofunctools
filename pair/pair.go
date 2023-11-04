package pair

type Pair[A, B any] struct {
	fst A
	snd B
}

// Create a new Pair.
func New[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{a, b}
}

// Extract a pair into its components.
func Unpair[A, B any](p Pair[A, B]) (A, B) {
	return p.fst, p.snd
}

// Duplicate a single value into a pair.
func Dupe[A any](a A) Pair[A, A] {
	return New[A, A](a, a)
}

// Swap the components of a pair.
func Swap[A, B any](p Pair[A, B]) Pair[B, A] {
	return New[B, A](p.snd, p.fst)
}

// Update the first component of a pair.
func First[A, B any](fn func(A) A, p Pair[A, B]) Pair[A, B] {
	return New[A, B](fn(p.fst), p.snd)
}

// Update the second component of a pair.
func Second[A, B any](fn func(B) B, p Pair[A, B]) Pair[A, B] {
	return New[A, B](p.fst, fn(p.snd))
}

// Apply a single function to both components of a pair.
func Both[A, B any](fn func(A) B, p Pair[A, A]) Pair[B, B] {
	return New[B, B](fn(p.fst), fn(p.snd))
}

// Given two functions, apply both to a single value to form a pair.
func Fanout[A, B, C any](f1 func(A) B, f2 func(A) C, x A) Pair[B, C] {
	return New[B, C](f1(x), f2(x))
}

// Given two functions, apply one to the first component and one to the second.
func Bimap[A, B, C, D any](f1 func(A) B, f2 func(C) D, p Pair[A, C]) Pair[B, D] {
	return New[B, D](f1(p.fst), f2(p.snd))
}

// Extract the first component of a pair.
func (p Pair[A, B]) Fst() A {
	return p.fst
}

// Extract the second component of a pair.
func (p Pair[A, B]) Snd() B {
	return p.snd
}

// Unpair as method.
func (p Pair[A, B]) Unpair() (A, B) {
	return p.fst, p.snd
}

// Swap as method.
func (p Pair[A, B]) Swap() Pair[B, A] {
	return New[B, A](p.snd, p.fst)
}

// First as method.
func (p Pair[A, B]) First(fn func(A) A) Pair[A, B] {
	return New[A, B](fn(p.fst), p.snd)
}

// Second as method.
func (p Pair[A, B]) Second(fn func(B) B) Pair[A, B] {
	return New[A, B](p.fst, fn(p.snd))
}
