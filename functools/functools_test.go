package functools

import (
	"reflect"
	"sort"
	"strings"
	"testing"
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
		t.Errorf("UnZipWith(split, %v) = %v, %v, expected %v, %v", datapoints, result1, result2, expect1, expect2)
	}
}

func TestPipe(t *testing.T) {
	input := "  Lorem ipsum dolor sit amet, consectetur  "
	expect := "lorem-ipsum-dolor-sit-amet-consectetur"
	rep := strings.NewReplacer(",", "", ".", "", " ", "-")
	slugify := Pipe(strings.TrimSpace, rep.Replace, strings.ToLower)
	result := slugify(input)
	if result != expect {
		t.Errorf("Pipe(s.TrimSpace, r.Replace, s.ToLower)(%q) = %q, expected %q", input, result, expect)
	}
}

func TestCompose(t *testing.T) {
	input := "Lorem ipsum dolor sit amet, consectetur..."
	expect := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	}
	rep := strings.NewReplacer(",", "", ".", "")
	split := func(sep string) func(string) []string {
		return func(s string) []string {
			return strings.Split(s, sep)
		}
	}
	tokenize := Compose(split(" "), Compose(strings.ToLower, rep.Replace))
	result := tokenize(input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Compose(split(" "), Compose(s.ToLower, r.Replace))(%q) = %#v, expected %#v`, input, result, expect)
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
		t.Errorf(`Curry2(Flip(s.Split))(" ")(%q) = %#v, expected %#v`, input, result, expect)
	}
}

func TestCurry3(t *testing.T) {
	input := "Lorem ipsum, dolor sit amet, consectetur."
	expect := []string{"Lorem ipsum", "dolor sit amet, consectetur."}
	splitN := Curry3(strings.SplitN)(input)(", ")
	result := splitN(2) // at most 2 substrings; the last substring is the unsplit remainder
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Curry3(s.SplitN)(%q)(",")(2) = %#v, expected %#v`, input, result, expect)
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
		t.Errorf(`Partial1(Flip(s.Split), " ")(%q) = %#v, expected %#v`, input, result, expect)
	}
}

func TestPartial2(t *testing.T) {
	input := "Lorem ipsum, dolor sit amet, consectetur."
	expect := []string{"Lorem ipsum", "dolor sit amet, consectetur."}
	splitN := Partial2(strings.SplitN, input, ", ")
	result := splitN(2) // at most 2 substrings; the last substring is the unsplit remainder
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Partial2(s.SplitN, %q, ",")(2) = %#v, expected %#v`, input, result, expect)
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

func TestMapToSlice(t *testing.T) {
	type Tuple2 struct {
		fst string
		snd int
	}
	asTuple2 := func(fst string, snd int) Tuple2 {
		return Tuple2{fst, snd}
	}
	items := map[string]int{
		"lorem": 1, "ipsum": 2, "dolor": 3,
	}
	expect := []Tuple2{
		{"lorem", 1}, {"ipsum", 2}, {"dolor", 3},
	}
	result := MapToSlice(items, asTuple2)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].snd < result[j].snd
	})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("MapToSlice(%v, asTuple2) = %v, expected %v", items, result, expect)
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
