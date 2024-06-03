package regex

import (
	"fmt"
	// Import the utils package
	"lenguagePr/utils"
)

// balanceExp function checks if the parentheses and square brackets are balanced in the regex
func balanceExp(regex string) ([]rune, bool) {
	// Initialize the boolean
	balanced := true
	// Convert the string to slice of runes
	runes := []rune(regex)
	// Create a result slice
	var result []rune

	isInBracket := 0
	inInSquareBracket := 0
	// Iterate over the runes
	for i := 0; i < len(runes); i++ {
		if runes[i] == '(' {
			result = append(result, runes[i])
			isInBracket++
		} else if runes[i] == ')' {
			if isInBracket > 0 {
				result = append(result, runes[i])
				// If stack is not empty, the expression is unbalanced
				isInBracket--

			} else {
				// If the stack is empty, the expression is unbalanced
				balanced = false
			}
		} else if runes[i] == '[' {
			inInSquareBracket++
			result = append(result, runes[i])
		} else if runes[i] == ']' {
			if inInSquareBracket > 0 {
				result = append(result, runes[i])
				// If stack is not empty, the expression is unbalanced
				inInSquareBracket--

			} else {
				// If the stack is empty, the expression is unbalanced
				balanced = false
				return result, balanced
			}
		} else if runes[i] == '\'' {
			// obtain next rune
			if i+1 < len(runes) {
				nextRune := runes[i+1]
				if nextRune == '\\' {
					if i+2 < len(runes) {
						result = append(result, runes[i+2])
						i += 3
						if i+3 < len(runes) {
							if runes[i+3] != '\'' {
								balanced = false
								return result, balanced
							}
						}
					} else {
						balanced = false
						return result, balanced
					}
				} else {
					result = append(result, nextRune)
					i += 2
				}
			}
		} else if runes[i] == '"' {
			// Iter next runes
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '\\' {
					// Add next rune
					if j+1 < len(runes) {
						result = append(result, runes[j+1])
						j++
					} else {
						balanced = false
						return result, balanced
					}

				} else if runes[j] == '"' {
					i = j
					break
				} else {
					result = append(result, runes[j])
				}
			}
		} else {
			if inInSquareBracket > 0 {
				if runes[i] != '-' && runes[i] != '^' {
					balanced = false
					return result, balanced
				}
				result = append(result, runes[i])
			} else {
				if contains([]string{"|", "*", "+", "?"}, string(runes[i])) {
					result = append(result, runes[i])
				} else if runes[i] == '\\' {
					result = append(result, runes[i])
				} else if runes[i] == '#' {
					// Verify in runes[i-1] is a ] and in runes[i+1] is a [
					if i-1 >= 0 && i+1 < len(runes) {
						if runes[i-1] == ']' && runes[i+1] == '[' {
							result = append(result, runes[i])
						} else {
							balanced = false
							return result, balanced
						}
					} else {
						balanced = false
						return result, balanced
					}
				} else {
					result = append(result, runes[i])
				}
			}
		}

	}

	return result, balanced
}

// makeSet function creates a set of runes from a range
func makeSet(start int, end int, runes []rune) []rune {
	// Verify if first element are ^
	negation:=false
	if runes[start] == '^' {
		negation = true
		start++
	}

	// Create a slice to store the result
	var result []rune
	// Iterate over the runes
	for i := start; i <= end; i++ {
		if runes[i] == '-' {
			// verify if has elements on result
			if len(result) > 0 {
				// Obtain last element
				last := result[len(result)-1]
				// Verify if has next rune
				if i+1 < len(runes) {
					next := runes[i+1]

					// Iterate over the runes, including the runes
					for j := last + 1; j <= next; j++ {
						// Verify if not repeated
						if !containsRune(result, j) {
							result = append(result, j)
						}
					}
				}
			}
		} else {
			// Verify if not repeated
			if !containsRune(result, runes[i]) {
				result = append(result, runes[i])
			}
		}
	}

	// Verify if has negation
	if negation {
		// Create a complete alphabet
		alphabet := [] rune{}
		for i := 0; i < 256; i++ {
			alphabet = append(alphabet, rune(i))
		}
		result = setDifference(alphabet, result)
	}
	return result
}

// setDifference function calculates the difference between two sets of runes
func setDifference(set1 []rune, set2 []rune) []rune {
	// Create a slice to store the result
	var result []rune
	// Iterate over the runes
	for i := 0; i < len(set1); i++ {
		if !containsRune(set2, set1[i]) {
			// Verify if not repeated
			if !containsRune(result, set1[i]) {
				result = append(result, set1[i])
			}
		}
	}
	return result
}

func FormatRegex(regexTex string) []interface{} {
	// Convert the string to a slice of chars
	runes, balanced := balanceExp(regexTex)
	if !balanced {
		fmt.Println("The expression is unbalanced")
		return nil
	} 
	// Create a slice to store the result
	var result []interface{}

	// Create a set of string to be used as operators
	operators := []string{"|", "(","[", "*", "."}

	// Iterate over the runes
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			// Verify if there is a next rune
			if i+1 < len(runes) {
				result = append(result, runes[i+1])
				// Skip the next rune
				i++
			}
		} else {
			if runes[i] == '(' {
				if len(result) > 0 {
					if !contains(operators, result[len(result)-1]) {
						// append the pipe to the result, to string
						result = append(result, ".")
					}
				}
				result = append(result, "(")
			} else if runes[i] == ')' {
				result = append(result, ")")
			} else if runes[i] == '[' {
				start:= i+1
				end := 0
				// Search for the closing bracket
				for j := i + 1; j < len(runes); j++ {
					if runes[j] == ']' {
						end = j-1
						i = j
						break
					}
				}

				// Obtain the set
				set := makeSet(start, end, runes)

				// Verify if the next rune is #
				if len(runes) > i+1 {
					if runes[i+1] == '#' {
						// Verify if the next next rune is [
						if len(runes) > i+2 {
							if runes[i+2] == '[' {
								start = i+3
								end = 0
								// Search for the closing bracket
								for j := i + 2; j < len(runes); j++ {
									if runes[j] == ']' {
										end = j-1
										i = j
										break
									}
								}
								// Obtain the set
								set2 := makeSet(start, end, runes)
								// Obtain the difference
								set = setDifference(set, set2)
								
							}
						}
					}
				}

				if len(result) > 0 {
					if !contains(operators, result[len(result)-1]) {
						// append the pipe to the result, to string
						result = append(result, ".")
					}
				}
				result = append(result, "(")
				// Iter over set to put values into result
				for j := 0; j < len(set); j++ {
					result = append(result, int32(set[j]))
					if j+1 < len(set) {
						result = append(result, "|")
					}
				}

				result = append(result, ")")


			} else{
				if runes[i] == '|' {
					// Append the pipe to the result, to string
					result = append(result, "|")
				} else if runes[i] == '*' {
					// append the pipe to the result, to string
					result = append(result, "*")
				} else if runes[i] == '+' {
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						if !contains(operators, result[len(result)-1]) {
							// Obtain the last element and pop
							last := result[len(result)-1]
							result = result[:len(result)-1]
							
							result = append(result, "(")
							result = append(result, last)
							result = append(result, ".")
							result = append(result, last)
							result = append(result, "*")
							result = append(result, ")")
						}
					}
				} else if runes[i] == '?' {
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						if !contains(operators, result[len(result)-1]) {
							// Obtain the last element and pop
							last := result[len(result)-1]
							result = result[:len(result)-1]
							
							result = append(result, "(")
							result = append(result, last)
							result = append(result, "|")
							result = append(result, "epsilon")
							result = append(result, ")")
						}
					}
				} else {
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						if !contains(operators, result[len(result)-1]) {
							// append the pipe to the result, to string
							result = append(result, ".")
						}
					}
					result = append(result, int32(runes[i]))
				}
			}
		}
	}

	return result
}
