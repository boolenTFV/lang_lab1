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
	var input string
	for true {
		fmt.Print("input regex: ")
		fmt.Scan(&input)
		// input = "(a|b)*abb"
		inputExt := input + "#"
		fmt.Println("extended regex: ", inputExt)
		lex := tree.ParseRegex(inputExt)
		flowpos := get_dfa.ComputeMetricsPos(lex)
		fmt.Println("tree: ")
		lex.DrawTree(0)
		fmt.Println("\nflowpos:")
		get_dfa.FlowposPrint(flowpos)
		q0, endState, sigma, states := get_dfa.GetDfa(inputExt, lex, flowpos);
		fmt.Println("\ndfa:")
		Q, F, E := get_minimized.GetMapSets(states, endState.ToString(), q0.ToString(), input)
		get_minimized.DrawDFA(Q, sigma, q0.ToString(), F, E)
		P := get_minimized.Minimize(Q, F, E, sigma)
		nQ, nSigma, nq0, nF := get_minimized.NewDfa(sigma, q0.ToString(), F, P, E)
		fmt.Println("\ndfa:")
		get_minimized.DrawDFA(nQ, nSigma, nq0, nF, E)
		fmt.Println("\nВведите строку для разбора:")
		fmt.Scan(&input)
		fmt.Println()
		analyze.ParseString(input, nQ, nSigma, nq0, nF, E)
		fmt.Println()
	}
}