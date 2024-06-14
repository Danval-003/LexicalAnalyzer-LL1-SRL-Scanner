package scanners

import (
	"backend/routes/login"
	"backend/src/afd"
	"backend/src/yalex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Struct to represent Yalex Priv
type YalexPriv struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}

// Struct to represent Yalex Public
type YalexPublic struct {
	Content string `json:"content"`
}

// Struct to represent a scanner Response
type ScannerResponse struct {
	Message string      `json:"message" example:"Scanner created successfully"`
	Status  int         `json:"status" example:"200"`
	Names    []string      `json:"names"`
}

// Struct to represent a Bad Request Response
type BadRequestResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status" example:"400"`
}

// Struct to represent a scanner Save
type ScannerSave struct {
	Scanner string `json:"scanner"`
	Name  string   `json:"name"`
	Username string `json:"username"`
	ImageURL string `json:"imageURL"`
}


// Function to create a Private scanner from a Yalex
// @Summary Create a Private scanner from a Yalex
// @Description Create a Private scanner from a Yalex
// @Tags yalex
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param yalex body YalexPublic true "Yalex object"
// @Success 200 {object} ScannerResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/scanners/priv/create [post]
func CreatePrivateScanner(w http.ResponseWriter, r *http.Request) {
	var yal YalexPriv
	err := json.NewDecoder(r.Body).Decode(&yal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Obtain user
	user, ok := r.Context().Value(login.UserContextKey).(string)
	if !ok {
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error getting user",
				Status:  http.StatusInternalServerError,
			},
		)
	}


	// Create the scanners
	scanners, err := yalex.Yal(yal.Content)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			})
		return
	}

	// Save the scanners on mongo
	collection:= login.Client.Database("GOLL1").Collection("scanners")

	names := []string{}

	for key, value := range scanners {
		_, err = collection.InsertOne(
			r.Context(),
			ScannerSave{
				Name: key,
				Scanner: value["Machine"],
				Username: user,
				ImageURL: "https://www.google.com",
			},
		)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				BadRequestResponse{
					Message: err.Error(),
					Status:  http.StatusInternalServerError,
				},
			)
			return
		}
		names = append(names, key)
	}

	json.NewEncoder(w).Encode(
		ScannerResponse{
			Message: "Scanner created successfully",
			Status:  http.StatusOK,
			Names:   names,
		},
	)

}


// Function to create a Public scanner from a Yalex
// @Summary Create a Public scanner from a Yalex
// @Description Create a Public scanner from a Yalex
// @Tags yalex
// @Accept json
// @Produce json
// @Param yalex body YalexPublic true "Yalex object"
// @Success 200 {object} ScannerResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/scanners/public/create [post]
func CreatePublicScanner(w http.ResponseWriter, r *http.Request) {
	var yal YalexPublic
	err := json.NewDecoder(r.Body).Decode(&yal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}
	// Create a Unique ID
	id := uuid.New()

	// Create the scanners
	scanners, err := yalex.Yal(yal.Content)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			})
		return
	}

	// Save the scanners on mongo
	collection:= login.Client.Database("GOLL1").Collection("scanners")

	names := []string{}

	for key, value := range scanners {
		_, err = collection.InsertOne(
			r.Context(),
			ScannerSave{
				Name: id.String()+":"+key,
				Scanner: value["Machine"],
				Username: "public",
				ImageURL: "https://www.google.com",
			},
		)
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				BadRequestResponse{
					Message: err.Error(),
					Status:  http.StatusInternalServerError,
				},
			)
			return
		}
		names = append(names, id.String()+":"+key)
	}

	json.NewEncoder(w).Encode(
		ScannerResponse{
			Message: "Scanner created successfully",
			Status:  http.StatusOK,
			Names:   names,
		},
	)

}

// Struct to represent a scannerrequest
type ScannerRequest struct {
	Content string `json:"content"`
	ScannerName string `json:"scannerName"`
}

// Struct to Represent Scanner Sim
type SimScan struct {
	Message string `json:"message"`
	Status  int    `json:"status" example:"200"`
	Name   string `json:"name"`
	SimPart afd.ResponseSim `json:"simPart"`
}


// Function to Simulate a scanner public
// @Summary Simulate a scanner
// @Description Simulate a scanner
// @Tags yalex
// @Accept json
// @Produce json
// @Param yalex body ScannerRequest true "Scanner object"
// @Success 200 {object} SimScan
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/scanners/public/simulate [post]
func SimulatePublicScanner(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Simulating scanner")
	var scan ScannerRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Verify if scanner exist
	collection:= login.Client.Database("GOLL1").Collection("scanners")
	var scanner ScannerSave
	err = collection.FindOne(r.Context(), bson.D{
		{Key: "name", Value: scan.ScannerName},
		{Key: "username", Value: "public"},
	}).Decode(&scanner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			})
		return
	}

	fmt.Println(" ", )

	scannerM := afd.DecodeAFD(scanner.Scanner)


	// Simulate the scanner
	simulate := afd.ResponseSimulateAFD(scannerM, scan.Content)



	json.NewEncoder(w).Encode(
		SimScan{
			Message: "Scanner created successfully",
			Status:  http.StatusOK,
			Name:   scanner.Name,
			SimPart: simulate,
		},
	)

}

// Function to Simulate Private Scanner
// @Summary Simulate a scanner
// @Description Simulate a scanner
// @Tags yalex
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param yalex body ScannerRequest true "Scanner object"
// @Success 200 {object} SimScan
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/scanners/priv/simulate [post]
func SimulatePrivateScanner(w http.ResponseWriter, r *http.Request) {
	var scan ScannerRequest
	err := json.NewDecoder(r.Body).Decode(&scan)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Obtain user
	user, ok := r.Context().Value(login.UserContextKey).(string)
	if !ok {
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error getting user",
				Status:  http.StatusInternalServerError,
			},
		)
	}

	// Verify if scanner exist
	collection:= login.Client.Database("GOLL1").Collection("scanners")
	var scanner ScannerSave
	err = collection.FindOne(r.Context(), bson.D{
		{Key: "name", Value: scan.ScannerName},
		{Key: "username", Value: user},
	}).Decode(&scanner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			})
		return
	}

	scannerM := afd.DecodeAFD(scanner.Scanner)

	// Simulate the scanner
	simulate := afd.ResponseSimulateAFD(scannerM, scan.Content)

	json.NewEncoder(w).Encode(
		SimScan{
			Message: "Scanner created successfully",
			Status:  http.StatusOK,
			Name:   scanner.Name,
			SimPart: simulate,
		},
	)
}





