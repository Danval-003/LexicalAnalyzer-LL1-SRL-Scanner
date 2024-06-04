package afd

import (
	"fmt"
	"strconv"

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


}


