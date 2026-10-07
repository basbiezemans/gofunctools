package functools

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
	"unicode"
)

func TestAny(t *testing.T) {
	type TestCase struct {
		input  []int
		expect bool
	}
	testcases := []TestCase{
		{nil, false},
		{[]int{}, false},
		{[]int{1}, false},
		{[]int{2}, true},
		{[]int{2, 4}, true},
		{[]int{1, 2, 1}, true},
		{[]int{1, 3, 4, 7}, true},
		{[]int{1, 3, 5, 7, 9}, false},
	}
	errorMsg := "Any(even, %v) = %t, expected %t"
	for _, test := range testcases {
		result := Any(even, test.input)
		if result != test.expect {
			t.Errorf(errorMsg, test.input, result, test.expect)
		}
	}
}

func TestAll(t *testing.T) {
	type TestCase struct {
		input  []int
		expect bool
	}
	testcases := []TestCase{
		{nil, true},
		{[]int{}, true},
		{[]int{1}, false},
		{[]int{2}, true},
		{[]int{2, 4}, true},
		{[]int{1, 2, 1}, false},
		{[]int{2, 4, 6, 8}, true},
		{[]int{2, 4, 7, 8, 2}, false},
	}
	errorMsg := "All(even, %v) = %t, expected %t"
	for _, test := range testcases {
		result := All(even, test.input)
		if result != test.expect {
			t.Errorf(errorMsg, test.input, result, test.expect)
		}
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

func TestReduceLeftEmptyPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("ReduceLeft did not panic")
		}
	}()
	ReduceLeft(subtract, []int{})
}

func TestReduceRight(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	expect := -2
	result := ReduceRight(subtract, numbers)
	if result != expect {
		t.Errorf("ReduceRight(subtract, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestReduceRightEmptyPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("ReduceRight did not panic")
		}
	}()
	ReduceRight(subtract, []int{})
}

func TestReduce(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		input  []int
		expect int
	}
	testcases := []TestCase{
		{add, []int{1}, 1},
		{add, []int{1, 2}, 3},
		{multiply, []int{2, 3, 4}, 24},
		{add, makeRange(0, 10_001), 50_005_000},
	}
	for _, tc := range testcases {
		result := Reduce(tc.callb, tc.input)
		if result != tc.expect {
			fname := funcName(tc.callb)
			xs := ""
			if len(tc.input) > 5 {
				xs = "[...]"
			} else {
				xs = fmt.Sprintf("%v", tc.input)
			}
			t.Errorf("Reduce(%s, %s) = %d, expected %d", fname, xs, result, tc.expect)
		}
	}
}

func TestReduceEmptyPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Reduce did not panic")
		}
	}()
	Reduce(add, []int{})
}

func TestFoldLeft(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) int
		init   int
		input  []int
		expect int
	}
	intdiv := func(acc, v int) int { return acc / v }
	testcases := []TestCase{
		{add, 0, []int{}, 0},
		{add, 0, []int{1}, 1},
		{add, 0, []int{1, 2, 3, 4}, 10},
		{multiply, 1, []int{1, 2, 3, 4}, 24},
		{intdiv, 100, []int{2, 2, 5}, 5},
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
	intdiv := func(v, acc int) int { return acc / v }
	testcases := []TestCase{
		{add, 0, []int{}, 0},
		{add, 0, []int{1}, 1},
		{add, 0, []int{1, 2, 3, 4}, 10},
		{multiply, 1, []int{1, 2, 3, 4}, 24},
		{intdiv, 100, []int{2, 2, 5}, 5},
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

func TestParallelMap(t *testing.T) {
	type TestCase struct {
		input  []int
		expect []int
	}
	var nilSlice []int = nil
	var bigSlice = makeRange(0, 10_001)
	testcases := []TestCase{
		{nilSlice, []int{}},
		{[]int{}, []int{}},
		{[]int{1}, []int{2}},
		{[]int{0, 1}, []int{0, 2}},
		{[]int{1, 2, 3, 4}, []int{2, 4, 6, 8}},
		{bigSlice, Map(double, bigSlice)},
	}
	for _, tc := range testcases {
		result := ParallelMap(double, tc.input)
		if !reflect.DeepEqual(result, tc.expect) {
			s1 := sliceToString(tc.input, 5)
			s2 := sliceToString(result, 5)
			s3 := sliceToString(tc.expect, 5)
			t.Errorf("ParallelMap(double, %s) = %s, expected %s", s1, s2, s3)
		}
	}
}

func TestParallelMapPanic(t *testing.T) {
	expect := []int{2, 4, 0, 8}
	result := ParallelMap(func(x int) int {
		if x == 3 {
			panic("test")
		}
		return double(x)
	}, []int{1, 2, 3, 4})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ParallelMap(double, []int{1,2,3,4}) = %v, expected %v", result, expect)
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
		t.Errorf("UnzipWith(split, %v) = %v, %v, expected %v, %v", datapoints, result1, result2, expect1, expect2)
	}
}

func TestPartial1PairUnzipWith(t *testing.T) {
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
	unzip := Partial1Pair(UnzipWith, split)
	result1, result2 := unzip(datapoints)
	if !reflect.DeepEqual(result1, expect1) || !reflect.DeepEqual(result2, expect2) {
		t.Errorf("Partial1Pair(UnzipWith, split)(%v) = %v, %v, expected %v, %v", datapoints, result1, result2, expect1, expect2)
	}
}

func TestPipe(t *testing.T) {
	input := "  Lorem ipsum dolor sit amet, consectetur  "
	expect := "lorem-ipsum-dolor-sit-amet-consectetur"
	replacer := strings.NewReplacer(",", "", ".", "", " ", "-")
	slugify := Pipe(strings.TrimSpace, replacer.Replace, strings.ToLower)
	result := slugify(input)
	if result != expect {
		t.Errorf("Pipe(TrimSpace, Replace, ToLower)(%q) = %q, expected %q", input, result, expect)
	}
}

func TestCompose(t *testing.T) {
	input := "Lorem ipsum dolor sit amet, consectetur..."
	expect := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	}
	replacer := strings.NewReplacer(",", "", ".", "")
	split := func(sep string) func(string) []string {
		return func(s string) []string {
			return strings.Split(s, sep)
		}
	}
	tokenize := Compose(split(" "), Compose(strings.ToLower, replacer.Replace))
	result := tokenize(input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Compose(split(" "), Compose(ToLower, Replace))(%q) = %#v, expected %#v`, input, result, expect)
	}
}

func TestFlipCurry2(t *testing.T) {
	input := "lorem ipsum dolor sit amet consectetur"
	expect := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	}
	split := Curry2(Flip(strings.Split))
	words := split(" ")
	result := words(input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Curry2(Flip(Split))(" ")(%q) = %#v, expected %#v`, input, result, expect)
	}
}

func TestCurry3(t *testing.T) {
	input := "Lorem ipsum, dolor sit amet, consectetur."
	expect := []string{"Lorem ipsum", "dolor sit amet, consectetur."}
	splitN := Curry3(strings.SplitN)(input)(", ")
	result := splitN(2) // at most 2 substrings; the last substring is the unsplit remainder
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Curry3(SplitN)(%q)(",")(2) = %#v, expected %#v`, input, result, expect)
	}
}

func TestFlipPartial1(t *testing.T) {
	input := "lorem ipsum dolor sit amet consectetur"
	expect := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	}
	words := Partial1(Flip(strings.Split), " ")
	result := words(input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Partial1(Flip(Split), " ")(%q) = %#v, expected %#v`, input, result, expect)
	}
}

func TestPartial2(t *testing.T) {
	input := "Lorem ipsum, dolor sit amet, consectetur."
	expect := []string{"Lorem ipsum", "dolor sit amet, consectetur."}
	splitN := Partial2(strings.SplitN, input, ", ")
	result := splitN(2) // at most 2 substrings; the last substring is the unsplit remainder
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Partial2(SplitN, %q, ",")(2) = %#v, expected %#v`, input, result, expect)
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

func TestPartial1PairPartition(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect1 := []int{0, 2, 4, 6, 8}
	expect2 := []int{1, 3, 5, 7, 9}
	splitEvenOdd := Partial1Pair(Partition, even)
	result1, result2 := splitEvenOdd(numbers)
	if !reflect.DeepEqual(result1, expect1) || !reflect.DeepEqual(result2, expect2) {
		t.Errorf("Partial1Pair(Partition, even)(%v) = %v, %v, expected %v, %v", numbers, result1, result2, expect1, expect2)
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

func TestHashMapToSlice(t *testing.T) {
	type Tuple2 struct {
		fst string
		snd int
	}
	type TestCase struct {
		input  map[string]int
		expect []Tuple2
	}
	newTuple2 := func(fst string, snd int) Tuple2 {
		return Tuple2{fst, snd}
	}
	testcases := []TestCase{
		{
			map[string]int{}, []Tuple2{},
		}, {
			map[string]int{"lorem": 1, "ipsum": 2, "dolor": 3},
			[]Tuple2{{"lorem", 1}, {"ipsum", 2}, {"dolor", 3}},
		},
	}
	for _, test := range testcases {
		result := HashMapToSlice(newTuple2, test.input)
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].snd < result[j].snd
		})
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("HashMapToSlice(asTuple2, %v) = %v, expected %v", test.input, result, test.expect)
		}
	}
}

func TestSliceToHashMap(t *testing.T) {
	type Tuple2 struct {
		fst string
		snd int
	}
	type TestCase struct {
		input  []Tuple2
		expect map[string]int
	}
	split := func(tuple Tuple2) (string, int) {
		return tuple.fst, tuple.snd
	}
	testcases := []TestCase{
		{
			[]Tuple2{}, map[string]int{},
		}, {
			[]Tuple2{{"a", 1}, {"b", 2}, {"a", 3}},
			map[string]int{"a": 3, "b": 2},
		}, {
			[]Tuple2{{"lorem", 1}, {"ipsum", 2}, {"dolor", 3}},
			map[string]int{"lorem": 1, "ipsum": 2, "dolor": 3},
		},
	}
	for _, test := range testcases {
		result := SliceToHashMap(split, test.input)
		if !reflect.DeepEqual(result, test.expect) {
			t.Errorf("SliceToHashMap(split, %v) = %v, expected %v", test.input, result, test.expect)
		}
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
			fname := funcName(test.callb)
			t.Errorf(errorMsg, fname, test.init, test.input, result, test.expect)
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

func TestGroupBy(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) bool
		input  []int
		expect [][]int
	}
	eq := func(x, y int) bool { return x == y }
	neq := func(x, y int) bool { return x != y }
	gt := func(x, y int) bool { return x > y }
	lt5 := func(x, y int) bool { return y-x < 5 }
	testcases := []TestCase{
		{eq, []int{}, [][]int{}},
		{eq, []int{1, 1, 1, 1}, [][]int{{1, 1, 1, 1}}},
		{eq, []int{1, 2, 2, 1}, [][]int{{1}, {2, 2}, {1}}},
		{neq, []int{1, 1, 1, 2, 3, 1, 4, 4, 5}, [][]int{{1}, {1}, {1, 2, 3}, {1, 4, 4, 5}}},
		{gt, []int{1, 3, 5, 1, 4, 2, 6, 5, 4}, [][]int{{1}, {3}, {5, 1, 4, 2}, {6, 5, 4}}},
		{lt5, makeRange(0, 20), [][]int{{0, 1, 2, 3, 4}, {5, 6, 7, 8, 9}, {10, 11, 12, 13, 14}, {15, 16, 17, 18, 19}}},
	}
	errorMsg := "GroupBy(%s, %v) = %v, expected %v"
	for _, tc := range testcases {
		result := GroupBy(tc.callb, tc.input)
		if !reflect.DeepEqual(result, tc.expect) {
			fname := funcName(tc.callb)
			t.Errorf(errorMsg, fname, tc.input, result, tc.expect)
		}
	}
}

func TestDeleteBy(t *testing.T) {
	type TestCase struct {
		callb  func(int, int) bool
		value  int
		input  []int
		expect []int
	}
	eq := func(x, y int) bool { return x == y }
	lte := func(x, y int) bool { return x <= y }
	neq := func(x, y int) bool { return x != y }
	testcases := []TestCase{
		{eq, 0, []int{}, []int{}},
		{eq, 0, []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
		{eq, 1, []int{1, 2, 3, 4}, []int{2, 3, 4}},
		{eq, 3, []int{1, 3, 3, 4}, []int{1, 3, 4}},
		{eq, 4, []int{1, 2, 3, 4}, []int{1, 2, 3}},
		{lte, 4, []int{1, 2, 3, 4, 5, 6, 7}, []int{1, 2, 3, 5, 6, 7}},
		{neq, 5, []int{5, 5, 4, 3, 5, 2}, []int{5, 5, 3, 5, 2}},
	}
	errorMsg := "DeleteBy(%s, %v, %v) = %v, expected %v"
	for _, test := range testcases {
		result := DeleteBy(test.callb, test.value, test.input)
		if !reflect.DeepEqual(result, test.expect) {
			fname := funcName(test.callb)
			t.Errorf(errorMsg, fname, test.value, test.input, result, test.expect)
		}
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

// Creates a slice containing a range of integers. The range is a half-open
// interval [x,y) where y is excluded.
func makeRange(x, y int) []int {
	if y <= x {
		return []int{}
	}
	var ys = make([]int, y-x)
	for i := range ys {
		ys[i] = i
	}
	return ys
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

func funcName(fn any) string {
	var fptr = reflect.ValueOf(fn).Pointer()
	var fname = runtime.FuncForPC(fptr).Name()
	if str, err := last(strings.Split(fname, ".")); err == nil {
		return str
	}
	return "N/A"
}

func sliceToString[T any](s []T, limit int) string {
	if len(s) <= limit {
		return fmt.Sprintf("%v", s)
	} else {
		str := fmt.Sprintf("%v", s[:limit])
		return str[:len(str)-1] + " ...]"
	}
}
