package yaparroutes

import (
	"backend/routes/login"
	"backend/src/yapar"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/google/uuid"
	//"github.com/gorilla/mux"
)

// Struct to represent yalp
type Yalp struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}

// Struct to represent a yalp Public
type YalpPublic struct {
	Content string `json:"content"`
}

// Struct to represent a table Response
type TableResponse struct {
	Table   yapar.SLR `json:"table"`
	Message string      `json:"message" example:"Table created successfully"`
	Status  int         `json:"status" example:"200"`
	ImageURL string     `json:"imageURL"`
	Name string `json:"name"`
}

// Struct to represent a Bad Request Response
type BadRequestResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status" example:"400"`
}

// TableSave struct represents the structure of the table to save
type TableSave struct {
	Table []byte `json:"table"`
	Name  string   `json:"name"`
	Username string `json:"username"`
	ImageURL string `json:"imageURL"`
}



// @Summary Create a Private table SRL from a Yalp
// @Description Create a Private table SRL from a Yalp
// @Tags yalp
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param yalp body Yalp true "Yalp object"
// @Success 200 {object} TableResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/yapar/priv/create [post]
func CreatePrivateTable(w http.ResponseWriter, r *http.Request) {
	var yalp Yalp
	_ = json.NewDecoder(r.Body).Decode(&yalp)

	urlImage, table, err := yapar.ObtainSrl(yalp.Content, yalp.Name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error creating table",
				Status:  http.StatusInternalServerError,
			},
		)
	}


	// Extract username from context
	username, ok := r.Context().Value(login.UserContextKey).(string)
	if !ok || username == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving user from context",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Conect to the database
	collection := login.Client.Database("GOLL1").Collection("SLRs")

	fmt.Println("Creating table")

	tableBinary, err := json.Marshal(table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error creating table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Verify if the table already exists
	count, err := collection.CountDocuments(
		r.Context(),
		bson.D{{Key: "name", Value: yalp.Name}, {Key: "username",
			Value: username}},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error verifying table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Create struct to save table
	tableSave := TableSave{tableBinary, yalp.Name, username, urlImage}

	if count > 0 {
		// Update table and image
		_, err = collection.UpdateOne(r.Context(), 
		bson.D{{Key: "name", Value: yalp.Name}, {Key: "username",
			Value: username}},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "table", Value: tableSave.Table},
				{Key: "imageURL", Value: tableSave.ImageURL},
			}},
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				BadRequestResponse{
					Message: "Error updating table",
					Status:  http.StatusInternalServerError,
				},
			)
			return
		}

	} else {

		// Save the table
		_, err = collection.InsertOne(r.Context(), 
		bson.D{
			{Key: "table", Value: tableSave.Table},
			{Key: "name", Value: tableSave.Name},
			{Key: "username", Value: tableSave.Username},
			{Key: "imageURL", Value: tableSave.ImageURL},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(
				BadRequestResponse{
					Message: "Error saving table",
					Status:  http.StatusInternalServerError,
				},
			)
			return
		}

	}

	

	json.NewEncoder(w).Encode(
		TableResponse{
			Table:   table,
			Message: "Table created successfully",
			Status:  http.StatusOK,
			ImageURL: urlImage,
			Name: yalp.Name,
		},
	)

}


// Function to create a Public Table
// @Summary Create a Public table SRL from a Yalp
// @Description Create a Public table SRL from a Yalp
// @Tags yalp
// @Accept json
// @Produce json
// @Param yalp body YalpPublic true "Yalp object"
// @Success 200 {object} TableResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/yapar/pub/create [post]
func CreatePublicTable(w http.ResponseWriter, r *http.Request) {
	// Obtain Data
	var yalp YalpPublic
	_ = json.NewDecoder(r.Body).Decode(&yalp)

	// Generate a random and unique id
	id := uuid.New().String()

	// Create the table
	urlImage, table, err := yapar.ObtainSrl(yalp.Content, "Public Table")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error creating table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Conect to the database
	collection := login.Client.Database("GOLL1").Collection("SLRs")

	// Save the table
	tableBinary, err := json.Marshal(table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error creating table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Save the table
	_, err = collection.InsertOne(r.Context(), 
	bson.D{
		{Key: "table", Value: tableBinary},
		{Key: "name", Value: id},
		{Key: "username", Value: "public"},
		{Key: "imageURL", Value: urlImage},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error saving table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	json.NewEncoder(w).Encode(
		TableResponse{
			Table:   table,
			Message: "Table created successfully",
			Status:  http.StatusOK,
			ImageURL: urlImage,
			Name: id,
		},
	)
}






// Struct to represent a Table Request
type TableRequest struct {
	Name string `json:"name"`
}


// @Summary Get a table SRL from user
// @Description Get a table SRL from user
// @Tags yalp
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name body TableRequest true "Table object"
// @Success 200 {object} TableResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/yapar/priv/get [post]
func GetPrivateTable(w http.ResponseWriter, r *http.Request) {
	var tableRequest TableRequest
	_ = json.NewDecoder(r.Body).Decode(&tableRequest)

	// Extract username from context
	username, ok := r.Context().Value(login.UserContextKey).(string)
	fmt.Println(username)
	if !ok || username == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving user from context",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	// Conect to the database
	collection := login.Client.Database("GOLL1").Collection("SLRs")

	// Find the table
	var table TableSave
	err := collection.FindOne(r.Context(), bson.D{{Key: "name", Value: tableRequest.Name}, {Key: "username", Value: username}}).Decode(&table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	var tableSRL yapar.SLR
	err = json.Unmarshal(table.Table, &tableSRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	json.NewEncoder(w).Encode(
		TableResponse{
			Table:   tableSRL,
			Message: "Table retrieved successfully",
			Status:  http.StatusOK,
			ImageURL: table.ImageURL,
			Name: tableRequest.Name,
		},
	)
}


// Function to get Public Table
// @Summary Get a table SRL from public
// @Description Get a table SRL from public
// @Tags yalp
// @Accept json
// @Produce json
// @Param name body TableRequest true "Table object"
// @Success 200 {object} TableResponse
// @Failure 400 {object} BadRequestResponse
// @Router /api/v1/yapar/pub/get [post]
func GetPublicTable(w http.ResponseWriter, r *http.Request) {
	var tableRequest TableRequest
	_ = json.NewDecoder(r.Body).Decode(&tableRequest)

	// Conect to the database
	collection := login.Client.Database("GOLL1").Collection("SLRs")

	// Find the table
	var table TableSave
	err := collection.FindOne(r.Context(), bson.D{{Key: "name", Value: tableRequest.Name}, {Key: "username", Value: "public"}}).Decode(&table)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	var tableSRL yapar.SLR
	err = json.Unmarshal(table.Table, &tableSRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			BadRequestResponse{
				Message: "Error retrieving table",
				Status:  http.StatusInternalServerError,
			},
		)
		return
	}

	json.NewEncoder(w).Encode(
		TableResponse{
			Table:   tableSRL,
			Message: "Table retrieved successfully",
			Status:  http.StatusOK,
			ImageURL: table.ImageURL,
			Name: tableRequest.Name,
		},
	)
}




