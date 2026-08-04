// Package main is a minimal reproduction case for the buildOutputs()
// non-determinism bug.
//
// The bug: buildOutputs() only checked if the direct failure target was an
// end node, without following the failure chain. When keywords were inserted
// in different orders (due to Go map iteration), some failure targets were
// non-end intermediate nodes, causing output links to be lost.
//
// This test uses a small keyword set with overlapping patterns where the
// failure chain goes through non-end intermediate nodes. It verifies that
// matching results are deterministic regardless of insertion order.
//
// Run with: go run .
package main

import (
	"fmt"
	"os"
	"sort"

	cedar "github.com/eugene-fedorenko/ahocorasick"
)

func main() {
	// Keywords with overlapping patterns.
	// "abc" and "bcd" share the suffix/prefix "bc" (which is NOT a keyword).
	// "d" is a suffix of both "bcd" and (via failure chain) "abc".
	// The failure chain for "abc" goes through the non-end node "bc" (prefix
	// of "bcd") before reaching "d". Before the fix, buildOutputs would not
	// propagate the output link from "bc" to "abc" because "bc" is not an
	// end node.
	keywords := []string{
		"abc",
		"bcd",
		"d",
	}

	// Build expected result with sorted (deterministic) insertion order.
	expected := matchWithSortedKeys("abcd", keywords)
	fmt.Println("keywords:", keywords)
	fmt.Println("text: abcd")
	fmt.Println("expected matches:", expected)

	// Run multiple rounds with random (map-based) insertion order.
	// Before the fix, some rounds would miss "d" or "bcd".
	keysMap := make(map[string]int, len(keywords))
	for i, k := range keywords {
		keysMap[k] = i
	}

	failures := 0
	const rounds = 50
	for i := 0; i < rounds; i++ {
		m := cedar.NewMatcher()
		for k, v := range keysMap {
			m.Insert([]byte(k), v)
		}
		m.Compile()

		result := match(m, []byte("abcd"))
		if !equal(result, expected) {
			failures++
			fmt.Printf("round %d: MISMATCH!\n  expected: %v\n  got:      %v\n", i, expected, result)
		}
	}

	if failures > 0 {
		fmt.Printf("\nFAIL: %d/%d rounds produced different results (non-deterministic!)\n", failures, rounds)
		os.Exit(1)
	}
	fmt.Printf("\nPASS: all %d rounds produced identical results\n", rounds)
}

func matchWithSortedKeys(text string, keywords []string) []string {
	m := cedar.NewMatcher()
	for i, k := range keywords {
		m.Insert([]byte(k), i)
	}
	m.Compile()
	return match(m, []byte(text))
}

func match(m *cedar.Matcher, data []byte) []string {
	var result []string
	resp := m.Match(data)
	for resp.HasNext() {
		for _, itr := range resp.NextMatchItem(data) {
			result = append(result, string(m.Key(data, itr)))
		}
	}
	resp.Release()
	return unique(result)
}

func unique(s []string) []string {
	seen := make(map[string]bool, len(s))
	var result []string
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	sort.Strings(result)
	return result
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
