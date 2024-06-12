package scanners

import (
	"backend/src/afd"
	"fmt"
)

// Function to search a scanner Machine into the list of scanners
func SearchScanner(scannerName string) (afd.State, bool) {
	path := "./src/Machines/Scaners/" + scannerName + ".json"
	if afd.MachineExists(path) {
		return *afd.LoadMachine(path), true
	} else {
		return afd.State{}, false
	}
}

// Function to simulate a scanner Machine
func SimulateScanner(scannerName string, input string) (string, error) {
	scanner, found := SearchScanner(scannerName)
	if found {
		// Simulated scanner and encoded using encoding function to SimulatedParts
		simulated:= afd.SimulateAFD(&scanner, input)

		// Encoded using encoding function to SimulatedParts
		encoded, err:= simulated.Encode()
		if err != nil {
			return "", err
		}

		return encoded, nil

	} else {
		return "", fmt.Errorf("scanner not found")
	}

}

