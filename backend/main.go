package main

import (
	"fmt"
	//"log"
	//"net/http"
	//"encoding/json"
	"github.com/Danval-003/LexicalAnalyzer-LL1-SRL-Scanner/backend/src/regex/regexFormated" 
)

func main() {
	// Llamar a la función FormatRegex desde el paquete regex
	result := regex.InfixToPostfix("c(a+b)|d")
	fmt.Println(result)
}



