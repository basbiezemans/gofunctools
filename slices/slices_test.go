package slices

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/basbiezemans/gofunctools/pair"
)

func TestAny(t *testing.T) {
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

func TestReduceLeft(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	expect := -8
	result := ReduceLeft(subtract, numbers)
	if result != expect {
		t.Errorf("ReduceLeft(subtract, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestReduceRight(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	expect := -2
	result := ReduceRight(subtract, numbers)
	if result != expect {
		t.Errorf("ReduceRight(subtract, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestFoldLeft(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		init   int
		input  []int
		expect int
	}
	testcases := []TestCase{
		{add, 0, []int{1, 2, 3, 4}, 10},
		{add, 42, []int{}, 42},
		{subtract, 100, []int{1, 2, 3, 4}, 90},
	}
	errorMsg := "FoldLeft(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		result := FoldLeft(test.callb, test.init, test.input)
		if result != test.expect {
			fn := funcName(test.callb)
			t.Errorf(errorMsg, fn, test.init, test.input, result, test.expect)
		}
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
	values := []int{5, 4, -3, 20, 17, -33, -4, 18}
	expect := []int{4, 1, 4, 20, 16, 1, 18}
	result := FoldLeft(posEvens, []int{}, values)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("FoldLeft(posEvens, []int{}, %v) = %v, expected %v", values, result, expect)
	}
}

func TestFoldRight(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		init   int
		input  []int
		expect int
	}
	testcases := []TestCase{
		{add, 0, []int{1, 2, 3, 4}, 10},
		{add, 42, []int{}, 42},
		{subtract, 100, []int{1, 2, 3, 4}, 98},
	}
	errorMsg := "FoldRight(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		result := FoldRight(test.callb, test.init, test.input)
		if result != test.expect {
			fn := funcName(test.callb)
			t.Errorf(errorMsg, fn, test.init, test.input, result, test.expect)
		}
	}
	append := func(x int, ys []int) []int {
		return append(ys, x)
	}
	values := []int{1, 2, 3, 4}
	expect := []int{4, 3, 2, 1}
	reverse := func(numbers []int) []int {
		return FoldRight(append, []int{}, numbers)
	}
	result := reverse(values)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("FoldRight(append, []int{}, %v) = %v, expected %v", values, result, expect)
	}
}

func TestMap(t *testing.T) {
	expect := []int{2, 4, 6, 8}
	result := Map(double, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Map(double, []int{1,2,3,4}) = %v, expected %v", result, expect)
	}
}

func TestMapMaybe(t *testing.T) {
	expect := []int{1, 2, 3, 4}
	slices := [][]int{
		{1}, {}, {1, 2}, {1, 2, 3}, {}, {4},
	}
	result := MapMaybe(last, slices)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("MapMaybe(last, [][]int{...}) = %v, expected %v", result, expect)
	}
}

func TestFilter(t *testing.T) {
	expect := []int{2, 4}
	result := Filter(even, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Filter(even, []int{1,2,3,4}) = %v, expected %v", result, expect)
	}
}

func TestDropWhile(t *testing.T) {
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
	testcases := []map[string][]int{
		{"nums1": {1, 2, 3, 4}, "nums2": {1, 2, 3, 4, 5}, "expect": {1, 4, 9, 16}},
		{"nums1": {1, 2, 3, 4}, "nums2": {}, "expect": {}},
	}
	for _, tc := range testcases {
		result := ZipWith(multiply, tc["nums1"], tc["nums2"])
		if !reflect.DeepEqual(result, tc["expect"]) {
			t.Errorf("ZipWith(multiply, %v, %v) = %v, expected %v", tc["nums1"], tc["nums2"], result, tc["expect"])
		}
	}
	type DataPoint struct {
		date string
		meas float64
	}
	makeDataPoint := func(date string, meas float64) DataPoint {
		return DataPoint{date, meas}
	}
	slice1 := []string{"2021-01-15", "2021-01-16"}
	slice2 := []float64{0.981, 0.973}
	expect := []DataPoint{
		{"2021-01-15", 0.981}, {"2021-01-16", 0.973},
	}
	result := ZipWith(makeDataPoint, slice1, slice2)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ZipWith(makeDataPoint, %v, %v) = %v, expected %v", slice1, slice2, result, expect)
	}
}

func TestZipLongest(t *testing.T) {
	expect := []int{2, 4, 3, 4}
	result := ZipLongest(add, []int{1, 2}, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ZipLongest(add, []int{1,2}, []int{1,2,3,4}) = %v, expected %v", result, expect)
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

	result1, result2 := UnzipWith(split, datapoints)
	if !reflect.DeepEqual(result1, expect1) || !reflect.DeepEqual(result2, expect2) {
		t.Errorf("UnZipWith(split, %v) = %v, %v, expected %v, %v", datapoints, result1, result2, expect1, expect2)
	}
}

func TestPartition(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect1 := []int{0, 2, 4, 6, 8}
	expect2 := []int{1, 3, 5, 7, 9}
	result1, result2 := Partition(even, numbers)
	if !reflect.DeepEqual(result1, expect1) || !reflect.DeepEqual(result2, expect2) {
		t.Errorf("Partition(even, %v) = %v, %v, expected %v, %v", numbers, result1, result2, expect1, expect2)
	}
}

func TestCount(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect := 5
	result := Count(even, numbers)
	if result != expect {
		t.Errorf("Count(even, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestToHashMap(t *testing.T) {
	items := []pair.Pair[string, int]{
		pair.New("lorem", 1),
		pair.New("ipsum", 2),
		pair.New("dolor", 3),
	}
	expect := map[string]int{
		"lorem": 1, "ipsum": 2, "dolor": 3,
	}
	result := ToHashMap(pair.Unpair, items)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ToHashMap(pair.Unpair, %v) = %v, expected %v", items, result, expect)
	}
}

func TestUnfold(t *testing.T) {
	expect := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	result := Unfold(decrement, 10)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Unfold(decrement, 10) = %v, expected %v", result, expect)
	}
}

func BenchmarkUnfold(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unfold(decrement, 1000)
	}
}

func TestFind(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n, ok := Find(greaterThan(4), input)
	if !ok || n != 5 {
		t.Errorf("Find(greaterThan(4), %v) = %d, expected 5", input, n)
	}
	n, ok = Find(lessThan(0), input)
	if ok || n != 0 {
		t.Errorf("Find(lessThan(0), %v) = %t, expected false", input, ok)
	}
}

func TestFindIndex(t *testing.T) {
	hello := "Hello World!"
	input := []rune(hello)
	i, ok := FindIndex(unicode.IsSpace, input)
	if !ok || i != 5 {
		t.Errorf("FindIndex(IsSpace, %q) = %d, expected 5", hello, i)
	}
	i, ok = FindIndex(unicode.IsDigit, input)
	if ok || i != -1 {
		t.Errorf("FindIndex(IsDigit, %q) = %t, expected false", hello, ok)
	}
}

func TestFindIndices(t *testing.T) {
	vowels := []rune("aeiou")
	isVowel := isOneOf(vowels)
	hello := "Hello World!"
	expect := []int{1, 4, 7}
	result := FindIndices(isVowel, []rune(hello))
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("FindIndices(isVowel, %q) = %v, expected %v", hello, result, expect)
	}
}

func TestScanLeft(t *testing.T) {
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
	errorMsg := "ScanLeft(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		result := ScanLeft(test.callb, test.init, test.input)
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
	result := ScanLeft(prepend, init, input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(errorMsg, "prepend", init, showRunes(input), result, expect)
	}
}

func TestScanRight(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		init   int
		input  []int
		expect []int
	}
	testcases := []TestCase{
		{add, 0, []int{1, 2, 3, 4}, []int{10, 9, 7, 4, 0}},
		{add, 42, []int{}, []int{42}},
		{subtract, 100, []int{1, 2, 3, 4}, []int{98, -97, 99, -96, 100}},
	}
	errorMsg := "ScanRight(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		result := ScanRight(test.callb, test.init, test.input)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf(errorMsg, funcName(test.callb), test.init, test.input, result, test.expect)
		}
	}
	prepend := func(r rune, s string) string {
		return string(r) + s
	}
	expect := []string{"abcdfoo", "bcdfoo", "cdfoo", "dfoo", "foo"}
	init := "foo"
	input := []rune{'a', 'b', 'c', 'd'}
	result := ScanRight(prepend, init, input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(errorMsg, "prepend", init, showRunes(input), result, expect)
	}
}

func TestConcatMap(t *testing.T) {
	testcases := []map[string][]int{
		{"input": {}, "expect": {}},
		{"input": {1, 2, 3}, "expect": {-1, 1, -2, 2, -3, 3}},
	}
	fn := func(i int) []int {
		return []int{-i, i}
	}
	for _, tc := range testcases {
		result := ConcatMap(fn, tc["input"])
		if !reflect.DeepEqual(result, tc["expect"]) {
			t.Errorf("ConcatMap(fn, %v) = %v, expected %v", tc["input"], result, tc["expect"])
		}
	}
}

func BenchmarkConcatMap(b *testing.B) {
	fn := func(i int) []int {
		return []int{-i, i}
	}
	for i := 0; i < b.N; i++ {
		ConcatMap(fn, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
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

func greaterThan(y int) func(int) bool {
	return func(x int) bool {
		return x > y
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

func isOneOf[T comparable](xs []T) func(T) bool {
	return func(e T) bool {
		for _, x := range xs {
			if x == e {
				return true
			}
		}
		return false
	}
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
