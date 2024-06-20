package afd

import (
	"bytes"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"backend/src/tree"

	"github.com/goccy/go-graphviz"

	// Import to visualize the AFD with graphviz
	"github.com/awalterschulze/gographviz"

	// Import to encoding and decoding json
	"encoding/json"

	"sync"
)

// Struct to represent a state in the AFD
type State struct {
	Name        string              `json:"Name"`
	Transitions map[rune]*State     `json:"-"` 
	Accept      bool                `json:"Accept"`
	Token       string              `json:"Token"`
	TokenPrecedence int             
}

type Transition struct {
	Symbol string     `json:"Symbol"`
	NextState string  `json:"NextState"`
}

// Struct to represent a AFD in JSON
type AFDJSON struct {
	States []State `json:"States"`
	Transitions map[string][]Transition `json:"Transitions"`
}

// Struct to represent a State with a list of Nodes
type StateNodes struct {
	State *State
	Nodes []*tree.Node
}

// Struct to represent a State with a list of States
type StateStates struct {
	Name  string
	States []*State
}


// Function to sort Nodes in a list of Nodes
func SortNodes(nodes []*tree.Node) {
	nodesWithOutRepeat := []*tree.Node{}
	for _, node := range nodes {
		if !InList(node, nodesWithOutRepeat) {
			nodesWithOutRepeat = append(nodesWithOutRepeat, node)
		}
	}

	sort.Slice(nodesWithOutRepeat, func(i, j int) bool {
		return nodes[i].Ident < nodes[j].Ident
	})

	nodes = nodesWithOutRepeat
}

// Function to conver list of Nodes to hash
func NodesToHash(nodes []*tree.Node) string {
	SortNodes(nodes)
	hash := ""
	for _, node := range nodes {
		hash += node.Ident
		hash += ","
	}
	return hash
}

// Search in the list a Node with like "TK"+(Number) in this Ident returns true if it exists
func SearchToken(nodes []*tree.Node) bool {
	for _, node := range nodes {
		if len(node.Ident) < 2 {
			continue
		}
		if node.Ident[:2] == "TK" {
			return true
		}
	}
	return false
}

// Obtain a best token in the list of Nodes
func BestToken(nodes []*tree.Node) (string, int) {
	type result struct {
		token      string
		precedence int
	}

	numWorkers := 4 // Número de goroutines
	nodeCount := len(nodes)
	chunkSize := (nodeCount + numWorkers - 1) / numWorkers

	results := make(chan result, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > nodeCount {
			end = nodeCount
		}
		if start >= end {
			start = end
		}
		wg.Add(1)
		go func(nodes []*tree.Node) {
			defer wg.Done()
			localBestToken := ""
			localBestPrecedence := -1
			for _, node := range nodes {
				if len(node.Ident) < 2 {
					continue
				}
				if node.Ident[:2] == "TK" {
					value, err := strconv.Atoi(node.Ident[2:])
					if err == nil {
						if localBestPrecedence == -1 || value < localBestPrecedence {
							localBestToken = node.Value.(string)
							localBestPrecedence = value
						}
					}
				}
			}
			results <- result{localBestToken, localBestPrecedence}
		}(nodes[start:end])
	}

	wg.Wait()
	close(results)

	globalBestToken := ""
	globalBestPrecedence := -1

	for res := range results {
		if res.precedence != -1 && (globalBestPrecedence == -1 || res.precedence < globalBestPrecedence) {
			globalBestToken = res.token
			globalBestPrecedence = res.precedence
		}
	}

	return globalBestToken, globalBestPrecedence
}

// List to verify if a Node is in the list
func InList(node *tree.Node, list []*tree.Node) bool {
	for _, n := range list {
		if n == node {
			return true
		}
	}
	return false
}


// Yal function verify if a Machine exist
func MachineExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}


func AddState(state *State, VisitedStates map[string]*State, Graph *gographviz.Graph) {
    shape := "circle"
    if state.Accept {
        shape = "doublecircle"
    }
    Graph.AddNode("G", state.Name, map[string]string{"shape": shape})

    for symbol, nextState := range state.Transitions {
        if _, ok := VisitedStates[nextState.Name]; !ok {
            VisitedStates[nextState.Name] = nextState
            AddState(nextState, VisitedStates, Graph)
        }
        symbolString := ""
        if symbol == '"' {
            symbolString = "\"\\\"\""
        } else if symbol == '\\' {
            symbolString = "\"\\\\\""
        } else if symbol < 32 || symbol == 127 {
            symbolString = "\"\\x" + strconv.FormatInt(int64(symbol), 16) + "\""
        } else {
            symbolString = "\"" + string(symbol) + "\""
        }
        Graph.AddEdge(state.Name, nextState.Name, true, map[string]string{"label": symbolString})
    }
}

func VisualizeAFD(state *State) ([]byte, error) {
    g := gographviz.NewGraph()
    g.SetName("G")
    g.SetDir(true)
    g.AddAttr("G", "rankdir", "LR")

    VisitedStates := map[string]*State{}
    VisitedStates[state.Name] = state

    AddState(state, VisitedStates, g)

    s := g.String()

    graphvizGraph, err := graphviz.ParseBytes([]byte(s))
    if err != nil {
        return nil, err
    }

    gv := graphviz.New()
    defer gv.Close()

    var buf bytes.Buffer
    if err := gv.Render(graphvizGraph, graphviz.PNG, &buf); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}


// Function to Make the AFD from the tree, obtain a optional rune to name the states, for default is "Q"
func MakeAFD(Tokens map[string]string) (*State, []rune, StateSlice){
	// Obtain keys on order, and make a map with precedence
	keys := []string{}
	for key := range Tokens {
		keys = append(keys, key)
	}

	// Precedence map
	precedence := map[string]int{}
	for i, key := range keys {
		precedence[key] = i
	}
	// Take time
	treeA, alphabet := tree.MakeTreeFromMap(Tokens)
	// Create a Map to store the States
	states := map[string]*StateNodes{}

	// Count the number of states
	count := 0

	// Create the initial state
	intialKey := NodesToHash(treeA.First)
	initialState := &State{Name: "Q"+strconv.Itoa(count) , Transitions: map[rune]*State{}, Accept: SearchToken(treeA.First)}
	initialState.Token, initialState.TokenPrecedence = BestToken(treeA.First)
	count++
	states[intialKey] = &StateNodes{State: initialState, Nodes: treeA.First}

	// Create a list to store the states to be processed
	statesToProcess := []*StateNodes{states[intialKey]}
	// List of states
	statesList := []*State{}

	statesList = append(statesList, initialState)

	// Iterate over the states to process
	for len(statesToProcess) > 0 {
		actualList := statesToProcess[0].Nodes
		actualState := statesToProcess[0].State

		for _, symbol := range alphabet {
			// Create a list to store the nodes that will be in the new state
			newStateNodes := []*tree.Node{}

			// Iterate over the nodes in the actual state
			for _, node := range actualList {

				if node.Value == symbol {
					for _, nextNode := range node.Follow {
						if !InList(nextNode, newStateNodes) {
							newStateNodes = append(newStateNodes, nextNode)
						}
					}
				}
			}

			// Check if the new state is empty
			if len(newStateNodes) == 0 {
				continue
			}
			// Create a new state
			newStateKey := NodesToHash(newStateNodes)
			var nextState *State
			// Verify if the statekey is in the states
			if _, ok := states[newStateKey]; !ok {
				nextState = &State{Name: "Q"+strconv.Itoa(count) , Transitions: map[rune]*State{}, Accept: SearchToken(newStateNodes)}
				states[newStateKey] = &StateNodes{State: nextState, Nodes: newStateNodes}
				statesList = append(statesList, nextState)
				nextState.Token, nextState.TokenPrecedence = BestToken(newStateNodes)
				statesToProcess = append(statesToProcess, &StateNodes{State: nextState, Nodes: newStateNodes})
				count++
			} else {
				nextState = states[newStateKey].State
			}

			// Add the transition to the actual state
			actualState.Transitions[int32(symbol)] = nextState
		}

		// Remove the actual state from the list
		if len(statesToProcess) > 0 {
			statesToProcess = statesToProcess[1:]
		
		} else {
			break
		}
	}



	return states[intialKey].State, alphabet, statesList
}

// Func to Encode the State to JSON
type StateSlice []*State

func (states StateSlice) Encode() string {
	AFDJSON := &AFDJSON{States: []State{}, Transitions: map[string][]Transition{}}
	// Iterate over the states
	for _, state := range states {
		AFDJSON.States = append(AFDJSON.States, *state)
		// Iterate over the transitions
		for symbol, nextState := range state.Transitions {
			if AFDJSON.Transitions[state.Name] == nil {
				AFDJSON.Transitions[state.Name] = []Transition{}
			}
			AFDJSON.Transitions[state.Name] = append(AFDJSON.Transitions[state.Name], Transition{Symbol: string(symbol), NextState: nextState.Name})
		}
	}

	// Encode the AFD to JSON
	AFDJSONBytes, _ := json.Marshal(AFDJSON)
	return string(AFDJSONBytes)
}

// Function to Decode the AFD from JSON
func DecodeAFD(jsonState string) *State {

	AFDJSON := &AFDJSON{}
	json.Unmarshal([]byte(jsonState), AFDJSON)

	// Create a map to store the states
	states := map[string]*State{}

	// Iterate over the states
	for _, state := range AFDJSON.States {
		states[state.Name] = &state
	}

	// Iterate over the transitions
	for state_From, transition := range AFDJSON.Transitions {
		// Iterate over the transitions
		for _, transition := range transition {
			if states[state_From].Transitions == nil {
				states[state_From].Transitions = map[rune]*State{}
			}

			states[state_From].Transitions[rune(transition.Symbol[0])] = states[transition.NextState]
		}
	}

	return states["Q0"]
}


// Function to Save the AFD to a file
func SaveMachine(filename string, jsonState StateSlice) {
	// Open the file
	file, _ := os.Create(filename)
	// Write the json to the file
	file.WriteString(jsonState.Encode())
	file.Close()
}

// Function to Load the AFD from a file
func LoadMachine(filename string) *State {
	// Open the file
	file, _ := os.Open(filename)
	// Read the content of the file
	content, _ := io.ReadAll(file)
	// Decode the AFD
	return DecodeAFD(string(content))
}

// Struct to represent a Simulated part
type SimulatedPart struct {
	Init int
	Runes []rune	`json:"runes"`
	Final int    	
	Token string	`json:"token"`
	Accepted bool   `json:"accepted"`
}

type SimulatedPartJSON struct {
	Text string     `json:"text"`
	Accepted bool   `json:"accepted"`
	Token string    `json:"token"`
}

// Alias to a list SimulatedPart
type SimulatedParts []*SimulatedPart

// Response Sim
type ResponseSim struct{
	SimulatedParts SimulatedParts `json:"simulatedParts"`
	Accepted bool `json:"accepted"`
	StringSummary string  `json:"stringSummary"`
}

// Parts in the simulation
type SimulationJSON struct {
	Parts []SimulatedPartJSON `json:"parts"`
	Accept bool              `json:"accept"`
}

// Function to convert SimulatedParts into a JSON
func (simulatedParts SimulatedParts) Encode() (string, error) {
	simulatedPartsJSON := []SimulatedPartJSON{}
	accep := true
	// Iterate over the SimulatedParts
	for _, simulatedPart := range simulatedParts {
		simulatedPartJSON := SimulatedPartJSON{Text: string(simulatedPart.Runes), Accepted: simulatedPart.Accepted, Token: simulatedPart.Token}
		simulatedPartsJSON = append(simulatedPartsJSON, simulatedPartJSON)
		accep = accep && simulatedPart.Accepted
	}

	simulation := SimulationJSON{Parts: simulatedPartsJSON, Accept: accep}

	// Encode the SimulatedParts to JSON
	simulatedPartsBytes, err := json.Marshal(simulation)

	return string(simulatedPartsBytes), err
}


// Function to simulate a part of string in the AFD
func SimulateAFDPart(state *State, runes []rune, init int, wg *sync.WaitGroup, simulate *SimulatedPart) {
	// Stack to store the saved states
	actualState := state
	// token to store the token
	token := ""

	simulate.Init = init
	simulate.Runes = runes[init: init+1]
	simulate.Final = init 
	// Runes simulated
	for i := init; i < len(runes); i++ {
		// Check if the rune is in the transitions
		if nextState := actualState.Transitions[runes[i]]; nextState != nil {
			actualState = nextState
			if actualState.Accept {
				token = actualState.Token
				simulate.Final = i
				simulate.Token = token
				simulate.Accepted = true
				simulate.Runes = runes[init:i+1]
			}
		} else {
			break
		}
	}

	wg.Done()

}


// Function to simulate the AFD with a string
func SimulateAFD(state *State, stringToSimulate string) SimulatedParts {
	// Change the string to a list of runes
	runes := []rune(stringToSimulate)
	// Stackto store the saved states
	actualState := state
	// User goroutines to simulate the AFD
	wg := sync.WaitGroup{}
	// Index of the string
	index := 0
	simulates := []*SimulatedPart{}
	// Iterate over the string
	for index < len(runes) {
		simulate := &SimulatedPart{}
		wg.Add(1)
		go SimulateAFDPart(actualState, runes, index, &wg, simulate)
		simulates = append(simulates, simulate)
		index++
	}
	wg.Wait()
	simulatedParts := SimulatedParts{}
	accepted:=true
	// Iterate over the simulates
	simu := -1
	for _, simulate := range simulates {
		if simulate.Init > simu {
			simu= simulate.Final
			simulatedParts = append(simulatedParts, simulate)
			accepted = accepted && simulate.Accepted
			
		}
	}



	return simulatedParts

}


// Function to simulate the AFD with a string
func ResponseSimulateAFD(state *State, stringToSimulate string) ResponseSim {
	// Change the string to a list of runes
	runes := []rune(stringToSimulate)
	// Stackto store the saved states
	actualState := state
	// User goroutines to simulate the AFD
	wg := sync.WaitGroup{}
	// Index of the string
	index := 0
	simulates := []*SimulatedPart{}
	// Iterate over the string
	for index < len(runes) {
		simulate := &SimulatedPart{}
		wg.Add(1)
		go SimulateAFDPart(actualState, runes, index, &wg, simulate)
		simulates = append(simulates, simulate)
		index++
	}
	wg.Wait()
	simulatedParts := SimulatedParts{}
	accepted:=true
	// Iterate over the simulates
	simu := -1
	simString := ""
	for _, simulate := range simulates {
		if simulate.Init > simu {
			simu= simulate.Final
			simulatedParts = append(simulatedParts, simulate)
			accepted = accepted && simulate.Accepted
			simString += simulate.Token+" "
		}
	}

	// Remove the spaces TrimSpace
	simString = strings.TrimSpace(simString)
	

	return ResponseSim{SimulatedParts: simulatedParts, Accepted: accepted, StringSummary: simString}

}

 