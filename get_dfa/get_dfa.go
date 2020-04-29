package get_dfa

import (
	"fmt"
	"lab1/tree"
	"lab1/get_token"
	"container/list"
	"strconv"
	"strings"
)
func ComputeMetricsPos(node *tree.Tree)(map[int]*list.List) {
	pos := 0
	flowpos := make(map[int]*list.List)
	var computeMetrics func (*tree.Tree)(*list.List, *list.List, bool)
	computeMetrics = func (node *tree.Tree)(*list.List, *list.List , bool) {
		firstPos, lastPos := list.New(), list.New()
		isNulble := false
		if node.NodeType == get_token.SYM {
			firstPos.PushBack(pos)
			lastPos.PushBack(pos)
			pos++
		} else if node.NodeType == get_token.OR {
			for temp := node.Trees.Front(); temp != nil; temp = temp.Next() {
				if val, ok := temp.Value.(*tree.Tree); ok {
					firstPosAdd, lastPosAdd, _ := computeMetrics(val)
					firstPos.PushBackList(firstPosAdd)
					lastPos.PushBackList(lastPosAdd)
				}
			}
		} else if node.NodeType == get_token.CAT {
			var firstPosL,lastPosL, lastPosR, firstPosR *list.List
			var isNulbleL, isNulbleR bool
			if val, ok := node.Trees.Front().Value.(*tree.Tree); ok {
				firstPosL, lastPosL, isNulbleL = computeMetrics(val)

			}
			if val, ok := node.Trees.Back().Value.(*tree.Tree); ok {
				firstPosR, lastPosR, isNulbleR = computeMetrics(val)

			}
			
			if(isNulbleL) {
				firstPos.PushBackList(firstPosL)
				for temp := firstPosR.Front(); temp != nil; temp = temp.Next() {
					if val, ok := temp.Value.(int); ok && !PoSExists(val, firstPos) {
						firstPos.PushBack(val)
					}
				}
			} else {
				firstPos.PushBackList(firstPosL)
			}

			if(isNulbleR) {
				lastPos.PushBackList(lastPosL)
				for temp := lastPosR.Front(); temp != nil; temp = temp.Next() {
					if val, ok := temp.Value.(int); ok && !PoSExists(val, lastPos) {
						lastPos.PushBack(val)
					}
				}
			} else {
				lastPos.PushBackList(lastPosR)
			}


			for temp := lastPosL.Front(); temp != nil; temp = temp.Next() {
				if val, ok := temp.Value.(int); ok {
					if flowpos[val] == nil {flowpos[val] = list.New()}
					flowpos[val].PushBackList(firstPosR)
				}
			}
		} else if node.NodeType == get_token.Q {
			isNulble = true
			if val, ok := node.Trees.Front().Value.(*tree.Tree); ok {
				firstPos, lastPos, _ = computeMetrics(val)
			}
			for temp := lastPos.Front(); temp != nil; temp = temp.Next() {
				if val, ok := temp.Value.(int); ok {
					if flowpos[val] == nil {flowpos[val] = list.New()}
					flowpos[val].PushBackList(firstPos)
				}
			}
		}
		node.FirstPos.PushBackList(firstPos)
		node.LastPos.PushBackList(lastPos)
		return firstPos, lastPos, isNulble
	}
	computeMetrics(node)
	return flowpos
}

func PoSExists(pos int, positions *list.List) bool {
	for temp := positions.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(int); ok {
			if(pos == val) {
				return true
			}
		}
	}
	return false
}

type State struct {
	Content *list.List
	Marked bool
}

type SymNPose struct {
	symbol byte
	poses *list.List
}

func (state *State)ToString()string {
	stateStr := ""
	for temp := state.Content.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(int); ok {
			stateStr = stateStr + strconv.Itoa(val+1)
		}
	}
	return stateStr;
}

func getUnmarkedState(states *list.List) *State {
	for temp := states.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*State); ok {
			if(!val.Marked) {
				return val
			}
		}
	}
	return nil
}

func get_all_uniq_sym_n_poses(input string, R *State) map[byte]*list.List {
	input = strings.ReplaceAll(input, "*", "")
	input = strings.ReplaceAll(input, "|", "")
	input = strings.ReplaceAll(input, "+", "")
	input = strings.ReplaceAll(input, "(", "")
	input = strings.ReplaceAll(input, ")", "")
	result := make(map[byte]*list.List)
	for i := 0; i < len(input); i++ {
		for temp := R.Content.Front(); temp != nil; temp = temp.Next() {
			if val, ok := temp.Value.(int); ok {
				if val == i {
					sym := input[i]
					if result[sym] == nil {result[sym] = list.New()}
					result[sym].PushBack(i)
				}
			}
		}
	}
	return result;
}

func SinQ(s *State, Q *list.List) bool {
	for temp := Q.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*State); ok {
			if(val.ToString() == s.ToString()) {
				return true
			}
		}
	}
	return false
}

func GetDfa(input string, root *tree.Tree, flowpos map[int]*list.List) (State, State, map[string]map[byte]string, list.List) {
	states := list.New()
	q0 := State{Content: &root.FirstPos, Marked: false}
	endState := q0
	states.PushBack(&q0)
	rulles := make(map[string]map[byte]string)
	for R := getUnmarkedState(states); R != nil; R = getUnmarkedState(states) {
		R.Marked = true;
		for key, value := range get_all_uniq_sym_n_poses(input, R) {
			newStatePoses := list.New()
			for temp := value.Front(); temp != nil; temp = temp.Next() {
				if val, ok := temp.Value.(int); ok && flowpos[val] != nil{
					for temp := flowpos[val].Front(); temp != nil; temp = temp.Next() {
						if val, ok := temp.Value.(int); ok { 
							isExist := false
							for temp := newStatePoses.Front(); temp != nil && !isExist; temp = temp.Next() {
								if pose, ok := temp.Value.(int); ok {
									if(val == pose ) {isExist = true}
								}
							}
							if !isExist {newStatePoses.PushBack(val)}
						}
					}
				}
			}
			if newStatePoses.Len() != 0 {
				newState := State{Content: newStatePoses, Marked: false}
				if !SinQ(&newState, states) {
					states.PushBack(&newState)
					endState = newState
				}
				if(rulles[R.ToString()] == nil) {rulles[R.ToString()] = make(map[byte]string)}
				rulles[R.ToString()][key]=newState.ToString();
			}
		}
	}
	return q0, endState, rulles, *states
}

func FlowposPrint(flowpos map[int]*list.List) {
	for k, v := range flowpos {
		fmt.Printf("%v: ", k+1)
		
		for temp := v.Front(); temp != nil; temp = temp.Next() {
			if val, ok := temp.Value.(int); ok {
				fmt.Print(val+1, "; ")
			}
		}
		fmt.Println()
	}
}
 