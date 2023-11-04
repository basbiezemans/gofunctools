package maps

import (
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestHashMapToSlice(t *testing.T) {
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
	result := HashMapToSlice(asTuple2, items)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].snd < result[j].snd
	})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("MapToSlice(%v, asTuple2) = %v, expected %v", items, result, expect)
	}
}

func TestMap(t *testing.T) {
	mapper := func(key string, value float64) (string, float64) {
		return strings.ToUpper(key), math.Exp2(value)
	}
	foomap := map[string]float64{
		"foo": 1, "bar": 2, "baz": 3,
	}
	expect := map[string]float64{
		"FOO": 2, "BAR": 4, "BAZ": 8,
	}
	result := Map(mapper, foomap)
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("Map(mapper, %v) = %v, expected %v", foomap, result, expect)
	}
}
