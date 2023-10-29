# Go functools

Package `funcs` provides generic higher-order functions. They can be used to build functions of functions in a concise manner.

[![Go Reference](https://pkg.go.dev/badge/github.com/basbiezemans/gofunctools.svg)](https://pkg.go.dev/github.com/basbiezemans/gofunctools)

## Go version

This package requires version 1.21 or later.

## Install

```bash
go get github.com/basbiezemans/gofunctools/funcs
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
sanitize := Compose(strings.ToLower, replacer.Replace)
splitstr := Partial1(Flip(strings.Split), " ")
tokenize := Compose(splitstr, sanitize)

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

fmt.Println(wordfreq(words)))

// Output: map["apple":2, "banana":2, "grapes":1, "kiwi":1, "mango":1, "pear":2]
```