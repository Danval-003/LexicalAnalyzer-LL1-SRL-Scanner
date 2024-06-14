package compare

import (
	"backend/routes/login"
	"backend/routes/scanners"
	"backend/routes/yaparRoutes"
	"backend/src/afd"
	"backend/src/yapar"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Struct to represent part to simulate
type Simulate struct {
	ScannerName string `json:"scannerName"`
	ContentSimulate string `json:"contentSimulate"`
	SLRName string `json:"slrName"`
}

// Struct to represent a simResponse
type SimResponse struct {
	Message string `json:"message"`
	Status int `json:"status" example:"200"`
	Sim yapar.SimResponse `json:"sim"`
	ScannerResult string `json:"scannerResult"`
}

// Struct to represent a BadResponse
type BadRequestResponse struct {
	Message string `json:"message"`
	Status int `json:"status" example:"400"`
}


// Func to simulate compile With public Resources
// @Summary Simulate compiler with public resources
// @Description Simulate compiler
// @Tags compare
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param simulate body Simulate true "Simulate object"
// @Success 200 {object} SimResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/compare/simulate [post]
func SimulateCompile(w http.ResponseWriter, r *http.Request) {
	init:= time.Now()
	// Obtain Data
	var simulate Simulate
	err := json.NewDecoder(r.Body).Decode(&simulate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "Error in the data", Status: http.StatusBadRequest})
		return
	}

	// Search Scanner and SLR
	collection:=login.Client.Database("GOLL1").Collection("scanners")
	scanner := scanners.ScannerSave{}
	err = collection.FindOne(r.Context(), bson.D{
		{Key: "name", Value: simulate.ScannerName},
		{Key: "username", Value: "public"},
	}).Decode(&scanner)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "Scanner not found", Status: http.StatusBadRequest})
		return
	}

	// Decoded Afd
	scannerM := afd.DecodeAFD(scanner.Scanner)

	// Simulate the scanner
	simScan := afd.ResponseSimulateAFD(scannerM, simulate.ContentSimulate)

	// Verify if simulation is accepted
	if !simScan.Accepted {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "Simulation not accepted", Status: http.StatusBadRequest})
		return
	}

	// Search SLR
	collection = login.Client.Database("GOLL1").Collection("SLRs")
	slr := yaparroutes.TableSave{}
	err = collection.FindOne(r.Context(), bson.D{ 
		{Key: "name", Value: simulate.SLRName},
		{Key: "username", Value: "public"},
	}).Decode(&slr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "SLR not found", Status: http.StatusBadRequest})
		return
	}

	var slrDecod yapar.SLR

	err = json.Unmarshal(slr.Table, &slrDecod)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "Error decoding SLR", Status: http.StatusBadRequest})
		return
	}

	// Simulate SLR with scanner simulation summary
	simResponse, err := yapar.SimulateTable(slrDecod, simScan.StringSummary)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "SLR: "+err.Error(), Status: http.StatusBadRequest})
		return
	}

	var simResp yapar.SimResponse
	err = json.Unmarshal(simResponse, &simResp)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BadRequestResponse{Message: "Error decoding Simulate SLR", Status: http.StatusBadRequest})
		return
	}

	fmt.Println("Time to simulate: ", time.Since(init))
	// Response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		SimResponse{
			Message: "Simulation completed", 
			Status: http.StatusOK, 
			Sim: simResp, 
			ScannerResult: simScan.StringSummary})


}



