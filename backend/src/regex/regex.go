// regex.go
package regex

/*
Daniel Armando Valdez Reyes | Danval-003
Description: 
This code is a package that contains a function to format a regex string into a list of runes. 
*/

import (
	"backend/src/regex/regexFormated"
)


// Function to pass infix a postfix
func InfixToPostfix(regex string) []interface{} {
	infix := regexFormated.FormatRegex(regex)
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
		} else if regexFormated.Contains(operators, infix[i]) {
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
