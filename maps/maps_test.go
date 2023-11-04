package maps

import (
	"math"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/basbiezemans/gofunctools/pair"
)

func TestToSlice(t *testing.T) {
	items := map[string]int{
		"lorem": 1, "ipsum": 2, "dolor": 3,
	}
	expect := []pair.Pair[string, int]{
		pair.New("lorem", 1),
		pair.New("ipsum", 2),
		pair.New("dolor", 3),
	}
	result := ToSlice(pair.New, items)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Snd() < result[j].Snd()
	})
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("ToSlice(pair.New, %v) = %v, expected %v", items, result, expect)
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
