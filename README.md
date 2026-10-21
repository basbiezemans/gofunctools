# Go functools

Functional programming utilities, the Go way.

Go isn't a functional language, but since Go 1.18 its generics and first-class functions make functional patterns both possible and practical. `gofunctools` fills the gap: it provides generic higher-order functions that let you compose functions in a concise way, without importing a framework or fighting the language.

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

## Package overview
The library is deliberately idiomatic Go. It operates on slices — the language's primary data structure — and uses pure higher-order functions without Result types, monads, or other abstractions foreign to Go. You read and write it the way you'd write any other Go code.

### Reduction and aggregation
Functions such as FoldLeft, FoldRight, ScanLeft, ScanRight, and Reduce combine slice elements into a single value or sequence of intermediate values. These are useful when computing totals, accumulations, running balances, or any summary where the order of evaluation matters.

Example 1: Running totals with ScanLeft
```go
changes := []int{25, -20, 30}
balances := ScanLeft(
    func(balance, change int) int { return balance + change },
    100,
    changes,
)

fmt.Println(balances)

// Output: [100 125 105 135]
```
Example 2: Word frequencies with FoldLeft
```go
fruit := "mango banana apple pear banana grapes pear kiwi apple"
words := strings.Split(fruit, " ")

type frequency map[string]int

func count(freq frequency, word string) frequency {
    freq[word] += 1
    return freq
}

wordfreq := FoldLeft(count, frequency{}, words)

fmt.Println(wordfreq)

// Output: map[apple:2 banana:2 grapes:1 kiwi:1 mango:1 pear:2]
```
### Transformations and filtering
Map, Filter, FilterMap, ConcatMap, Partition, UnzipWith, and ZipWith make common slice transformations more concise.

Example: Filtering and transforming numbers with Filter, Map
```go
numbers := []int{1, 2, 3, 4, 5, 6}

even := func(n int) bool { return n%2 == 0 }
evens := Filter(even, numbers)

square := func(n int) int { return n * n }
squares := Map(square, evens)

fmt.Println(squares)

// Output: [4 16 36]
```
### Composition and currying
Pipe, Compose, Curry1, Curry2, Curry3, Partial1, Partial2, Flip, and Rotate help build reusable function pipelines in a style that fits Go's function values. These helpers let you compose functions and partially apply arguments to build reusable data-processing pipelines.

Example 1: Data transformation pipeline with Pipe
```go
text := "... LOREM IPSUM, dolor sit amet, consectetur ...  "

replacer := strings.NewReplacer(",", "", ".", "")

normalize := Pipe(replacer.Replace, strings.TrimSpace, strings.ToLower)

fmt.Println(normalize(text))

// Output: lorem ipsum dolor sit amet consectetur
```
Example 2: Building a complex function with Compose, Partial1, Partial2
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
### Search, grouping, and discovery
Find, FindIndex, FindIndices, GroupBy, Count, Any, and All locate values or classify elements according to predicates. These functions are especially helpful when working with data processing.

Example: Consecutive runs of equal values with GroupBy
```go
numbers := []int{1, 1, 3, 3, 3, 1, 4, 4, 4, 5}

equal := func(a, b int) bool { return a == b }

groups := GroupBy(equal, numbers)

fmt.Println(groups)

// Output: [[1 1] [3 3 3] [1] [4 4 4] [5]]
```