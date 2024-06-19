package login

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)
type contextKey string

const UserContextKey contextKey = "user"

// Secret key to sign the JWT
var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token
func GenerateToken(username, password string) (string, error) {
	// Set the expiration time of the token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the claims
	claims := &Claims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// Middleware to check JWT
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Store the username in the request context
		ctx := context.WithValue(r.Context(), UserContextKey, claims.Username)
		endpoint(w, r.WithContext(ctx))
	})
}


// Function to hash a password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Function to check a hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Struct to represent a simple login
type User struct {
	Username string `json:"username" example:"daniel"`
	Password string `json:"password" example:"1234"`
}

var Client *mongo.Client

// Function to create client and connect to the database
func CreateClientLogin() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI")).
		SetTLSConfig(tlsConfig)

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)

	return err
}

// Response struct represents the structure of the response
type GoodResponse struct {
	Message string `json:"message" example:"Success"`
	Status  int    `json:"status" example:"200"`
	Data    map[string]interface{} `json:"data"`
}

// Response struct represents the structure of the bad response
type BadResponse struct {
	Message string `json:"message" example:"Error"`
	Status  int    `json:"status" example:"400"`
}

// @Summary Create a user
// @Description Create a user
// @Tags login
// @Accept json
// @Produce json
// @Param login body User true "User object that needs to be created"
// @Success 200 {object} GoodResponse "User created successfully"
// @Failure 400 {object} BadResponse "Error creating user"
// @Router /api/v1/register [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var login User
	_ = json.NewDecoder(r.Body).Decode(&login)

	collection := Client.Database("GOLL1").Collection("users")

	// Check if the user already exists
	count, err := collection.CountDocuments(
		context.Background(),
		bson.D{{Key: "username", Value: login.Username}},
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: "User already exists",
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(login.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	_, err = collection.InsertOne(
		context.Background(),
		bson.D{
			{Key: "username", Value: login.Username},
			{Key: "password", Value: hashedPassword},
		})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Tokenize the user
	token, err := GenerateToken(login.Username, login.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	json.NewEncoder(w).Encode(
		GoodResponse{
			Message: "User created successfully",
			Status:  http.StatusOK,
			Data:    map[string]interface{}{"token": token},
		},
	)
}

// @Summary Login
// @Description Login
// @Tags login
// @Accept json
// @Produce json
// @Param login body User true "User object that needs to be created"
// @Success 200 {object} GoodResponse "User created successfully"
// @Failure 400 {object} BadResponse "Error creating user"
// @Router /api/v1/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	collection := Client.Database("GOLL1").Collection("users")

	// Find user for username
	var result bson.M
	err := collection.FindOne(
		context.Background(),
		bson.D{{Key: "username", Value: user.Username}},
	).Decode(&result)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Check if the password is correct
	if !CheckPasswordHash(user.Password, result["password"].(string)) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: "Invalid password",
				Status:  http.StatusBadRequest,
			})
		return
	}

	// Tokenize the user
	token, err := GenerateToken(user.Username, user.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			BadResponse{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			})
		return
	}

	json.NewEncoder(w).Encode(
		GoodResponse{
			Message: "User logged in successfully",
			Status:  http.StatusOK,
			Data:    map[string]interface{}{"token": token},
		},
	)
}
