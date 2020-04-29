package get_minimized

import (
	"github.com/deckarep/golang-set"
	"container/list"
	"fmt"
	"lab1/get_dfa"
	"strings"
)
type ClPair struct{
	C mapset.Set
	a byte
}

type ReplaceSet struct{
	R mapset.Set
	R1 mapset.Set
	R2 mapset.Set
}


func split(R mapset.Set, C mapset.Set, a byte, dfa map[string]map[byte]string) (mapset.Set, mapset.Set) {
	R1, R2 := mapset.NewSet(), mapset.NewSet()
	for r := range R.Iterator().C {
		state := dfa[r.(string)][a]
		if(C.Contains(state)) {
			R1.Add(r)
		} else {
			R2.Add(r)
		}
	}
	return R1, R2
}

func setToS(set mapset.Set) string {
	str := ""
	for item := range set.Iterator().C {
		str = str + item.(string) + ";"
	}
	return str
}

func setBToS(set mapset.Set) string {
	str := ""
	for item := range set.Iterator().C {
		str = str + string(item.(byte)) + ";"
	}
	return str
}

func replaceRWithR1R2 (S *list.List, R ClPair, R1 ClPair, R2 ClPair ) bool {
	var findedR *list.Element = nil
	for temp := S.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(ClPair); ok && val.C.Equal(R.C) && val.a == R.a {
			findedR = temp
		}
	}

	if findedR != nil {
		S.InsertAfter(R1, findedR)
		S.InsertAfter(R2, findedR)
		S.Remove(findedR)
	}
	return findedR != nil
}

func popFront(queue *list.List) *list.Element {
	elem := queue.Front()
	queue.Remove(elem)
	return elem
}

/**
*	Q — множество состояний ДКА,
*	F — множество терминальных состояний,
*	E — множество терменалов
*	S — очередь пар ⟨C, a⟩,
*	P — разбиение множества состояний ДКА,
*	R — класс состояний ДКА.
* 	dfa - функция перехода
*/
func Minimize(Q mapset.Set, F mapset.Set, E mapset.Set, dfa map[string]map[byte]string) mapset.Set {
	P := mapset.NewSet()
	P.Add(Q)
	QdF := Q.Difference(F)
	P.Add(QdF)
	var S list.List
	for c := range E.Iterator().C {
		S.PushBack(ClPair{C: Q, a: c.(byte)})
		S.PushBack(ClPair{C: QdF, a: c.(byte)})
	}
    for S.Len() != 0 {
		var qurrentPair ClPair
		var C  mapset.Set
		var a byte
		if val, ok := popFront(&S).Value.(ClPair); ok {
			qurrentPair = val
			C = qurrentPair.C
			a = qurrentPair.a
		}
		var replaceClases list.List
		for R := range P.Iterator().C {
			R1, R2 := split(R.(mapset.Set), C, a, dfa)
			if R1.Cardinality() != 0 &&  R2.Cardinality() != 0 {
				replaceClases.PushBack(ReplaceSet{R: R.(mapset.Set), R1: R1, R2: R2} )
				for c := range E.Iterator().C {
					if !replaceRWithR1R2(&S, ClPair{C: R.(mapset.Set), a: c.(byte)}, ClPair{C: R1, a: c.(byte)}, ClPair{C: R2, a: c.(byte)}) {
						if R1.Cardinality() >=  R2.Cardinality() {
							S.PushBack(ClPair{C: R1, a: c.(byte)})
						} else {
							S.PushBack(ClPair{C: R2, a: c.(byte)})
						}
					}
				}

			}
		}
		for replaceClases.Len() != 0 {
			if val, ok := popFront(&replaceClases).Value.(ReplaceSet); ok {
				P.Remove(val.R)
				addOnce(P, val.R1)
				addOnce(P, val.R2)
			}
		}
	}
	return P
}

func PrintPClasses(P mapset.Set) {
	fmt.Println("P (" + string(P.Cardinality()) + "): ")
	for Class := range P.Iterator().C {
		fmt.Println("-- класс: ", setToS(Class.(mapset.Set)))
    }
}

func Test() mapset.Set {
	dfa := map[string]map[byte]string {
		"1": map[byte]string{'0': "2", '1': "3"},
		"2": map[byte]string{'0': "2", '1': "3"},
		"3": map[byte]string{'0': "4", '1': "5"},
		"4": map[byte]string{'0': "6", '1': "7"},
		"5": map[byte]string{'0': "2", '1': "3"},
		"6": map[byte]string{'0': "8", '1': "9"},
		"7": map[byte]string{'0': "6", '1': "7"},
		"8": map[byte]string{'0': "6", '1': "7"},
		"9": map[byte]string{'0': "10", '1': "11"},
		"10": map[byte]string{'0': "10", '1': "11"},
		"11": map[byte]string{'0': "8", '1': "9"}}
	Q := mapset.NewSet("1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11")
	F := mapset.NewSet("4", "7", "8")
	E := mapset.NewSet(byte('0'),byte('1'))
	return Minimize(Q, F, E, dfa)
}

func addOnce(P mapset.Set, R mapset.Set) {
	it := P.Iterator()
	for item := range it.C {
		if(R.Equal(item.(mapset.Set))) {
			it.Stop()
			return
		}
	}
	P.Add(R)
}

func GetMapSets(Q list.List, F string, q0 string, reg string) (mapset.Set, mapset.Set, mapset.Set) {
	reg = strings.ReplaceAll(reg, "*", "")
	reg = strings.ReplaceAll(reg, "|", "")
	reg = strings.ReplaceAll(reg, "+", "")
	reg = strings.ReplaceAll(reg, "(", "")
	reg = strings.ReplaceAll(reg, ")", "")

	QN, FN, EN := mapset.NewSet(), mapset.NewSet(), mapset.NewSet()
	FN.Add(F)
	for i := 0; i < len(reg); i++ {
        EN.Add(reg[i])
    }
	for temp := Q.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*get_dfa.State); ok{
			QN.Add(val.ToString())
		}
	}
	return QN, FN, EN
}

/**
* Создает новый ДКА из списка классов и старого DFA
*/
func NewDfa(oldSigma map[string]map[byte]string, q0 string, F mapset.Set, P mapset.Set, E mapset.Set) (mapset.Set, map[string]map[byte]string, string, mapset.Set) {
	mapOfClasses := make(map[mapset.Set]string)
	var NF mapset.Set
	Q, NF := mapset.NewSet(), mapset.NewSet()
	for class := range P.Iterator().C {
		mapOfClasses[class.(mapset.Set)] = setToS(class.(mapset.Set))
	}
	for class := range P.Iterator().C {
		Q.Add(setToS(class.(mapset.Set)))
	}
	sigma := make(map[string]map[byte]string)
	for classOut := range P.Iterator().C {
		for stateOut := range classOut.(mapset.Set).Iterator().C {
			for c := range E.Iterator().C {
				for classIn := range P.Iterator().C {
					for stateIn := range classIn.(mapset.Set).Iterator().C {
						if oldSigma[stateOut.(string)][c.(byte)] == stateIn.(string) {
							if sigma[mapOfClasses[classOut.(mapset.Set)]] == nil {
								sigma[mapOfClasses[classOut.(mapset.Set)]] = make(map[byte]string)
							}
							sigma[mapOfClasses[classOut.(mapset.Set)]][c.(byte)] = mapOfClasses[classIn.(mapset.Set)]
						}
					}
				}
			}
		}
	}

	for f := range F.Iterator().C {
		for class := range P.Iterator().C {
			for state := range class.(mapset.Set).Iterator().C {
				if f.(string) == state.(string) {
					NF.Add(mapOfClasses[class.(mapset.Set)])
				}
			}
		}
	}
	var nq0 string
	for class := range P.Iterator().C {
		for state := range class.(mapset.Set).Iterator().C {
			if q0 == state.(string) && !NF.Contains(q0) {
				nq0 = mapOfClasses[class.(mapset.Set)]
			}
		}
	}

	return Q, sigma, nq0, NF
}
func DFAToString(Q mapset.Set, sigma map[string]map[byte]string, q0 string, F mapset.Set, E mapset.Set) string {
	str := fmt.Sprintf("%v\n", E)
	str += fmt.Sprintf("%v\n", Q)
	str += fmt.Sprintf("%v\n", q0)
	fmt.Print("asafsaf", q0)
	str += fmt.Sprintf("%v\n", F)
	str += fmt.Sprintf("%v\n", sigma)
	return str
}
func DrawDFA(Q mapset.Set, sigma map[string]map[byte]string, q0 string, F mapset.Set, E mapset.Set) {
	fmt.Println("Множество терминальных символов: ", setBToS(E))
	fmt.Println("Все состояния: ", setToS(Q))
	fmt.Println("Начальное состояние: ", q0)
	fmt.Println("Конечные состояния: ", setToS(F))
	fmt.Println("Функция перехода: ")
	for key, value := range sigma {
		fmt.Println(key, " : ")
		for key, value := range value {
			fmt.Print(string(key), " -> ", value, "; ")
		}
		fmt.Println()
	}
}