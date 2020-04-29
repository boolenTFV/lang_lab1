package analyze

import (
	"github.com/deckarep/golang-set"
	"fmt"
)

func ParseString(input string, Q mapset.Set, sigma map[string]map[byte]string, q0 string, F mapset.Set, E mapset.Set) {
	qState := q0
	for pos, char := range input {
		if sigma[qState][byte(char)] != "" {
			fmt.Print(qState + "-" + string(char) + "->" + sigma[qState][byte(char)])
			qState = sigma[qState][byte(char)];
		} else {
			fmt.Print("Ошибка в при разборе строки в позиции " + string(pos) + ";")
			return;
		}
	}

	if !F.Contains(qState) {
		fmt.Print("Ошибка в при разборе строки в позиции " + string(len(input)) + "; Разбор закончился на не конечном состоянии;")
	}
}