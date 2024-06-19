package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"backend/routes/compare"
	"backend/routes/login"
	scanners "backend/routes/scanners"
	yaparroutes "backend/routes/yaparRoutes"

	_ "backend/docs" // This is to import generated docs

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/manucorporat/stats"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

var (
	ips        = stats.New()
	messages   = stats.New()
	users      = stats.New()
	mutexStats sync.RWMutex
	savedStats map[string]uint64
)

func statsWorker() {
	c := time.Tick(1 * time.Second)
	var lastMallocs uint64
	var lastFrees uint64
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		mutexStats.Lock()
		savedStats = map[string]uint64{
			"timestamp":  uint64(time.Now().Unix()),
			"HeapInuse":  stats.HeapInuse,
			"StackInuse": stats.StackInuse,
			"Mallocs":    stats.Mallocs - lastMallocs,
			"Frees":      stats.Frees - lastFrees,
			"Inbound":    uint64(messages.Get("inbound")),
			"Outbound":   uint64(messages.Get("outbound")),
			"Connected":  connectedUsers(),
		}
		lastMallocs = stats.Mallocs
		lastFrees = stats.Frees
		messages.Reset()
		mutexStats.Unlock()
	}
}

func connectedUsers() uint64 {
	connected := users.Get("connected") - users.Get("disconnected")
	if connected < 0 {
		return 0
	}
	return uint64(connected)
}

// Stats returns savedStats data.
func Stats() map[string]uint64 {
	mutexStats.RLock()
	defer mutexStats.RUnlock()

	return savedStats
}

// Wrapper to convert http.HandlerFunc to gin.HandlerFunc
func WrapHandler(h http.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
        h(c.Writer, c.Request)
    }
}

// Wrapper to convert http.Handler to gin.HandlerFunc
func WrapHandlerWithHandler(h http.Handler) gin.HandlerFunc {
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

// StartWorkers start starsWorker by goroutine.
func StartWorkers() {
	go statsWorker()
}

// @title GO-Api API
// @version 1.0
// @description This is an API to create a language ll1
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

	fmt.Println("Starting server...")
	fmt.Println(os.Getenv("MONGO_URI"), os.Getenv("PORT"))

	// Configure the runtime
	ConfigRuntime()

	// Start the workers
	StartWorkers()



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

	gin.SetMode(gin.ReleaseMode)

    r := gin.Default()

    // CORS middleware
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Define the API routes
    api := r.Group("/api/v1")
    {
        api.POST("/register", WrapHandler(login.CreateUser))
        api.POST("/login", WrapHandler(login.Login))
        api.POST("/yapar/priv/create", WrapHandlerWithHandler(login.IsAuthorized(yaparroutes.CreatePrivateTable)))
        api.POST("/yapar/priv/get", WrapHandlerWithHandler(login.IsAuthorized(yaparroutes.GetPrivateTable)))
        api.POST("/yapar/pub/create", WrapHandler(yaparroutes.CreatePublicTable))
        api.POST("/yapar/pub/get", WrapHandler(yaparroutes.GetPublicTable))

        api.POST("/scanners/priv/create", WrapHandlerWithHandler(login.IsAuthorized(scanners.CreatePrivateScanner)))
        api.POST("/scanners/priv/simulate", WrapHandlerWithHandler(login.IsAuthorized(scanners.SimulatePrivateScanner)))
        api.POST("/scanners/public/create", WrapHandler(scanners.CreatePublicScanner))
        api.POST("/scanners/public/simulate", WrapHandler(scanners.SimulatePublicScanner))

        api.POST("/compare/simulate", WrapHandler(compare.SimulateCompile))
    }

    // Route for image with id param
    r.GET("/image/:id", WrapHandler(scanners.GetImageHandler))

    // Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := os.Getenv("PORT")
    // Start the server
    fmt.Println("Starting server at port", port)
    err = r.Run(":" + port)
    if err != nil {
        fmt.Println(err)
        // Finish the program
        return
    }
}
