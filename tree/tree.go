package tree
import (
	"fmt"
	"lab1/get_token"
	"container/list"
)

type NodeValue struct {
	nodeType int
	value byte
}

type Tree struct {
	Trees list.List
	FirstPos list.List
	LastPos list.List
	Pos int
	Parrent *Tree
	NodeType int
	Value byte
}

func New() *Tree{
	return new(Tree)
}

func (tree *Tree) AddItem(nodeType int, value byte, pos int) *Tree{
	newTree := new(Tree)
	newTree.NodeType = nodeType;
	newTree.Value = value;
	newTree.Pos = pos;
	newTree.Parrent = tree;
	tree.Trees.PushBack(newTree)
	return newTree
}

func (tree *Tree) AddTree(newTree *Tree) *Tree{
	newTree.Parrent = tree;
	tree.Trees.PushBack(newTree)
	return newTree
}

func copyTree(tree *Tree) *Tree {
	nodeCopy := &Tree{NodeType: tree.NodeType, Value: tree.Value, Pos: tree.Pos}
	for temp := tree.Trees.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*Tree); ok {
			subNodeCopy := copyTree(val)
			nodeCopy.Trees.PushBack(subNodeCopy)
		}
	}
	return nodeCopy
}

func (tree *Tree) DrawTree(level int) {
	prefix := "\n+"
	for i := 0 ; i < level ; i++ {
		prefix += "---"
	}
	fmt.Printf(prefix + "> %c, type: %v", tree.Value, tree.NodeType )
	fmt.Print(" FirstPos:")
	for temp := tree.FirstPos.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(int); ok {
			fmt.Printf(" %v", val)
		}
	}
	fmt.Print(" LastPos:")
	for temp := tree.LastPos.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(int); ok {
			fmt.Printf(" %v", val)
		}
	}
	for temp := tree.Trees.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*Tree); ok {
			val.DrawTree(level+1)
		}
    }
}

func (tree *Tree) ToString() string{
	var str string
	for temp := tree.Trees.Front(); temp != nil; temp = temp.Next() {
		if val, ok := temp.Value.(*Tree); ok {
			str += val.ToString()
		}
	}
	return str + string(tree.Value)
}

type GetToken func(bool)(int, byte, int)

func Cat(getToken GetToken, stack *list.List) *Tree{
	token, sym, pos := getToken(false)
	var node *Tree = nil
	// catNode := Tree{NodeType: get_token.CAT, Value: '&'}
	if token == get_token.Q {
		getToken(true)
		node = &Tree{NodeType: get_token.Q, Value: '*'}
		node.AddTree(Cat(getToken, stack))
	} else if node != nil && token == get_token.PQ {
		getToken(true)
		qnode := &Tree{NodeType: get_token.Q, Value: '*'}
		node = &Tree{NodeType: get_token.CAT, Value: '&'}
		cat := Cat(getToken, stack)
		node.AddTree(cat)
		qnode.AddTree(copyTree(cat))
		node.AddTree(qnode)
	} else if token == get_token.SYM {
		getToken(true)
		// subNode := Cat(getToken, stack)
		symNode := &Tree{NodeType: get_token.SYM, Value: sym, Pos: pos}
		node = symNode
		/* if(subNode != nil) {
			catNode.AddTree(subNode)
			catNode.AddTree(symNode)
			node = &catNode
		} else {
			node = symNode
		} */
		
	}else if token == get_token.RP {
		getToken(true)
		stack.PushBack(get_token.RP)
		node = Block(getToken, stack)
		getToken(true)
	}
	return node;
}

func symBlock(getToken GetToken, stack *list.List) *Tree {
	token, _, _ := getToken(false)
	// fmt.Println("block:", string(sym))
	var node *Tree = nil
	if token == get_token.SYM {
		catNode := Tree{NodeType: get_token.CAT, Value: '&'}
		symNode := Cat(getToken, stack)
		subNode := symBlock(getToken, stack)
		node = symNode
		if(subNode != nil) {
			catNode.AddTree(subNode)
			catNode.AddTree(symNode)
			node = &catNode
		}
	} else if token == get_token.RP || token == get_token.Q {
		node = Cat(getToken, stack)
	}

	return node
}

func Block(getToken GetToken, stack *list.List) *Tree {
	token, _, _ := getToken(false)
	// fmt.Println("block:", string(sym))
	var node *Tree = nil
	if token == get_token.SYM {
		catNode := Tree{NodeType: get_token.CAT, Value: '&'}
		symNode := Cat(getToken, stack)
		subNode := symBlock(getToken, stack)
		node = symNode
		if(subNode != nil) {
			catNode.AddTree(subNode)
			catNode.AddTree(symNode)
			node = &catNode
		}
	} else if token == get_token.RP || token == get_token.Q {
		node = Cat(getToken, stack)
	}

	token, _, _ = getToken(false)
	if node != nil && (token == get_token.OR){
		getToken(true)
		subNode := node
		node = &Tree{NodeType: get_token.OR, Value: '|'}
		node.AddTree(Block(getToken, stack))
		node.AddTree(subNode)
		Block(getToken, stack)
	}
	return node
}
func Reg(getToken GetToken, stack *list.List) *Tree {
	node := Block(getToken, stack)
	root := node;
	for node != nil {
		token, _, _ := getToken(false)
		if token == get_token.RP {
			getToken(true)
		}
		node = Block(getToken, stack)
		if node != nil {
			subNode := root;
			root = &Tree{NodeType: get_token.CAT, Value: '&'}
			root.AddTree(node)
			root.AddTree(subNode)
		}
	}
	return root;
}

// TODO: ввести pos, которую увеличивать только для не терминалов
func ParseRegex (input string) *Tree{
	currSym := len(input)-1;
	getTokenFunc := func (isShifting bool)(int, byte, int) {
		i := currSym;
		if i < 0 {return get_token.OUT, '#', i}
		if(isShifting) {currSym--}
		return get_token.GetToken(input[i]), input[i] ,i
	}
	stak := list.New()
	tree := Reg(getTokenFunc, stak)
	return tree
}