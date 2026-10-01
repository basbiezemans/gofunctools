# Go functools

Package `gofunctools` provides generic higher-order functions. They can be used to build functions of functions in a concise manner.

[![Go Reference](https://pkg.go.dev/badge/github.com/basbiezemans/gofunctools.svg)](https://pkg.go.dev/github.com/basbiezemans/gofunctools)

## Go version

This package requires version 1.23 or later.

## Install

```bash
go get github.com/basbiezemans/gofunctools
```

## Examples

#### String sanitizer with Pipe
```go
replacer := strings.NewReplacer(",", "", ".", "")
sanitize := Pipe(strings.TrimSpace, replacer.Replace, strings.ToLower)

text := "  Lorem ipsum dolor sit amet, consectetur.  "

fmt.Println(sanitize(text))

// Output: "lorem ipsum dolor sit amet consectetur"
```
#### String tokenizer with Curry2, Flip
```go
split := Curry2(Flip(strings.Split))
words := split(" ")

text := "lorem ipsum dolor sit amet consectetur"

fmt.Println(words(text))

// Output: ["lorem", "ipsum", "dolor", "sit", "amet", "consectetur"]
```
#### Sanitizer-tokenizer with Compose, Partial1, Flip
```go
replacer := strings.NewReplacer(",", "", ".", "")
tokenize := Compose(
    Partial1(Flip(strings.Split), " "),
    Compose(strings.ToLower, replacer.Replace),
)

text := "Lorem ipsum dolor sit amet, ...consectetur."

fmt.Println(tokenize(text))

// Output: ["lorem", "ipsum", "dolor", "sit", "amet", "consectetur"]
```
#### Word frequency counter with Partial2, FoldLeft
```go
type frequency map[string]int

func count(freq frequency, word string) frequency {
    freq[word] += 1
    return freq
}

wordfreq := Partial2(FoldLeft, count, frequency{})

fruit := "mango banana apple pear banana grapes pear kiwi apple"
words := strings.Split(fruit, " ")

fmt.Println(wordfreq(words))

// Output: map["apple":2, "banana":2, "grapes":1, "kiwi":1, "mango":1, "pear":2]
```
#### Maximum consecutive group with GroupBy, Compose, Partial1, Map, and Filter
This is a more elaborate example that requires two small helper functions to
carry out the task. We want to find a segment containing the maximum number of
consecutive identical elements from a sequence of numbers. In this example, we
choose the first result.
```go
numbers := []int{1, 1, 3, 3, 3, 1, 4, 4, 4, 5}

func equal(a, b int) bool {
    return a == b
}

func length(s []int) int {
	return len(s)
}

// A slice of groups with consecutive equal elements
groups := GroupBy(equal, numbers)

// Function maxLen returns the maximum group length
maxLen := Compose(slices.Max, Partial1(Map, length))

// Function eqMaxLen determines if len(group) == maxlen
eqMaxLen := Compose(Partial1(equal, maxLen(groups)), length)

// Filter the groups and pick the first result
firstMaxConsecutiveGroup := Filter(eqMaxLen, groups)[0]

fmt.Println(firstMaxConsecutiveGroup)
// Output: [3 3 3]
```
