package fuzzymatch

import "slices"

// Levenshtein distance calculates the number of edits between two strings
func Levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	if a == b {
		return 0
	}

	// Initialize the distance matrix
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = minimum(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[la][lb]
}

// Helper function to find the minimum of three integers
func minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// SuggestClosestMatch finds the closest string match using Levenshtein distance
func SuggestClosestMatch(s string, possibleMatches []string, minDist int) string {
	// Check for exact match
	if slices.Contains(possibleMatches, s) {
		return ""
	}

	// Find the closest match
	var closest string
	for _, p := range possibleMatches {
		dist := Levenshtein(s, p)
		if dist < minDist {
			minDist = dist
			closest = p
		}
	}

	// Suggest only if the distance is small (e.g., <= 3)
	if minDist <= 3 {
		return closest
	}

	return ""
}
