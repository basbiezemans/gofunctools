// This package defines various functions that correspond to operators of Go.
// They are meant to be used with higher-order functions like Filter, Reduce, etc.
package operators

import (
	"errors"
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

// Modulo (remainder)
func Modulo[T constraints.Integer](x, y T) T {
	return x % y
}

// Determine if a given integer is an even number.
func Even[T constraints.Integer](x T) bool {
	return x%2 == 0
}

// Determine if a given integer is an odd number.
func Odd[T constraints.Integer](x T) bool {
	return x%2 != 0
}

// Addition / string concatenation.
func Add[T Number | ~string](x, y T) T {
	return x + y
}

// Multiplication
func Multiply[T Number](x, y T) T {
	return x * y
}

// Subtraction
func Subtract[T Number](x, y T) T {
	return x - y
}

// Division. The divisor must not be zero. If the divisor is zero at run time,
// a run-time panic occurs. Function SafeDiv will prevent this from happening.
func Divide[T Number](x, y T) T {
	return x / y
}

// Safe Division
func SafeDiv[T Number](x, y T) (T, error) {
	if y == 0 {
		return 0, errors.New("division by zero")
	}
	return x / y, nil
}

// Pow returns x ** y, the base-x exponential of y.
func Pow[T Number](x, y T) float64 {
	base := float64(x)
	exponent := float64(y)
	return math.Pow(base, exponent)
}

// Logical AND
func AND(x, y bool) bool {
	return x && y
}

// Logical OR
func OR(x, y bool) bool {
	return x || y
}

// Equality
func Equal[T constraints.Ordered](x, y T) bool {
	return x == y
}

// Difference
func NotEqual[T constraints.Ordered](x, y T) bool {
	return x != y
}

// Ordering: less than
func LessThan[T constraints.Ordered](x, y T) bool {
	return x < y
}

// Ordering: greater than
func GreaterThan[T constraints.Ordered](x, y T) bool {
	return x > y
}

// Ordering: less than or equal
func LessThanOrEqual[T constraints.Ordered](x, y T) bool {
	return x <= y
}

// Ordering: greater than or equal
func GreaterThanOrEqual[T constraints.Ordered](x, y T) bool {
	return x >= y
}
