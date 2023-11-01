package gofunctools

import (
	"reflect"
	"strings"
	"testing"
)

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
