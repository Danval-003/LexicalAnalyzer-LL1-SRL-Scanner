package main


import (
	"fmt"
	"github.com/Danval-003/LexicalAnalyzer-LL1-SRL-Scanner/backend/src/afd"
)

func main() {
	YalexTokens := map[string]string{
		"COMMET": "\\(\\*[^'}']+\\*\\)",
	}
	YalexMachine,_, _ := afd.MakeAFD(YalexTokens)
	afd.VisualizeAFD(YalexMachine, "YalexMachine", "YalexMachine")
	simulate:= afd.SimulateAFD(YalexMachine, "(*Hola Mundo*)")

	// Iterate over the result
	for _, token := range simulate {
		fmt.Println("Token: ", token.Token, "Value: ", token.Accepted)
	}


}

