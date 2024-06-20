package yapar

import (
	"backend/src/yapar/Reader"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"encoding/json"

	"github.com/bndr/gotabulate"
	"github.com/go-gota/gota/dataframe"
)

// Struct to represent action
type Action struct {
	Action string  `json:"action"`
	Symbol string  `json:"symbol"`
	Number int     `json:"number"`
}


// Function to convert Action to string
func (a Action) String() string {
	return fmt.Sprintf("Action: %s, Symbol: %s, Number: %d", a.Action, a.Symbol, a.Number)
}

// Alias to map[int]map[string]Action
type Table map[int]map[string]Action

// Struct to represent SLR
type SLR struct {
	Table Table  `json:"table"`
	Ignored map[string]string  `json:"ignored"`
}

// Function to encoding table to json
func EncodeTable(table Table) (string, error) {
	// Convert table to json
	encoded, err := json.Marshal(table)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

// Funtion to convert json to Table
func DecodeTable(encoded string) (Table, error) {
	// Convert json to table
	var table Table
	err := json.Unmarshal([]byte(encoded), &table)
	if err != nil {
		return nil, err
	}

	return table, nil
}

// Function to print table
func PrintTable(table map[int]map[string]Action) {
	// Convert table into a map 
	var records []map[string]interface{}
    for key, innerMap := range table {
        for action, data := range innerMap {
            record := map[string]interface{}{
                "Key":    key,
                "Action": action,
                "Data":   data.String(),
            }
            records = append(records, record)
        }
    }

	// Create a dataframe from the records
	df := dataframe.LoadMaps(records)

	// Imprimir el DataFrame
    fmt.Println(df)

    // Convertir el DataFrame a un slice de slices para usar con gotabulate
    recordsForTabulate := df.Records()

    // Crear una nueva instancia de Tabulate
    tabulate := gotabulate.Create(recordsForTabulate)
    tabulate.SetHeaders(recordsForTabulate[0])
    tabulate.SetAlign("left")
    tabulate.SetWrapStrings(false)

    // Imprimir la tabla en formato tabular
    fmt.Println(tabulate.Render("grid"))
}


func ObtainSrl(yapar string, fileName string) (string, SLR, error) {
	var slr SLR
    // ReadYapar
    init, states, _, productsMap, ignoreMap, _ := reader.ReadFileYalp(yapar)
	// Delete spaces on fileName
	fileName = strings.ReplaceAll(fileName, " ", "")

	digraph, err:=reader.VisualizeLR0(&states, init.ToString(), fileName)
	if err != nil {
		return "", slr, err
	}

	urlImage:= "https://quickchart.io/graphviz?graph=" + digraph

    // Create the table -> map of maps of strings
    table := make(map[int]map[string]Action)
    var mu sync.Mutex

    var wg sync.WaitGroup
    errChan := make(chan error, len(states))

    // Create the table
    for _, state := range states {
        wg.Add(1)
        go func(state *reader.StateLR0) {
            defer wg.Done()
            localTable := make(map[string]Action)

            // Iterate over productions
            for _, production := range state.Productions {
                // Verify if the point is at the end
                if production.Pointer == len(production.Production) {
                    // Iterate over Follow
                    for _, follow := range production.Follow {
                        mu.Lock()
                        if _, ok := localTable[follow]; ok {
                            mu.Unlock()
                            errChan <- fmt.Errorf("error en la tabla de análisis sintáctico, tipo: " + localTable[follow].Action + "R")
                            return
                        }
                        localTable[follow] = Action{Action: "R", Symbol: production.NoTerminal, Number: len(production.Production)}
                        mu.Unlock()
                    }
                } else {
                    symbolToTrans := production.Production[production.Pointer]

                    if symbolToTrans == "$" {
                        mu.Lock()
                        localTable[symbolToTrans] = Action{Action: "A", Symbol: symbolToTrans}
                        mu.Unlock()
                    } else {
                        mu.Lock()
                        stateToTrans, ok := state.Transitions[symbolToTrans]
                        mu.Unlock()
                        if !ok {
                            errChan <- fmt.Errorf("error en la tabla de análisis sintáctico, no existe transición para: " + production.Production[production.Pointer])
                            return
                        }
                        mu.Lock()
                        if _, ok := productsMap[symbolToTrans]; ok {
                            localTable[symbolToTrans] = Action{Action: "GOTO", Symbol: symbolToTrans, Number: stateToTrans.Number}
                        } else {
                            localTable[symbolToTrans] = Action{Action: "S", Symbol: symbolToTrans, Number: stateToTrans.Number}
                        }
                        mu.Unlock()
                    }
                }
            }
            mu.Lock()
            table[state.Number] = localTable
            mu.Unlock()
        }(state)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return "",slr, err
        }
    }

    return urlImage, SLR{Table: table, Ignored: ignoreMap}, nil
}


// Function to save table
func SaveTable(table map[int]map[string]Action, name string) error {
	// Create a path
	path := "./src/tables/"+name+".json"
	// Convert table to json
	encoded, err := EncodeTable(table)
	if err != nil {
		return err
	}

	// Verify if file exist
	if _, err := os.Stat(path); err == nil {
		// Delete file
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	// Save table in file
	_, err = file.WriteString(encoded)
	if err != nil {
		return err
	}
	// Close file
	err = file.Close()

	return err

}

// Load Table
func LoadTable(name string) (map[int]map[string]Action, error) {
	// Create a path
	path := "./src/tables/"+name+".json"
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// Read file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// Convert json to table
	table, err := DecodeTable(string(content))
	if err != nil {
		return nil, err
	}

	return table, nil

}


// Function to print the simulate map
func PrintSimulateMap(table map[string][]interface{}) {
	maxLength := 0
    for _, column := range table {
        if len(column) > maxLength {
            maxLength = len(column)
        }
    }

    // Crear registros y llenar con valores nulos donde sea necesario
    records := make([]map[string]interface{}, maxLength)
    for i := 0; i < maxLength; i++ {
        record := make(map[string]interface{})
        for key, column := range table {
            if i < len(column) {
                record[key] = column[i]
            } else {
                record[key] = nil
            }
        }
        records[i] = record
    }

    // Crear un DataFrame a partir de los registros
    df := dataframe.LoadMaps(records)



    // Convertir el DataFrame a un slice de slices para usar con gotabulate
    recordsForTabulate := df.Records()

    // Crear una nueva instancia de Tabulate
    tabulate := gotabulate.Create(recordsForTabulate[1:])
    tabulate.SetHeaders(recordsForTabulate[0])
    tabulate.SetAlign("left")
    tabulate.SetWrapStrings(false)

    // Imprimir la tabla en formato tabular
    fmt.Println(tabulate.Render("grid"))
	// Convert Table into String
}

// Struct to represent Step
type Step struct {
	State map[string][]interface{}  `json:"state"`
	Number int                      `json:"number"`
}

// Struct to represent Simulate Response
type SimResponse struct {
	Steps []Step  `json:"steps"`
	Accept bool  `json:"accept"`
}


// Function to simulate Table
func SimulateTable(slr SLR, input string) ([]byte ,error) {
	table := slr.Table
	toIgnore := slr.Ignored
	// Create a map with the stack and action
	mapSimulate := make(map[string][]interface{})
	// Create a stack
	mapSimulate["stack"] = []interface{}{0}
	// Create a input, split by space
	input = input + " $"
	inputSplit := strings.Split(input, " ")
	interFaceInput := make([]interface{}, len(inputSplit))
	for i, v := range inputSplit {
		interFaceInput[i] = v
	}
	mapSimulate["input"] = interFaceInput
	// Create a action
	mapSimulate["action"] = []interface{}{"-"}

	// List of steps
	steps := []Step{}

	// Accept
	accept := false


	// Iterate over the input
	for len(mapSimulate["input"]) > 0 {
		state:= mapSimulate["stack"][0].(int)
		symbol := mapSimulate["input"][0].(string)

		// Verify if symbol is in the ignored
		if _, ok := toIgnore[symbol]; ok {
			mapSimulate["input"] = mapSimulate["input"][1:]
			continue
		}

		// Verify if exists action
		action, ok := table[state][symbol]
		// Change the action
		mapSimulate["action"] = []interface{}{action.Action}
		

		if !ok {
			step := Step{ Number: len(steps)}
			// Create a copy of mapSimulate
			step.State = make(map[string][]interface{})
			for key, value := range mapSimulate {
				step.State[key] = make([]interface{}, len(value))
				copy(step.State[key], value)
			}
			steps = append(steps, step)
			_ = steps
			// List for expected tokens
			runeExpected := []string{}
			// Iterate over the table
			for _, value := range table[state] {
				// Verify if value not are in Uppercase
				if value.Symbol != strings.ToUpper(value.Symbol) {
					continue
				}
				runeExpected = append(runeExpected, value.Symbol)
			}
			// Convert list to string
			expected := strings.Join(runeExpected, ", ")

			// Error
			if symbol == "$" {
				return nil, fmt.Errorf("error en la tabla de análisis sintáctico, se esperaban los tokens: "+expected)
			}
			return nil, fmt.Errorf("error en la tabla de análisis sintáctico, no existe acción para: "+symbol+" se esperaban los tokens: "+expected)
		}
		step := Step{ Number: len(steps)}
		// Create a copy of mapSimulate
		step.State = make(map[string][]interface{})
		for key, value := range mapSimulate {
			step.State[key] = make([]interface{}, len(value))
			copy(step.State[key], value)
		}
		steps = append(steps, step)

		// Verify if action is accept
		if action.Action == "A" {
			accept = true
			break
		}

		// Verify if action is shift
		if action.Action == "S" {
			// Add to stack number on first
			mapSimulate["stack"] = append([]interface{}{action.Number}, mapSimulate["stack"]...)
			// Remove from input
			mapSimulate["input"] = mapSimulate["input"][1:]
		}

		// Verify if action is reduce
		if action.Action == "R" {
			// Remove from stack number
			number := action.Number
			// Remove from stack
			mapSimulate["stack"] = mapSimulate["stack"][number:]
			// Add to stack
			state = mapSimulate["stack"][0].(int)

			// Verify if exists action
			action, ok := table[state][action.Symbol]
			// Change the action
			mapSimulate["action"] = []interface{}{action.Action+" "+action.Symbol}

			if !ok {
				step := Step{ Number: len(steps)}
				// Create a copy of mapSimulate
				step.State = make(map[string][]interface{})
				for key, value := range mapSimulate {
					step.State[key] = make([]interface{}, len(value))
					copy(step.State[key], value)
				}
				steps = append(steps, step)
				_=steps
				// Error
				return nil, fmt.Errorf("error en la tabla de análisis sintáctico, no existe acción para: "+action.Symbol)
			}
			step := Step{ Number: len(steps)}
			// Create a copy of mapSimulate
			step.State = make(map[string][]interface{})
			for key, value := range mapSimulate {
				step.State[key] = make([]interface{}, len(value))
				copy(step.State[key], value)
			}
			steps = append(steps, step)

			// Add to stack in the first
			mapSimulate["stack"] = append([]interface{}{action.Number}, mapSimulate["stack"]...)
		}
	}
	// Create a response
	response := SimResponse{Steps: steps, Accept: accept}

	// Convert List of Steps into a Json
	encoded, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error al convertir los pasos en json:", err)
		return nil, err
	}

	return encoded, nil
}


