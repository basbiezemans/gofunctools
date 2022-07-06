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
