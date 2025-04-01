## fuzzymatch
A simple Go package for fuzzy string matching using the Levenshtein distance
algorithm. Find the closest match to a string from a list of possibilities,
perfect for handling typos or approximate searches.

### Installation

To use fuzzymatch in your Go project, run:

    go get github.com/MarkusZoppelt/fuzzymatch

### Usage

Here's a quick example of how to use fuzzymatch:

```go
package main

import (
	"fmt"
	"github.com/MarkusZoppelt/fuzzymatch"
)

func main() {
	possibleMatches := []string{"foo", "bar", "baz", "foobar"}
	input := "fo"
	closest := fuzzymatch.SuggestClosestMatch(input, possibleMatches, 9999)
	if closest != "" {
		fmt.Printf("Did you mean %q?\n", closest) // Output: Did you mean "foo"?
	} else {
		fmt.Println("No close match found or exact match detected.")
	}
}
```

### Key Functions
* `SuggestClosestMatch(s string, possibleMatches []string, minDist int)
  string`: Returns the closest match from possibleMatches to s if within a
  Levenshtein distance of 3, or an empty string if no close match is found or s
  is an exact match.
* `Levenshtein(a, b string) int`: Computes the Levenshtein distance between two
  strings.

### Features

* Lightweight and dependency-free.
* Handles exact matches, close matches (up to 3 edits), and no-match cases.
* Ideal for configuration key validation, search suggestions, or typo correction.
