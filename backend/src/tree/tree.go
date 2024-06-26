package tree

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"backend/src/regex"
	"backend/src/regex/regexFormated"

	// Import json
	"encoding/json"

	"sync"

	// Use graphivz to visualize the tree
	"github.com/awalterschulze/gographviz"
)

// Node struct to represent a node in the tree
type Node struct {
	// Define json key is Value
	Value    interface{} `json:"Value"`
	Left     *Node       `json:"Left"`
	Right    *Node       `json:"Right"`
	Ident    string      `json:"ident"` // Change the field name to be exported
	Nullable bool        `json:"Nullable"`

	First  []*Node `json:"First"`
	Last   []*Node `json:"Last"`
	Follow []*Node `json:"Follow"`
}

func NumberToLetter(num int) string {
    if num <= 26 {
        return string(rune('a' + num - 1))
    } else {
        first := string(rune('a' + (num-1)/26 - 1))
        second := string(rune('a' + (num-1)%26))
        return first + second
    }
}

func EvalAlphabet(alphabet *[]rune, toeval []interface{}) {
	// Iter over interface
	for _, v := range toeval {
		switch val := v.(type) {
			case int32:
				// Check if the value is not in the alphabet
				if !regexFormated.ContainsRune(*alphabet, val) {
					// Append the value to the alphabet with try
					*alphabet = append(*alphabet, val)
				}
			default:
				break
		}
	}
}

// Struct to represent a tree
type Tree struct {
	Root *Node
	Counter int
}


func MakeTree(tree *Tree, regex_ string, token string, alphabet *[]rune, wg *sync.WaitGroup) {
	// Define counter
	counter := tree.Counter
	// Convert the infix to postfix
	postfix := regex.InfixToPostfix(regex_)
	// Eval the alphabet
	EvalAlphabet(alphabet, postfix)
	// Create a stack to store the nodes
	stack := []*Node{}

	// Iterate over the postfix
	for i := 0; i < len(postfix); i++ {
		if postfix[i] == "." {
			counter++
			// Create a new node
			node := &Node{Value: ".", Left: stack[len(stack)-2], Right: stack[len(stack)-1], Ident: NumberToLetter(counter)}
			// Set nullable
			node.Nullable = stack[len(stack)-1].Nullable && stack[len(stack)-2].Nullable
			// Calc First
			if stack[len(stack)-2].Nullable {
				node.First = append(stack[len(stack)-2].First, stack[len(stack)-1].First...)
			} else {
				node.First = stack[len(stack)-2].First
			}

			// Calc Last
			if stack[len(stack)-1].Nullable {
				node.Last = append(stack[len(stack)-2].Last, stack[len(stack)-1].Last...)
			} else {
				node.Last = stack[len(stack)-1].Last
			}
			// Pop the last two elements
			stack = stack[:len(stack)-2]
			// Append the node to the stack
			stack = append(stack, node)
		} else if postfix[i] == "|" {
			counter++
			// Create a new node
			node := &Node{Value: "|", Left: stack[len(stack)-2], Right: stack[len(stack)-1], Ident: NumberToLetter(counter)}
			// Calc First
			node.First = append(stack[len(stack)-2].First, stack[len(stack)-1].First...)
			// Calc Last
			node.Last = append(stack[len(stack)-2].Last, stack[len(stack)-1].Last...)
			// Pop the last two elements
			stack = stack[:len(stack)-2]
			node.Nullable = node.Left.Nullable || node.Right.Nullable
			// Append the node to the stack
			stack = append(stack, node)
		} else if postfix[i] == "*" {
			counter++
			// Create a new node
			node := &Node{Value: "*", Left: stack[len(stack)-1], Ident: NumberToLetter(counter)}
			// Calc First
			node.First = stack[len(stack)-1].First
			// Calc Last
			node.Last = stack[len(stack)-1].Last
			// Set nullable
			node.Nullable = true
			// Pop the last element
			stack = stack[:len(stack)-1]
			// Append the node to the stack
			stack = append(stack, node)
		} else {
			counter++
			// Create a new node
			node := &Node{Value: postfix[i], Ident: strconv.Itoa(counter)}
			if postfix[i] == "epsilon" {
				node.Nullable = true
			} else {
				node.Nullable = false
			}
			// Calc First
			node.First = append(node.First, node)
			// Calc Last
			node.Last = append(node.Last, node)

			// Append the node to the stack
			stack = append(stack, node)
		}
	}

	// Create a token node
	tokenNode := &Node{Value: token, Ident: "TK"+strconv.Itoa(counter)+"", Nullable: false}
	tokenNode.First = append(tokenNode.First, tokenNode)
	tokenNode.Last = append(tokenNode.Last, tokenNode)
	// Create a new node
	node := &Node{Value: ".", Left: stack[0] , Right: tokenNode, Ident: "T"+strconv.Itoa(counter)}

	// Calc nullable
	node.Nullable = stack[0].Nullable && tokenNode.Nullable
	// Calc First
	if stack[0].Nullable {
		node.First = append(stack[0].First, tokenNode.First...)
	} else {
		node.First = stack[0].First
	}
	// Calc Last
	if tokenNode.Nullable {
		node.Last = append(tokenNode.Last, stack[0].Last...)
	} else {
		node.Last = tokenNode.Last
	}

	// Calc Follow
	CalcFollow(node)

	// Return the node
	tree.Root = node

	// Wait for the group
	wg.Done()

}

func CalcFollow(n *Node){
	if n.Value == "." {
		for _, node := range n.Left.Last {
			node.Follow = append(node.Follow, n.Right.First...)
		}
	} else if n.Value == "*" {
		for _, node := range n.Last {
			node.Follow = append(node.Follow, n.First...)
		}
	}

	if n.Left != nil {
		CalcFollow(n.Left)
	}
	if n.Right != nil {
		CalcFollow(n.Right)
	}
}


func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Value)
}

func CodeToJson(n *Node) string {
	// Convert the node to json
	jsonNode, _ := json.Marshal(n)
	fmt.Println(string(jsonNode))
	return string(jsonNode)
}

func AddNode(n *Node, graphAst *gographviz.Graph) {
	// Create a new node Ident is a string
	nodeName := "N"+n.Ident
	// Add the node to the graph
	graphAst.AddNode("AST", nodeName, map[string]string{"label": fmt.Sprintf("\"%v\"", n.Value)})
	// Check if the left node is not nil
	if n.Left != nil {
		// Add the left node to the graph
		AddNode(n.Left, graphAst)
		// Create a new edge
		graphAst.AddEdge(nodeName, "N"+n.Left.Ident, true, nil)
	}
	// Check if the right node is not nil
	if n.Right != nil {
		// Add the right node to the graph
		AddNode(n.Right, graphAst)
		// Create a new edge
		graphAst.AddEdge(nodeName, "N"+n.Right.Ident, true, nil)
	}

}

func ToGraph(n *Node){
	// Create a new graph
	graphAst := gographviz.NewGraph()
	// Set the name of the graph
	graphAst.SetName("AST")
	// Set the type of the graph
	graphAst.SetDir(true)
	// Add the node to the graph
	AddNode(n, graphAst)
	// Create a new file
	f, _ := os.Create("ast.dot")
	// Write the graph to the file
	f.WriteString(graphAst.String())
	// Close the file
	f.Close()
	

	// Make a pdf to visualize the graph
	cmd := "dot -Tpdf ast.dot -o ast.pdf"
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("cmd", "/C", cmd).Run()
	default: // Assume Unix-like system
		err = exec.Command("bash", "-c", cmd).Run()
	}

	if err != nil {
		fmt.Println(err)
	}

}


func MakeTreeFromMap(Tokens map[string]string) (*Node, []rune) {
	// Define counter
	counter := 0
	// Create a Alphabet rune
	var alphabet []rune
	nodes := []*Node{}

	var topTree *Node
	var wg sync.WaitGroup
	// Iterate over the Tokens
	for key, value := range Tokens {
		wg.Add(1)
		tr := &Tree{Counter: counter}
		MakeTree(tr, value, key, &alphabet, &wg)
		// Sum the counter number of values
		counter += len(value)
		nodes = append(nodes, tr.Root)
	}

	wg.Wait()

	for _, node := range nodes {
		if topTree == nil {
			topTree = node
		} else {
			topTree = &Node{Value: "|", Left: topTree, Right: node, Ident: "O"+strconv.Itoa(counter)}
			topTree.Nullable = topTree.Left.Nullable || topTree.Right.Nullable
			topTree.First = append(topTree.Left.First, topTree.Right.First...)
			topTree.Last = append(topTree.Left.Last, topTree.Right.Last...)
			CalcFollow(topTree) 
		}
	}

	// View Tree
	ToGraph(topTree)

	// Return the top tree
	return topTree, alphabet
}

