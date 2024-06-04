package main

import (
	"fmt"

	"github.com/Danval-003/LexicalAnalyzer-LL1-SRL-Scanner/backend/src/tree"
)

// Struct to represent a state in the AFD
type State struct {
	Name string
	Transitions map[string]*State
	Accept bool
	Token string
}

func MakeAFD(Tokens map[string]string) {

	treeA := tree.MakeTreeFromMap(Tokens)

}


func main() {
	Tokens := map[string]string{}
	Tokens["ID"] = "['a'-'z']['a'-'z']*"
	Tokens["NUM"] = "['0'-'9']['0'-'9']*"
	MakeAFD(Tokens)
}


