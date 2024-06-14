package regexFormated

import (
	"fmt"
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

// Function to verify if a rune is in a slice of runes
func ContainsRune(s []rune, e rune) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


// BalanceExp function checks if the parentheses and square brackets are balanced in the regex
func BalanceExp(regex string) ([]rune, bool) {
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
		if runes[i] == '\\' {
			if i+1 < len(runes) {
				result = append(result, runes[i])
				result = append(result, runes[i+1])
				i++
			}
			continue
		}


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
					if Contains([]string{"|", "*", "+", "?", "_", "(", ")", "[", "]"}, string(nextRune)) {
						result = append(result, '\\')
						result = append(result, nextRune)
					} else {
						result = append(result, nextRune)
					}
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
					if Contains([]string{"|", "*", "+", "?", "_", "(", ")", "[", "]"}, string(runes[j])) {
						result = append(result, '\\')
						result = append(result, runes[j])
					} else {
						result = append(result, runes[j])
					}
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
				if Contains([]string{"|", "*", "+", "?", "_"}, string(runes[i])) {
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
					balanced = false
					return result, balanced
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
						if !ContainsRune(result, j) {
							result = append(result, j)
						}
					}
				}
			}
		} else {
			// Verify if not repeated
			if !ContainsRune(result, runes[i]) {
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
		if !ContainsRune(set2, set1[i]) {
			// Verify if not repeated
			if !ContainsRune(result, set1[i]) {
				result = append(result, set1[i])
			}
		}
	}
	return result
}


func FormatRegex(regexTex string) []interface{} {
	// Convert the string to a slice of chars
	runes, balanced := BalanceExp(regexTex)
	if !balanced {
		fmt.Println("The expression is unbalanced")
		return nil 
	} 
	// Create a slice to store the result
	var result []interface{}

	// Create a set of string to be used as operators
	operators := []string{"|", "(","[", "."}
	end_index :=1 
	var tempIndex [] int
	tempIndex = append(tempIndex, 0)

	// Iterate over the runes
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' {
			if len(result) > 0 {
				// Verify if type from last element in result is String
				last := result[len(result)-1]
				if _, ok := last.(string); ok {
					if !Contains(operators, result[len(result)-1]) {
						// append the pipe to the result, to string
						result = append(result, ".")
					}
				} else {
					result = append(result, ".")
				}
			}
			// Verify if there is a next rune
			if i+1 < len(runes) {
				result = append(result, runes[i+1])
				// Skip the next rune
				i++
				tempIndex[len(tempIndex)-1] = len(result)-1
				end_index = len(result)
			}
			continue
		} else {
			if runes[i] == '(' {
				if len(result) > 0 {
					// Verify if type from last element in result is String
					last := result[len(result)-1]
					if _, ok := last.(string); ok {
						if !Contains(operators, result[len(result)-1]) {
							// append the pipe to the result, to string
							result = append(result, ".")
						}
					} else {
						result = append(result, ".")
					}
				}
				result = append(result, "(")
				tempIndex[len(tempIndex)-1] = len(result)-1
				tempIndex = append(tempIndex, len(result)-1)
			} else if runes[i] == ')' {
				result = append(result, ")")
				end_index = len(result)
				tempIndex = tempIndex[:len(tempIndex)-1]
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
					// Verify if type from last element in result is String
					last := result[len(result)-1]
					if _, ok := last.(string); ok {
						if !Contains(operators, result[len(result)-1]) {
							// append the pipe to the result, to string
							result = append(result, ".")
						}
					} else {
						result = append(result, ".")
					}
				}
				result = append(result, "(")
				tempIndex = append(tempIndex, len(result)-1)
				// Iter over set to put values into result
				for j := 0; j < len(set); j++ {
					result = append(result, int32(set[j]))
					if j+1 < len(set) {
						result = append(result, "|")
					}
				}

				result = append(result, ")")
				end_index = len(result)

			} else{
				if runes[i] == '|' {
					// Copy the last elements into last
					var last []interface{}
					for j := tempIndex[len(tempIndex)-1]; j < end_index; j++ {
						last = append(last, result[j])
					}
					result = result[:tempIndex[len(tempIndex)-1]]
					// Append the pipe to the result, to string
					result = append(result, "(")
					tempIndex[len(tempIndex)-1] = len(result)-1
					result = append(result, last...)
					result = append(result, ")")
					result = append(result, "|")
				} else if runes[i] == '*' {
					// Copy the last elements into last
					var last []interface{}
					for j := tempIndex[len(tempIndex)-1]; j < end_index; j++ {
						last = append(last, result[j])
					}
					result = append(result, "(")
					result = result[:tempIndex[len(tempIndex)-1]]
					// Append the pipe to the result, to string
					result = append(result, "(")
					tempIndex[len(tempIndex)-1] = len(result)-1
					result = append(result, last...)
					result = append(result, ")")
					result = append(result, "*")
					result = append(result, ")")
					end_index = len(result)
				} else if runes[i] == '+' {
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						if !Contains([]string{"|", "(","[", ".", "*"}, result[len(result)-1]) {
							// Copy the last elements into last
							var last []interface{}
							for j := tempIndex[len(tempIndex)-1]; j < end_index; j++ {
								last = append(last, result[j])
							}
							
							result = result[:tempIndex[len(tempIndex)-1]]

							text:=""
							for j:=0; j<len(last); j++ {
								if _, ok := last[j].(int32); ok {
									text += string(rune(last[j].(int32)))
								} else {
									text += last[j].(string)
								}
							}
							
							
							
							result = append(result, "(")
							
							
							tempIndex[len(tempIndex)-1] = len(result)-1
							result = append(result, "(")

							result = append(result, last...)

							result = append(result, ")")
							

							result = append(result, ".")
							result = append(result, "(")
							result = append(result, "(")
							result = append(result, last...)
							result = append(result, ")")
							result = append(result, "*")
							result = append(result, ")")
							result = append(result, ")")
							end_index = len(result)
						}
					}
				} else if runes[i] == '?' {
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						if !Contains([]string{"|", "(","[", ".", "*"}, result[len(result)-1]) {
							// Copy the last elements into last
							var last []interface{}
							for j := tempIndex[len(tempIndex)-1]; j < end_index; j++ {
								last = append(last, result[j])
							}
							result = result[:tempIndex[len(tempIndex)-1]]
							
							result = append(result, "(")
							tempIndex[len(tempIndex)-1] = len(result)-1
							result = append(result, "(")
							result = append(result, last...)
							result = append(result, ")")
							result = append(result, "|")
							result = append(result, "epsilon")
							result = append(result, ")")
							end_index = len(result)
						}
					}
				} else if runes[i] == '_' { 
					// Verify if the last rune is not a operator, (not in the set operators)
					if len(result) > 0 {
						// Verify if type from last element in result is String
						last := result[len(result)-1]
						if _, ok := last.(string); ok {
							if !Contains(operators, result[len(result)-1]) {
								// append the pipe to the result, to string
								result = append(result, ".")
							}
						} else {
							result = append(result, ".")
						}
					}

					result = append(result, "(")
					// Add All the runes from 0 to 255
					for j := 0; j < 256; j++ {
						result = append(result, int32(j))
						if j+1 < 256 {
							result = append(result, "|")
						}
					}
					result = append(result, ")")
				
				}else {
					
					if len(result) > 0 {
						// Verify if type from last element in result is String
						last := result[len(result)-1]
						if _, ok := last.(string); ok {
							if !Contains(operators, result[len(result)-1]) {
								// append the pipe to the result, to string
								result = append(result, ".")
							}
						} else {
							result = append(result, ".")
						}
					}
					result = append(result, int32(runes[i]))
					if len(tempIndex) >0 {
						if len(result)>1{
							if result[len(result)-1] == "*" {
								tempIndex[len(tempIndex)-1] = len(result)-2
								end_index = len(result)
								continue
							}
						}
	
						tempIndex[len(tempIndex)-1] = len(result)-1
						end_index = len(result)
					}
				}

			}
		}
	}


	return result
}
