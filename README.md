# Go functools

Package `gofunctools` provides generic higher-order functions. They can be used to build functions of functions in a concise manner.

[![Go Reference](https://pkg.go.dev/badge/github.com/basbiezemans/gofunctools.svg)](https://pkg.go.dev/github.com/basbiezemans/gofunctools)
[![Test](https://github.com/basbiezemans/gofunctools/actions/workflows/test.yml/badge.svg)](https://github.com/basbiezemans/gofunctools/actions/workflows/test.yml)

## Minimum Go version

This package requires version 1.23 or later.

## How to use this library

Use `go get` to add this library to your project's go.mod file.

```bash
go get github.com/basbiezemans/gofunctools
```
Then import the library in your Go code with:

```go
// Without a package prefix
import . github.com/basbiezemans/gofunctools

// With a short import alias
import fn github.com/basbiezemans/gofunctools
```

## Examples

#### String sanitizer with Pipe
```go
replacer := strings.NewReplacer(",", "", ".", "")
sanitize := Pipe(strings.TrimSpace, replacer.Replace, strings.ToLower)

text := "  Lorem Ipsum dolor sit amet, consectetur.  "

fmt.Println(sanitize(text))

// Output: "lorem ipsum dolor sit amet consectetur"
```
#### Sanitizer-tokenizer with Compose
```go
replacer := strings.NewReplacer(",", "", ".", "")
sanitize := Compose(strings.ToLower, replacer.Replace)
tokenize := Compose(strings.Fields, sanitize)

text := "  Lorem Ipsum dolor sit amet, consectetur.  "

fmt.Println(tokenize(text))

// Output: ["lorem", "ipsum", "dolor", "sit", "amet", "consectetur"]
```
#### Word frequency counter with Partial2 and FoldLeft
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
#### Sum of squared even numbers with Compose, Partial1, Partial2, Map, Filter, and FoldLeft
This pipeline filters the input, squares each remaining number, then folds the
squares into a total.
```go
numbers := []int{1, 2, 3, 4, 5, 6}

isEven := func(n int) bool { return n%2 == 0 }
square := func(n int) int { return n * n }
add := func(total, n int) int { return total + n }

sumSquaresOfEven := Compose(
    Partial2(FoldLeft, add, 0),
    Compose(
        Partial1(Map, square),
        Partial1(Filter, isEven),
    ),
)

fmt.Println(sumSquaresOfEven(numbers))

// Output: 56
```
#### Maximum consecutive group with GroupBy, Compose, Partial1, Map, and Filter
This is a more elaborate example that requires two small helper functions to
carry out the task. We want to find a segment containing the maximum number of
consecutive identical elements from a sequence of numbers. In this example, we
will choose the first result.
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
