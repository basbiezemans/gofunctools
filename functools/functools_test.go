package functools

import (
	"reflect"
	"testing"
)

func TestAny(t *testing.T) {
	even := func(x int) bool {
		return x%2 == 0
	}
	expect := true
	result := Any(even, []int{1, 2, 3, 4})
	if result != expect {
		t.Errorf("Any(even, []int{1,2,3,4}) = %t, expected %t", result, expect)
	}
	expect = false
	result = Any(even, []int{1, 3, 5, 7})
	if result != expect {
		t.Errorf("Any(even, []int{1,3,5,7}) = %t, expected %t", result, expect)
	}
}

func TestAll(t *testing.T) {
	even := func(x int) bool {
		return x%2 == 0
	}
	expect := true
	result := All(even, []int{2, 4, 6, 8})
	if result != expect {
		t.Errorf("All(even, []int{2,4,6,8}) = %t, expected %t", result, expect)
	}
	expect = false
	result = All(even, []int{2, 4, 7, 8})
	if result != expect {
		t.Errorf("All(even, []int{2,4,7,8}) = %t, expected %t", result, expect)
	}
}

func TestReduce(t *testing.T) {
	add := func(x, y int) int {
		return x + y
	}
	t1_expect := 10
	t1_result := Reduce(add, []int{1, 2, 3, 4}, 0)
	if t1_result != t1_expect {
		t.Errorf("Reduce(add, []int{1,2,3,4}, 0) = %d, expected %d", t1_result, t1_expect)
	}
	even := func(x int) bool {
		return x%2 == 0
	}
	// Remove negative numbers and split odd numbers into an even number and 1
	posEvens := func(ys []int, x int) []int {
		switch {
		case x < 0:
			return ys
		case even(x):
			return append(ys, x)
		default:
			return append(ys, x-1, 1)
		}
	}
	t2_values := []int{5, 4, -3, 20, 17, -33, -4, 18}
	t2_expect := []int{4, 1, 4, 20, 16, 1, 18}
	t2_result := Reduce(posEvens, t2_values, []int{})
	if !reflect.DeepEqual(t2_result, t2_expect) {
		t.Errorf("Reduce(posEvens, %v, []int{}) = %v, expected %v", t2_values, t2_result, t2_expect)
	}

}

func TestReduceRight(t *testing.T) {
	append := func(ys []int, x int) []int {
		return append(ys, x)
	}
	data := []int{1, 2, 3, 4}
	expect := []int{4, 3, 2, 1}
	result := ReduceRight(append, data, []int{})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ReduceRight(append, %v, []int{}) = %d, expected %d", data, result, expect)
	}
}

func TestApply(t *testing.T) {
	double := func(x int) int {
		return 2 * x
	}
	expect := []int{2, 4, 6, 8}
	result := Apply(double, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Apply(double, []int{1,2,3,4}) = %v, expected %v", result, expect)
	}
}

func TestFilter(t *testing.T) {
	even := func(x int) bool {
		return x%2 == 0
	}
	expect := []int{2, 4}
	result := Filter(even, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Filter(even, []int{1,2,3,4}) = %v, expected %v", result, expect)
	}
}

func TestDropWhile(t *testing.T) {
	lessThan := func(y int) func(x int) bool {
		return func(x int) bool {
			return x < y
		}
	}
	data := []int{1, 2, 3, 4}
	input := []int{3, 8, 1}
	expect := [][]int{
		{3, 4}, {}, {1, 2, 3, 4},
	}
	var result []int
	for i, n := range input {
		result = DropWhile(lessThan(n), data)
		if !reflect.DeepEqual(result, expect[i]) {
			t.Errorf("DropWhile(lessThan(%d), %v) = %v, expected %v", n, data, result, expect[i])
		}
	}
}

func TestTakeWhile(t *testing.T) {
	lessThan := func(y int) func(x int) bool {
		return func(x int) bool {
			return x < y
		}
	}
	data := []int{1, 2, 3, 4}
	input := []int{3, 0, 5}
	expect := [][]int{
		{1, 2}, {}, {1, 2, 3, 4},
	}
	var result []int
	for i, n := range input {
		result = TakeWhile(lessThan(n), data)
		if !reflect.DeepEqual(result, expect[i]) {
			t.Errorf("TakeWhile(lessThan(%d), %v) = %v, expected %v", n, data, result, expect[i])
		}
	}
}

func TestZipWith(t *testing.T) {
	multiply := func(x, y int) int {
		return x * y
	}
	t1_slice1 := []int{1, 2, 3, 4}
	t1_slice2 := []int{1, 2, 3, 4, 5}
	t1_expect := []int{1, 4, 9, 16}
	t1_result := ZipWith(multiply, t1_slice1, t1_slice2)
	if !reflect.DeepEqual(t1_result, t1_expect) {
		t.Errorf("ZipWith(multiply, %v, %v) = %v, expected %v", t1_slice1, t1_slice2, t1_result, t1_expect)
	}
	type DataPoint struct {
		date string
		meas float64
	}
	makeDataPoint := func(date string, meas float64) DataPoint {
		return DataPoint{date, meas}
	}
	t2_slice1 := []string{"2021-01-15", "2021-01-16"}
	t2_slice2 := []float64{0.981, 0.973}
	t2_expect := []DataPoint{
		{"2021-01-15", 0.981}, {"2021-01-16", 0.973},
	}
	t2_result := ZipWith(makeDataPoint, t2_slice1, t2_slice2)
	if !reflect.DeepEqual(t2_result, t2_expect) {
		t.Errorf("ZipWith(makeDataPoint, %v, %v) = %v, expected %v", t2_slice1, t2_slice2, t2_result, t2_expect)
	}
}
