package functools

import (
	"reflect"
	"strings"
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

func TestReduceLeft(t *testing.T) {
	subtract := func(x, y int) int {
		return x - y
	}
	numbers := []int{1, 2, 3, 4}
	expect := -8
	result := ReduceLeft(subtract, numbers)
	if result != expect {
		t.Errorf("ReduceLeft(subtract, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestReduceRight(t *testing.T) {
	subtract := func(x, y int) int {
		return x - y
	}
	numbers := []int{1, 2, 3, 4}
	expect := -2
	result := ReduceRight(subtract, numbers)
	if result != expect {
		t.Errorf("ReduceRight(subtract, %v) = %d, expected %d", numbers, result, expect)
	}
}

func TestFoldLeft(t *testing.T) {
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

func TestPipe(t *testing.T) {
	input := "  Lorem ipsum dolor sit amet, consectetur  "
	expect := "lorem-ipsum-dolor-sit-amet-consectetur"
	rep := strings.NewReplacer(",", "", ".", "", " ", "-")
	slugify := Pipe(strings.TrimSpace, rep.Replace, strings.ToLower)
	result := slugify(input)
	if result != expect {
		t.Errorf(`Pipe(s.TrimSpace, r.Replace, s.ToLower)("%s") = %s, expected %s`, input, result, expect)
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
		t.Errorf(`Compose(split(" "), Compose(s.ToLower, r.Replace))("%s") = %v, expected %v`, input, result, expect)
	}
}

func TestFlipCurry2(t *testing.T) {
	input := "lorem ipsum dolor sit amet consectetur"
	expect := []string{
		"lorem", "ipsum", "dolor", "sit", "amet", "consectetur",
	}
	split := Curry2(Flip(strings.Split))
	tokenize := split(" ")
	result := tokenize(input)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(`Curry2(Flip(s.Split))(" ")("%s") = %v, expected %v`, input, result, expect)
	}
}
