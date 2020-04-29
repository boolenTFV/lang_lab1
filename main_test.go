package main

import (
    "testing"
    "lab1/tree"
    "lab1/get_dfa"
    "lab1/get_minimized"
    "sort"
    "strings"
)

type TreeTestpair struct {
    input string
    expected string
}

var ParseRegTestSet = []TreeTestpair{
    {"(a|b)*abb", "Set{97, 98}\nSet{1235, 123, 1234}\n123\nSet{1235}\nmap[123:map[97:1234 98:123] 1234:map[97:1234 98:1235] 1235:map[97:1234 98:123]]\n"},
}

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

func TestGetDfa(t *testing.T) {
	for _, pair := range ParseRegTestSet {
        extInput := pair.input+"#"
        root := tree.ParseRegex(pair.input)
        flowpos := get_dfa.ComputeMetricsPos(root)
		get_dfa.FlowposPrint(flowpos)
        q0, endState, sigma, states := get_dfa.GetDfa(extInput, root, flowpos);
        Q, F, E := get_minimized.GetMapSets(states, endState.ToString(), q0.ToString(), pair.input)
        v := get_minimized.DFAToString(Q, sigma, q0.ToString(), F, E)
        if SortString(v) != SortString(pair.expected) {
            t.Error(
                "For", pair.input, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}