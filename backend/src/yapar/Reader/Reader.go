package reader

import (
	"backend/src/afd"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	// Import to visualize the LR0 with graphviz
	"github.com/awalterschulze/gographviz"
)

type PRODUCTION struct {
	NoTerminal string   `json:"noTerminal"`
	Production []string `json:"production"`
	Pointer int         `json:"pointer"`
	First []string      
	Follow []string     
}

type StateLR0 struct{
	Productions []PRODUCTION
	Transitions map[string]*StateLR0
	Number int
}


// Function to Visualize LR0 in Graphviz
func VisualizeLR0(states *map[string]*StateLR0, initStateName string, fileName string) (string, error) {
	// Create a new graph
	graph := gographviz.NewGraph()
	// Set the name
	graph.SetName(fileName)
	// Set the rankdir
	graph.SetDir(true)

	// Set horizontal
	graph.Attrs.Add("rankdir", "LR")

	// Create a map of string
	statesMap := make(map[string]string)

	// Iterate over states
	for _, state := range *states {
		// Create a node
		nodeName := "S" + fmt.Sprint(state.Number)
		// Set the node with box shape and escape special characters in labels
		label := strings.ReplaceAll(state.ToString(), "\"", "\\\"")
		graph.AddNode(fileName, nodeName, map[string]string{"label": "\"" + nodeName + "\\n" + label + "\"", "shape": "box"})
		// Append in statesMap
		statesMap[state.ToString()] = nodeName
	}

	// Iterate over states
	for _, state := range *states {
		// Iterate over state.Transitions
		for symbol, transition := range state.Transitions {
			// Set the edge and escape special characters in labels
			symbolLabel := strings.ReplaceAll(symbol, "\"", "\\\"")
			graph.AddEdge(statesMap[state.ToString()], statesMap[transition.ToString()], true, map[string]string{"label": "\"" + symbolLabel + "\""})
		}
	}

	// Create an arrow to start in nothing and focus the init state
	graph.AddNode(fileName, "nothing", map[string]string{"label": "\"\"", "shape": "point"})
	graph.AddEdge("nothing", statesMap[initStateName], true, map[string]string{"label": "\"\""})

	s := graph.String()

	// URL encode the string to make it safe for the API
	encoded := url.QueryEscape(s)

	return encoded, nil
}

// Function to Calc First of productions
func CalcFirst(productionsMap map[string][]PRODUCTION, noTerminal string) []string {
	// Create a list of string
	first := make([]string, 0)

	// Create a List with noTerminals to evaluate
	noTerminalsToEvaluate := []string{noTerminal}

	// Create a List with noTerminals evaluated	
	noTerminalsEvaluated := make(map[string]bool)

	for len(noTerminalsToEvaluate) > 0 {
		// Obtain the first element
		noTerminal := noTerminalsToEvaluate[0]

		// Remove the first element
		noTerminalsToEvaluate = noTerminalsToEvaluate[1:]

		// Verify if the noTerminal is in noTerminalsEvaluated
		if _, ok := noTerminalsEvaluated[noTerminal]; !ok {
			// Append in noTerminalsEvaluated
			noTerminalsEvaluated[noTerminal] = true

			// Verify if the noTerminal is in productionsMap
			if productions, ok := productionsMap[noTerminal]; ok {
				for _, production := range productions {
					// Verify if the pointer is in the end
					if production.Pointer < len(production.Production) {
						// Obtain the symbol in the pointer
						symbol := production.Production[production.Pointer]

						// Verify if the symbol is in tokensMap
						if _, ok := productionsMap[symbol]; ok {
							// Verify if the symbol is in noTerminalsEvaluated
							if _, ok := noTerminalsEvaluated[symbol]; !ok {
								// Append in noTerminalsToEvaluate
								noTerminalsToEvaluate = append(noTerminalsToEvaluate, symbol)
							}
						} else {
							
							// Append in first
							first = append(first, symbol)
						}
					}
				}
			}
		}
	
	}

	

	return first

}

// Function to calc Follow
func CalcFollow(productionsMap map[string][]PRODUCTION, noTerminal string, evaled map[string]bool) []string {
	// Create a list of string
	follow := make([]string, 0)

	// Create a List with noTerminals to evaluate
	noTerminalsToEvaluate := []string{noTerminal}

	// Create a List with noTerminals evaluated	
	noTerminalsEvaluated := make(map[string]bool)

	for len(noTerminalsToEvaluate) > 0 {
		// Obtain the first element
		noTerminal := noTerminalsToEvaluate[0]

		// Remove the first element
		noTerminalsToEvaluate = noTerminalsToEvaluate[1:]

		// Verify if the noTerminal is in noTerminalsEvaluated
		if _, ok := noTerminalsEvaluated[noTerminal]; !ok {
			// Append in noTerminalsEvaluated
			noTerminalsEvaluated[noTerminal] = true

			// Iterate over productionsMap
			for keyPr, productions := range productionsMap{
				for _, production := range productions{
					for i, symbol := range production.Production{
						if symbol == noTerminal{
							if i+1 < len(production.Production){
								// Obtain the next symbol
								nextSymbol := production.Production[i+1]

								// Verify if nextSymbol is in productionsMap
								if _, ok := productionsMap[nextSymbol]; ok{
									// Combine the list with First
									follow = append(follow, CalcFirst(productionsMap, nextSymbol)...)
								} else {
									follow = append(follow, nextSymbol)
								}
							} else {
								// Verify if the keyPr is different to noTerminal
								if keyPr != noTerminal{
									// Verify if the keyPr not are in evaled
									if _, ok := evaled[keyPr]; !ok{
										// Append in evaled
										evaled[keyPr] = true
										// Combine the list with Follow
										follow = append(follow, CalcFollow(productionsMap, keyPr, evaled)...)
									}

								}
							}
						}
					}
				}
			}
		}
	
	}

	

	return follow

}

// Function to convert PRODUCTION to JSON
func (p PRODUCTION) ToJSON() string {
	json, _ := json.Marshal(p)
	return string(json)
}


// Func to Convert PRODUCTION to String
func (p PRODUCTION) ToString() string {
	// Copy production
	production:= make([]string, len(p.Production))
	copy(production, p.Production)
	// Add a "." into the pointer
	production = append(production[:p.Pointer], append([]string{"."}, production[p.Pointer:]...)...)

	return p.NoTerminal + " -> " + strings.Join(production, " ")
}

// Func to sort a Production List on base into ToString() func
func SortProductionList(productionList []PRODUCTION) []PRODUCTION {
	// Create a map of PRODUCTION
	productionMap := make(map[string]PRODUCTION)

	// Create a list of PRODUCTION
	productionListSorted := make([]PRODUCTION, 0)

	// Strings
	sProduct := []string{}

	// Iterate over productionList and append in productionMap
	for _, production := range productionList {
		productionMap[production.ToString()] = production
		sProduct = append(sProduct, production.ToString())
	}

	// Sort the strings
	sort.Strings(sProduct)

	// Iterate over sProduct and append in productionListSorted
	for _, s := range sProduct {
		productionListSorted = append(productionListSorted, productionMap[s])
	}

	return productionListSorted
}


// Function to EClouser
func EClouser(productionsMap map[string][]PRODUCTION, noTerminal string, original PRODUCTION) []PRODUCTION {
	// Create a list of PRODUCTION
	productionList := make([]PRODUCTION, 0)

	// Append the original production
	productionList = append(productionList, original)

	// Create a List with nonTerminals explored
	noTerminals := make(map[string]bool)

	// Create a List with nonTerminals to explore
	noTerminalsToExplore := []string{noTerminal}

	for len(noTerminalsToExplore) > 0 {
		// Obtain the first element
		noTerminal := noTerminalsToExplore[0]

		// Remove the first element
		noTerminalsToExplore = noTerminalsToExplore[1:]

		// Explore the list in Productions
		if productions, ok := productionsMap[noTerminal]; ok {
			for _, production := range productions {
				// Append in productionList
				productionList = append(productionList, production)
				// Append in noTerminals
				noTerminals[noTerminal] = true

				// Verify if the pointer is in the end
				if production.Pointer < len(production.Production) {
					// Obtain the symbol in the pointer
					symbol := production.Production[production.Pointer]

					// Verify if the symbol is in Productionsmap
					if _, ok := productionsMap[symbol]; ok {
						// Verify if the symbol is in noTerminals
						if _, ok := noTerminals[symbol]; !ok {
							// Append in noTerminalsToExplore
							noTerminalsToExplore = append(noTerminalsToExplore, symbol)
						}
					}
				}
			}
		}

	}

	return SortProductionList(productionList)
}


// Function to combine two lists of PRODUCTION, removing duplicates
func CombineProductionLists(list1 []PRODUCTION, list2 []PRODUCTION) []PRODUCTION {
	// Create a map of PRODUCTION
	productionMap := make(map[string]PRODUCTION)

	// Create a list of PRODUCTION
	productionList := make([]PRODUCTION, 0)

	// Iterate over list1 and append in productionMap
	for _, production := range list1 {
		productionMap[production.ToString()] = production
	}

	// Iterate over list2 and append in productionMap
	for _, production := range list2 {
		productionMap[production.ToString()] = production
	}

	// Iterate over productionMap and append in productionList
	for _, production := range productionMap {
		productionList = append(productionList, production)
	}

	return SortProductionList(productionList)
}

// Function to convert StateLR0 into String
func (s StateLR0) ToString() string {
	// Create a list of string
	productions := make([]string, 0)

	// Iterate over s.Productions and append in productions
	for _, production := range s.Productions {
		productions = append(productions, production.ToString())
	}

	// Sort the productions
	sort.Strings(productions)

	// Create a string
	sProductions := ""

	// Iterate over productions and append in sProductions
	for _, production := range productions {
		sProductions += production + "\n"
	}

	return sProductions
}


// Function to Read a file .yalp
func ReadFileYalp(content string) (*StateLR0, map[string]*StateLR0, []string, map[string][]PRODUCTION, map[string]string, error){

	// Create the machine
	YalpMachine := map[string]string{
		"TOKENS": "'%''t''o''k''e''n'(' '['A'-'Z']+)+",
		"COMMENTS":"'/'\\*|\\*'/'",
		"PRODUCTION": "['a'-'z']+':'(['a'-'z''A'-'Z']|' '|'\n'|\\|)+';'",
		"IGNORE": "\"IGNORE\"(' '['A'-'Z']+)+",
		"WS": "[' ''\n''\t']",
	}

	machine, _, _ := afd.MakeAFD(YalpMachine)


	// Simulate content with machine
	simulate:= afd.SimulateAFD(machine, content)

	inComment := false

	// Create a tokens List
	tokensMap:= make(map[string]string)

	// Create a IGNORE list
	ignoreMap:= make(map[string]string)

	// Create a PRODUCTION list
	productionsMap:= make(map[string][]PRODUCTION)

	// prsList
	prsList:= make(map[string][]string)

	// Init exp
	initPr:= ""

	// Symbols
	symbols:= []string{}

	// Iterate over simulate and print accepted tokens
	for _, token := range simulate {
		switch token.Token {
		case "COMMENTS":
			if string(token.Runes) == "*/" {
				inComment = false
			} else {
				inComment = true
			}
		case "TOKENS":
			if !inComment {
				tokenString := string(token.Runes)
				// Remove %token and spaces (strip)
				tokenString = strings.TrimSpace(tokenString[6:])
				// Split by spaces
				tokens := strings.Split(tokenString, " ")

				for _, tok := range tokens {
					tokensMap[tok] = "\"" + tok + "\""
					symbols = append(symbols, tok)
					_ = symbols // Assign the result of append to a variable to avoid the SA4010 error
				}

			}
		case "IGNORE":
			if !inComment {
				tokenString := string(token.Runes)
				// Remove IGNORE and spaces (strip)
				tokenString = strings.TrimSpace(tokenString[6:])
				// Split by spaces
				ignores := strings.Split(tokenString, " ")

				for _, ignore := range ignores {
					ignoreMap[ignore] = "\"" + ignore + "\""
				}
			}

		case "PRODUCTION":
			if !inComment {
				// Delete ';' in the end
				token.Runes = token.Runes[:len(token.Runes)-1]
				// Split by ':'
				production := strings.SplitN(string(token.Runes), ":", 2)
				if len(production) < 2 {
					fmt.Println("Error: invalid production format")
					continue
				}

				// The other part split by '|' and TrimSpace all
				productions := strings.Split(production[1], "|")
				for i, prod := range productions {
					productions[i] = strings.TrimSpace(prod)
					if initPr == "" {
						initPr = strings.TrimSpace(production[0])
					}
				}
				production[0] =strings.TrimSpace(production[0])

				symbols = append(symbols, production[0])

				// Append in prsList
				prsList[production[0]] = productions
			}

		case "WS":
			// Do nothing

		default:
			if !inComment {
				fmt.Printf("error: Token not recognized - %s\n", token.Token)
			}

		}
	}


	// Iterate over prsList and create a PRODUCTION list
	for noTerminal, productions := range prsList{
		// Create a list of PRODUCTION
		var productionList []PRODUCTION

		for _, production := range productions{
			
			// Create a PRODUCTION
			prod:= PRODUCTION{
				NoTerminal: noTerminal,
				Pointer: 0,
				Production: []string{},
			}

			for _, symbol := range strings.Split(production, " ") {
				symbolString:= strings.TrimSpace(string(symbol))
				// Verify if are in tokens Map or prsList
				if _, ok := tokensMap[symbolString]; ok{
					prod.Production = append(prod.Production, symbolString)
				} else if _, ok := prsList[symbolString]; ok{
					prod.Production = append(prod.Production, symbolString)
				} else {
					return nil,nil, nil, nil, nil, fmt.Errorf("error: Symbol not recognized")
				}
			}

			// Append in productionList
			productionList = append(productionList, prod)
		}

		// Append in productionsMap
		productionsMap[noTerminal] = productionList
	}

	// Create Init Production
	initProduction := PRODUCTION{NoTerminal: initPr+"'", Pointer: 0, Production: []string{initPr, "$"}}
	productionsMap[initPr+"'"] = []PRODUCTION{ initProduction}

	initPr = initPr + "'"

	// Iterate over Production Map and calc First
	for keyPr, productions := range productionsMap{
		for IndexPr, production := range productions{
			production.First = CalcFirst(productionsMap, production.NoTerminal)
			productionsMap[keyPr][IndexPr] = production
		}
	}

	// Iterate over Production Map and calc Follow
	for keyPr, productions := range productionsMap{
		for IndexPr, production := range productions{
			production.Follow = CalcFollow(productionsMap, production.NoTerminal, map[string]bool{keyPr: true})
			productionsMap[keyPr][IndexPr] = production
		}
	}


	initState := EClouser(productionsMap, initPr, productionsMap[initPr][0])

	cout:= 0
	// Create a Map to store the States
	initStateLR0 := StateLR0{Productions: initState, Transitions: map[string]*StateLR0{}, Number: cout} 
	statesMap := make(map[string]*StateLR0)
	statesMap[initStateLR0.ToString()] = &initStateLR0

	// Create a list with states to explore
	statesToExplore := []StateLR0{initStateLR0}

	// Iterate over statesToExplore
	for len(statesToExplore) > 0 {
		// Obtain the first element
		state := &statesToExplore[0]

		// Remove the first element
		statesToExplore = statesToExplore[1:]

		// Create a map of string
		transitions := &state.Transitions

		// Iterate over Symbols
		for _, symbol := range symbols{
			// Create a list of PRODUCTION
			productionList := make([]PRODUCTION, 0)

			// Iterate over state.Productions
			for _, production := range state.Productions{
				// Verify if the pointer is in the end
				if production.Pointer < len(production.Production){
					// Obtain the symbol in the pointer
					symbolPointer := production.Production[production.Pointer]

					// Verify if the symbolPointer is equal to symbol
					if symbolPointer == symbol{
						// Create a PRODUCTION
						prod:= PRODUCTION{
							NoTerminal: production.NoTerminal,
							Pointer: production.Pointer + 1,
							Production: production.Production,
							First: production.First,
							Follow: production.Follow,
						}

						

						// verify if Pointer is in the end
						if prod.Pointer < len(prod.Production){
							nonTerminal := prod.Production[prod.Pointer]
							// Verify if nonTerminal is in productionsMap
							if _, ok := productionsMap[nonTerminal]; ok{
								// Combine the list with EClouser
								productionList = CombineProductionLists(productionList, EClouser(productionsMap, nonTerminal, prod))
							} else {
								productionList = append(productionList, prod)
							}
						} else {
							productionList = append(productionList, prod)
						}
						
					}
				}
			}

			// Verify if productionList is not empty
			if len(productionList) > 0{
				// Create LR0
				lr0:= StateLR0{Productions: productionList, Transitions: map[string]*StateLR0{}, Number: cout+1}
				// Verify if new lr0 Exist
				if _, ok := statesMap[lr0.ToString()]; !ok{
					// Append in statesMap
					statesMap[lr0.ToString()] = &lr0
					// Append in statesToExplore
					statesToExplore = append(statesToExplore, lr0)
					// Append in transitions
					(*transitions)[symbol] = &lr0
					// Increment cout
					cout += 1
				} else {
					lro:= statesMap[lr0.ToString()]
					// Append
					(*transitions)[symbol] = lro
					
				}
			}
		}
		
	
	}

	VisualizeLR0(&statesMap, initStateLR0.ToString(), "LR0")


	return &initStateLR0, statesMap, symbols, productionsMap, ignoreMap, nil
}

