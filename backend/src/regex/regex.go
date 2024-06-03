// regex.go
package regex

/*
Daniel Armando Valdez Reyes | Danval-003
Description: 
This code is a package that contains a function to format a regex string into a list of runes. 
*/

import (
	"fmt"
	"github.com/Danval-003/LexicalAnalyzer-LL1-SRL-Scanner/backend/src/regex/regexFormated"
)

// Function to verify if a string is in a slice of strings
func Contains(s []string, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Function to pass infix a postfix
func InfixToPostfix(regex string) []interface{} {
	infix := regexFormated.FormatRegex(regex)
	fmt.Println(infix)
	// Create a slice to store the result
	result := []interface{}{}
	// Create a stack to store the operators
	stack := []interface{}{}
	// Create a set of string to be used as operators
	operators := []string{"|", "(", "[", "*", "."}
	// Create a map to operators with precedence
	precedence := map[string]int{
		"(": 1, ")": 1, "[": 1, "]": 1, "|": 2, ".": 3, "*": 4,
	}

	// Iterate over the infix
	for i := 0; i < len(infix); i++ {
		if infix[i] == "(" {
			stack = append(stack, infix[i])
		} else if infix[i] == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		} else if Contains(operators, infix[i]) {
			for len(stack) > 0 && precedence[stack[len(stack)-1].(string)] >= precedence[infix[i].(string)] {
				result = append(result, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, infix[i])
		} else {
			result = append(result, infix[i])
		}
	}

	for len(stack) > 0 {
		result = append(result, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return result
}
