package pair

import (
	"reflect"
	"testing"

	fts "github.com/basbiezemans/gofunctools"
	opr "github.com/basbiezemans/gofunctools/operators"
)

func TestSwap(t *testing.T) {
	pair := New(1, 2)
	want := New(2, 1)
	have := Swap(pair)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Swap(%v) = %v, expected %v", pair, have, want)
	}
}

func TestFirst(t *testing.T) {
	pair := New(1, 2)
	want := New(2, 2)
	add1 := fts.Partial1(opr.Add, 1)
	have := First(add1, pair)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("First(add1, %v) = %v, expected %v", pair, have, want)
	}
}

func TestSecond(t *testing.T) {
	pair := New(1, 2)
	want := New(1, 4)
	mul2 := fts.Partial1(opr.Multiply, 2)
	have := Second(mul2, pair)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Second(mul2, %v) = %v, expected %v", pair, have, want)
	}
}

func TestBoth(t *testing.T) {
	pair := New(1, 2)
	want := New(2, 4)
	mul2 := fts.Partial1(opr.Multiply, 2)
	have := Both(mul2, pair)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Both(mul2, %v) = %v, expected %v", pair, have, want)
	}
}

func TestBimap(t *testing.T) {
	pair := New(1, 2)
	want := New(2, 4)
	add1 := fts.Partial1(opr.Add, 1)
	mul2 := fts.Partial1(opr.Multiply, 2)
	have := Bimap(add1, mul2, pair)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Bimap(add1, mul2, %v) = %v, expected %v", pair, have, want)
	}
}

func TestFanout(t *testing.T) {
	ival := 2
	want := New(3, 4)
	add1 := fts.Partial1(opr.Add, 1)
	mul2 := fts.Partial1(opr.Multiply, 2)
	have := Fanout(add1, mul2, ival)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("Fanout(add1, mul2, %d) = %v, expected %v", ival, have, want)
	}
}
