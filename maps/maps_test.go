package maps

import (
	"reflect"
	"sort"
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
