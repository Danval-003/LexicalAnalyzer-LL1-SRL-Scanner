package main

import (
	"fmt"
	//"log"
	//"net/http"
	//"encoding/json"
	"backend/regex" // Importa el paquete regex
)

func main() {
	// Llamar a la función FormatRegex desde el paquete regex
	result := regex.InfixToPostfix("c(a+b)|d")
	fmt.Println(result)
}



