# Golang functools

Package functools provides generic higher-order functions. They can be used to build functions of functions in a concise manner. The goal is to combine proven functional programming concepts with Go's practical programming style.

[![Go Reference](https://pkg.go.dev/badge/github.com/basbiezemans/gofunctools.svg)](https://pkg.go.dev/github.com/basbiezemans/gofunctools/functools)

## Go version

1.21 or later

## Install

```bash
go get github.com/basbiezemans/gofunctools/functools
```

## Examples

#### String sanitizer w/ Pipe
```go
replacer := strings.NewReplacer(",", "", ".", "")
sanitize := Pipe(strings.TrimSpace, replacer.Replace, strings.ToLower)

text := "  Lorem ipsum dolor sit amet, consectetur.  "

fmt.Println(sanitize(text))

// Output: "lorem ipsum dolor sit amet consectetur"
```
#### String tokenizer w/ Curry2, Flip
```go
split := Curry2(Flip(strings.Split))
words := split(" ")

text := "lorem ipsum dolor sit amet consectetur"

fmt.Println(words(text))

// Output: ["lorem", "ipsum", "dolor", "sit", "amet", "consectetur"]
```
#### Sanitizer-tokenizer w/ Partial1, Flip, Compose
```go
replacer := strings.NewReplacer(",", "", ".", "")
words := Partial1(Flip(strings.Split), " ")
tokenize := Compose(words, Compose(strings.ToLower, replacer.Replace))

text := "Lorem ipsum dolor sit amet, ...consectetur."

fmt.Println(tokenize(text))

// Output: ["lorem", "ipsum", "dolor", "sit", "amet", "consectetur"]
```
#### Word frequency counter w/ Partial2, FoldLeft
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