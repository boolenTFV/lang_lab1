package main

import (
	"fmt"
	"lab1/tree"
	"lab1/get_dfa"
	"lab1/get_minimized"
	"lab1/analyze"
)

func main() {
	fmt.Println("start")
	var in []string = []string{"(a|b)*abb", "b*|ba", "abb*(c|de)", "(ab)cd*" }
	test := [][]string{[]string{"abb", "abababb", "abbb"}, []string{"bbbb", "ba", "abba"}, []string{"abc", "abde", "abce"}, []string{"abc", "abcdd", "ab"}}
	for i := 0; i<4; i++ {
		input := in[i]
		inputExt := input + "#"
		fmt.Println("\n -- Регулярное выражение: ", input)
		// fmt.Println("extended regex: ", inputExt)
		lex := tree.ParseRegex(inputExt)
		flowpos := get_dfa.ComputeMetricsPos(lex)
		// fmt.Println("tree: ")
		// lex.DrawTree(0)
		// fmt.Println("\nflowpos:")
		// get_dfa.FlowposPrint(flowpos)
		q0, F, sigma, states := get_dfa.GetDfa(inputExt, lex, flowpos);
		// fmt.Println("\ndfa:")
		Q, _, E := get_minimized.GetMapSets(states, "", q0.ToString(), input)
		// get_minimized.DrawDFA(Q, sigma, q0.ToString(), F, E)
		P := get_minimized.Minimize(Q, F, E, sigma)
		nQ, nSigma, nq0, nF := get_minimized.NewDfa(sigma, q0.ToString(), F, P, E)
		// fmt.Println("\ndfa:")
		// get_minimized.DrawDFA(nQ, nSigma, nq0, nF, E)
		// fmt.Println()
		for j := 0; j<3; j++ {
			fmt.Println("Строка для разбора:", test[i][j])
			analyze.ParseString(test[i][j], nQ, nSigma, nq0, nF, E)
			fmt.Println("\n")
		}
		fmt.Println()
	}
}