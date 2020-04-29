package tree

import (
    "testing"
)

type TreeTestpair struct {
    input string
    expected string
}

var ParseRegTestSet = []TreeTestpair{
    {"(a|b)*abb", "ab|*a&b&b&"},
    {"ab|cd", "ab&cd&|"},
    {"(ab)cd(ef)", "ab&c&d&ef&&"},
}

func TestParseRegex(t *testing.T) {
	for _, pair := range ParseRegTestSet {
        v := ParseRegex(pair.input).ToString()
        if v != pair.expected {
            t.Error(
                "For", pair.input, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}