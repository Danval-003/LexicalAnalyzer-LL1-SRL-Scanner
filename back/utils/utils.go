package utils

// Function to verify if a string is in a slice of strings
func contains(s []string, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Function to verify if a rune is in a slice of runes
func containsRune(s []rune, e rune) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

