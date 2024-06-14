package yalex

import (
	"backend/src/afd"
	"fmt"
	"strings"
)

// Function to manage declarations
func Declarations(declarations string, declares *map[string]string) error {
	// Split using '='
	dec := strings.Split(declarations, "=")
	// Check if the length is not 2
	if len(dec) != 2 {
		// Return an error
		return fmt.Errorf("error in declaration: %s", declarations)
	}

	// Create the machine
	DeclareTokens := map[string]string{
		"STRING": "\\\"([^'\"']|\\\\\\\")+\\\"",
		"SYMBOL": "\\'([^\\']|\\')\\'",
		"WS":"[' ''\n''\t']+",
	}
	// Add into machine keys to declares
	for key := range *declares {
		DeclareTokens[key] = "\""+key+"\""
	}
	machine,_, _ := afd.MakeAFD(DeclareTokens)

	// Define the value of declaration
	value:= strings.TrimSpace(dec[1])

	// Remplace literal use to \n to '\n'
	value = strings.ReplaceAll(value, "\\n", "\n")
	// Remplace literal use to \t to '\t'
	value = strings.ReplaceAll(value, "\\t", "\t")

	//afd.VisualizeAFD(DeclareMachine, "DeclareMachine", "DeclareMachine")
	simulate:= afd.SimulateAFD(machine, value)

	formatValue := ""

	// Eval tokens
	for _, token := range simulate {
		if token.Accepted {
			if token.Token != "STRING" && token.Token != "SYMBOL" && token.Token != "WS" {
				// Add the value into declares
				formatValue += (*declares)[token.Token]
			} else {
				formatValue += string(token.Runes)
			}
		} else {
			formatValue += string(token.Runes)
		}
	}
	if formatValue == "" {
		formatValue = value
	}


	// Save the declaration without word "let"
	(*declares)[strings.TrimSpace(dec[0][3:])] = formatValue


	// Return nil
	return nil
}

// Function to manage rules
func TokenRules(rules string, rule *map[string]string, declares map[string]string) error {
	// Split definition to part to return if return exists using a machine
	var machine *afd.State
	// Check if the machine exists
	if afd.MachineExists("./src/Machines/Readers/RuleMachine.json") {
		// Load the machine
		machine = afd.LoadMachine("./src/Machines/Readers/RuleMachine.json")
	} else {
		// Create the machine
		RuleTokens := map[string]string{
			"DEFINITION": "([^' ''\n''\t']|(\\'[' ''\n''\t']\\'))+",
			"RETURN": "'{'' '*\"return\"' '+['A'-'Z']+' '*'}'",
			"WS":"[' ''\n''\t']+",
		}
		var States afd.StateSlice
		machine,_, States = afd.MakeAFD(RuleTokens)

		// Save the machine
		afd.SaveMachine("./src/Machines/Readers/RuleMachine.json", States)
	}

	var definition string
	var returnRule string
	// Iterate over the rules
	for _, token := range afd.SimulateAFD(machine, rules) {
		if token.Accepted {
			if token.Token == "DEFINITION" {
				definition = string(token.Runes)
				// Quit spaces
				definition = strings.TrimSpace(definition)
			} else if token.Token == "RETURN" {
				returnRule = string(token.Runes)
				// Quit {}
				returnRule = returnRule[1:len(returnRule)-1]
				// Quit return
				returnRule = returnRule[7:]
				// Quit spaces
				returnRule = strings.TrimSpace(returnRule)

			}
		} else {
			return fmt.Errorf("token: %s not accepted", token.Token)
		}
	}

	dec := map[string]string{}

	// ItersOverDeclares
	for key := range declares {
		dec[key] = "\""+key+"\""
	}

	// Create a machine with declares
	machineDeclare, _, _ := afd.MakeAFD(dec)

	// Simulate the machine using the definition
	simulate := afd.SimulateAFD(machineDeclare, definition)

	// Format the definition
	formatDefinition := ""
	for _, token := range simulate {
		if token.Accepted {
			if token.Token != "STRING" && token.Token != "SYMBOL" && token.Token != "WS" {
				// Add the value into declares
				formatDefinition += declares[token.Token]
			} else {
				formatDefinition += string(token.Runes)
			}
		} else {
			formatDefinition += string(token.Runes)
		}
	}
	if formatDefinition == "" {
		formatDefinition = definition
	}

	if returnRule == "" {
		returnRule = definition
		// Delete {}
		returnRule = returnRule[1:len(returnRule)-1]
		// Delete return
		returnRule = returnRule[7:]
		// Quit spaces
		returnRule = strings.TrimSpace(returnRule)
	}

	// Save the rule
	(*rule)[returnRule] = formatDefinition

	// Return nil
	return nil
}


// Yal function to read a yal file
func Yal(yal string) (map[string]map[string]string,error) {


	// Make a var to save error
	var err error
	var machine *afd.State
	// Check if the machine exists
	if afd.MachineExists("./src/Machines/Readers/YalexMachine.json") {
		// Load the machine
		machine = afd.LoadMachine("./src/Machines/Readers/YalexMachine.json")
	} else {
		// Create the machine
		YalexTokens := map[string]string{
			"COMMET": "\\(\\*[^'*']+\\*\\)",
			"DECLARE": "\"let\"' '+['a'-'z']+' '*'='' '*([^'\n']|\\''\n'\\')+",
			"RULES": "\"rule\"' '+['a'-'z']+' '*'='[' ''\n''\t']*([^' ''\n''\t']|(\\'[' ''\n''\t']\\'))([^' ''\n''\t']|(\\'[' ''\n''\t']\\'))*(' '*\\{' '*\"return\"' '+['A'-'Z']+' '*\\})?",
			"EXTRARULE": "\\|' '*([^' ''\n''\t']|(\\'[' ''\n''\t']\\'))([^' ''\n''\t']|(\\'[' ''\n''\t']\\'))*(' '*\\{' '*\"return\"' '+['A'-'Z']+' '*\\})?",
			"WS":"[' ''\n''\t']+",
		}
		var States afd.StateSlice
		machine,_, States = afd.MakeAFD(YalexTokens)
		afd.VisualizeAFD(machine, "YalexMachine", "YalexMachine")

		// Save the machine
		afd.SaveMachine("./src/Machines/Readers/YalexMachine.json", States)
	}

	//afd.VisualizeAFD(YalexMachine, "YalexMachine", "YalexMachine")
	simulate:= afd.SimulateAFD(machine, yal)
	declares := make(map[string]string)

	actualRule := ""

	Rules := map[string]map[string]string{}

	// Iterate over the result
	for _, token := range simulate {
		
		if token.Accepted {

			if token.Token == "DECLARE" {
				err = Declarations(string(token.Runes), &declares)
				if err != nil {
					return nil,err
				}
			} else if token.Token == "RULES" {
				// Split using '='
				rules := strings.Split(string(token.Runes), "=")
				// Check if the length is not 2
				if len(rules) != 2 {
					// Return an error
					return nil, fmt.Errorf("error in rules: %s", string(token.Runes))
				}
				// Save the actual rule
				actualRule = strings.TrimSpace(rules[0][5:])
				// Create a map to save the rules
				Rules[actualRule] = make(map[string]string)
				// Add the rule
				rule := Rules[actualRule]
				err = TokenRules(rules[1], &rule, declares)
				Rules[actualRule] = rule

			} else if token.Token == "EXTRARULE" {
				// Delete '|' in the first position and quit spaces
				extraRule := strings.TrimSpace(string(token.Runes[1:]))
				
				// Verify if the actual rule exists
				if _, ok := Rules[actualRule]; !ok {
					return nil, fmt.Errorf("error the rule: %s not exists", actualRule)
				}

				// Add the rule in actual rule
				rule := Rules[actualRule]

				err = TokenRules(extraRule, &rule, declares)
				Rules[actualRule] = rule
			}
		} else { 
			return nil, fmt.Errorf("token: %s not accepted", token.Token)
		}
	}

	// Create a map to save the machines
	machines := make(map[string]map[string]string)


	// Iterate over the rules
	for key, value := range Rules {
		machines[key] = map[string]string{}
		// Create the machine
		init,_, states := afd.MakeAFD(value)
		afd.SimulateAFD(init, " ")

		imageUrll:= "https://quickchart.io/graphviz?graph="

		bytesBuffer := states.Encode()
		// Save the machine
		machines[key]["Machine"] = bytesBuffer

		// Save the machine url
		machines[key]["Image"] = imageUrll
		
	}

	return machines, err
}


