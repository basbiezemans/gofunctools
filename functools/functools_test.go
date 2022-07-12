package functools

import (
	"math/rand"
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
	expect := 10
	result := Reduce(add, []int{1, 2, 3, 4}, 0)
	if result != expect {
		t.Errorf("Reduce(add, []int{1,2,3,4}, 0) = %d, expected %d", result, expect)
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

func BenchmarkApply(b *testing.B) {
	square := func(x int) int {
		return x ^ 2
	}
	for i := 0; i < b.N; i++ {
		Apply(square, rand.Perm(10000))
	}
}

func BenchmarkFilter(b *testing.B) {
	even := func(x int) bool {
		return x%2 == 0
	}
	for i := 0; i < b.N; i++ {
		Filter(even, rand.Perm(10000))
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
		date string  "ISO 8601"
		meas float64 "measurement"
	}
	makeDataPoint := func(date string, meas float64) DataPoint {
		return DataPoint{date, meas}
	}
	t2_slice1 := []string{"2021-01-15", "2021-01-16"}
	t2_slice2 := []float64{0.981, 0.973}
	t2_expect := []DataPoint{{"2021-01-15", 0.981}, {"2021-01-16", 0.973}}
	t2_result := ZipWith(makeDataPoint, t2_slice1, t2_slice2)
	if !reflect.DeepEqual(t2_result, t2_expect) {
		t.Errorf(`ZipWith(makeDataPoint, %v, %v) = %v, expected %v`, t2_slice1, t2_slice2, t2_result, t2_expect)
	}
}
