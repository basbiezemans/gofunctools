package iters

import (
	"errors"
	"iter"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	expect := []int{2, 4, 6, 8}
	result := slices.Collect(Map(double, slices.Values(numbers)))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, expected %v", result, expect)
	}
}

func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	expect := []int{2, 4}
	result := slices.Collect(Filter(even, slices.Values(numbers)))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, expected %v", result, expect)
	}
}

func TestDropWhile(t *testing.T) {
	data := []int{1, 2, 3, 4}
	input := []int{3, 8, 1}
	expect := [][]int{
		{3, 4}, nil, {1, 2, 3, 4},
	}
	var result []int
	for i, n := range input {
		result = slices.Collect(DropWhile(lessThan(n), slices.Values(data)))
		if !reflect.DeepEqual(result, expect[i]) {
			t.Errorf("result = %v, expected %v", result, expect[i])
		}
	}
}

func TestTakeWhile(t *testing.T) {
	data := []int{1, 2, 3, 4}
	input := []int{3, 0, 5}
	expect := [][]int{
		{1, 2}, nil, {1, 2, 3, 4},
	}
	var result []int
	for i, n := range input {
		result = slices.Collect(TakeWhile(lessThan(n), slices.Values(data)))
		if !reflect.DeepEqual(result, expect[i]) {
			t.Errorf("result = %v, expected %v", result, expect[i])
		}
	}
}

func TestZipWith(t *testing.T) {
	testcases := []map[string][]int{
		{"nums1": {1, 2, 3, 4}, "nums2": {1, 2, 3, 4, 5}, "expect": {1, 4, 9, 16}},
		{"nums1": {1, 2, 3, 4}, "nums2": {}, "expect": nil},
	}
	for _, tc := range testcases {
		it1 := slices.Values(tc["nums1"])
		it2 := slices.Values(tc["nums2"])
		result := slices.Collect(ZipWith(multiply, it1, it2))
		if !reflect.DeepEqual(result, tc["expect"]) {
			t.Errorf("result = %v, expected %v", result, tc["expect"])
		}
	}
	type DataPoint struct {
		date string
		meas float64
	}
	makeDataPoint := func(date string, meas float64) DataPoint {
		return DataPoint{date, meas}
	}
	it1 := slices.Values([]string{"2021-01-15", "2021-01-16"})
	it2 := slices.Values([]float64{0.981, 0.973})
	expect := []DataPoint{
		{"2021-01-15", 0.981}, {"2021-01-16", 0.973},
	}
	result := slices.Collect(ZipWith(makeDataPoint, it1, it2))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, expected %v", result, expect)
	}
}

func TestUnzipWith(t *testing.T) {
	type DataPoint struct {
		date string
		meas float64
	}
	split := func(datum DataPoint) (string, float64) {
		return datum.date, datum.meas
	}
	datapoints := []DataPoint{
		{"2021-01-15", 0.981}, {"2021-01-16", 0.973},
	}
	expect1 := []string{"2021-01-15", "2021-01-16"}
	expect2 := []float64{0.981, 0.973}

	result1, result2 := collect(UnzipWith(split, slices.Values(datapoints)))
	if !reflect.DeepEqual(result1, expect1) || !reflect.DeepEqual(result2, expect2) {
		t.Errorf("result = %v, %v, expected %v, %v", result1, result2, expect1, expect2)
	}
}

func TestUnfold(t *testing.T) {
	expect := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	result := slices.Collect(Unfold(decrement, 10))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result = %v, expected %v", result, expect)
	}
}

func TestScan(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		init   int
		input  []int
		expect []int
	}
	testcases := []TestCase{
		{add, 0, []int{1, 2, 3, 4}, []int{0, 1, 3, 6, 10}},
		{add, 42, []int{}, []int{42}},
		{subtract, 100, []int{1, 2, 3, 4}, []int{100, 99, 97, 94, 90}},
	}
	errorMsg := "Scan(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		it := slices.Values(test.input)
		result := slices.Collect(Scan(test.callb, test.init, it))
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf(errorMsg, funcName(test.callb), test.init, test.input, result, test.expect)
		}
	}
	prepend := func(s string, r rune) string {
		return string(r) + s
	}
	expect := []string{"foo", "afoo", "bafoo", "cbafoo", "dcbafoo"}
	init := "foo"
	input := []rune{'a', 'b', 'c', 'd'}
	result := slices.Collect(Scan(prepend, init, slices.Values(input)))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(errorMsg, "prepend", init, showRunes(input), result, expect)
	}
}

// Helper functions

func even(x int) bool {
	return x%2 == 0
}

func double(x int) int {
	return 2 * x
}

func add(x, y int) int {
	return x + y
}

func multiply(x, y int) int {
	return x * y
}

func subtract(x, y int) int {
	return x - y
}

func lessThan(y int) func(int) bool {
	return func(x int) bool {
		return x < y
	}
}

func decrement(x int) (int, int, bool) {
	if x > 0 {
		return x, x - 1, true
	}
	return 0, -1, false
}

func showRunes(rs []rune) []string {
	var ys = make([]string, len(rs))
	for i, r := range rs {
		ys[i] = string(r)
	}
	return ys
}

func last[T any](xs []T) (T, error) {
	var zero T
	if len(xs) > 0 {
		return xs[len(xs)-1], nil
	}
	return zero, errors.New("empty slice")
}

func funcName[T, U any](fn func(T, U) T) string {
	var fptr = reflect.ValueOf(fn).Pointer()
	var fname = runtime.FuncForPC(fptr).Name()
	if str, err := last(strings.Split(fname, ".")); err == nil {
		return str
	}
	return "N/A"
}

func collect[T, U any](seq iter.Seq2[T, U]) ([]T, []U) {
	var xs = make([]T, 0)
	var ys = make([]U, 0)
	for x, y := range seq {
		xs = append(xs, x)
		ys = append(ys, y)
	}
	return xs, ys
}
