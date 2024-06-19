package main

import (
	_ "backend/docs" // This is to import generated docs
	"backend/routes/login"
	yaparroutes "backend/routes/yaparRoutes"
	scanners "backend/routes/scanners"
	compare "backend/routes/compare"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers" // Import the missing package
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

)


// @title GO-Api API
// @version 1.0
// @description This is a Api to create a language ll1
// @termsOfService http://swagger.io/terms/
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme bearer
// @contact.name Daniel Valdez
// @contact.email danarvare@outlook.com
// @BasePath /api/v1
func main() {

	// Create var to store the error
	var err error

	// Load the .env file
	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
		// Finish the program
		return
	}

	err = login.CreateClientLogin()
	if err != nil {
		fmt.Println(err)
		// Finish the program
		return
	}

	r := mux.NewRouter()

	
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	// Define the routes and their handlers
	// Define the API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/register", login.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/login", login.Login).Methods(http.MethodPost)
	api.Handle("/yapar/priv/create", login.IsAuthorized(yaparroutes.CreatePrivateTable)).Methods(http.MethodPost)
	api.Handle("/yapar/priv/get", login.IsAuthorized(yaparroutes.GetPrivateTable)).Methods(http.MethodPost)
	api.HandleFunc("/yapar/pub/create", yaparroutes.CreatePublicTable).Methods(http.MethodPost)
	api.HandleFunc("/yapar/pub/get", yaparroutes.GetPublicTable).Methods(http.MethodPost)

	api.Handle("/scanners/priv/create", login.IsAuthorized(scanners.CreatePrivateScanner)).Methods(http.MethodPost)
	api.Handle("/scanners/priv/simulate", login.IsAuthorized(scanners.SimulatePrivateScanner)).Methods(http.MethodPost)
	api.HandleFunc("/scanners/public/create", scanners.CreatePublicScanner).Methods(http.MethodPost)
	api.HandleFunc("/scanners/public/simulate", scanners.SimulatePublicScanner).Methods(http.MethodPost)

	api.HandleFunc("/compare/simulate", compare.SimulateCompile).Methods(http.MethodPost)

	// Do route image with id for param
	r.HandleFunc("/image/{id}", scanners.GetImageHandler).Methods(http.MethodGet)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the serverc
	fmt.Println("Starting server at port 8000")
	r.Use(mux.CORSMethodMiddleware(r))
	err = http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))

	if err != nil {
		fmt.Println(err)
		// Finish the program
		return
	}
}
