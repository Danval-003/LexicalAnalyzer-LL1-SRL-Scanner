package main

import (
	//"log"
	//"net/http"
	//"encoding/json"
	"backend/src/yalex"
	"backend/src/scanners"
	"fmt"
	"io"
	"os"
)

func main() {
	// Leer todo el contenido del archivo
	// Abrir el archivo
	file, err := os.Open("./examples/srl-1.yal")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	// Leer todo el contenido del archivo
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}


	yalex.Yal(string(content), "srl-1")

	// Simulate scanner
	encoded, err := scanners.SimulateScanner("srl-1/tokens", "var3*var4+var5")

	fmt.Println(encoded)


}



