package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	r := gin.Default()

	// Enable CORS for frontend access
	r.Use(cors.Default())

	mongo, err := createMongoDb()

	if err != nil {
		panic(err)
	}

	database := mongo.Database(os.Getenv("MONGODB_DATABASE"))
	log.Print("Connected")
	inputCollections := database.Collection("inputs")

	inputHandlers := handlers.NewInputHandlers(inputCollections)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World from Gin!",
		})
	})
	inputHandlers.AddHandlers(r)

	r.Run(":8080")
}

func createMongoDb() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	return mongo.Connect(opts)
}
