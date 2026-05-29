package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.iain.rocks/alectryon/backend/channels"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/engine/provider"
	"go.iain.rocks/alectryon/backend/entities"
	"go.iain.rocks/alectryon/backend/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// Enable CORS for frontend access
	r.Use(cors.Default())

	database := createMongoDb()
	channelCollections := database.Collection("channels")
	historyCollection := database.Collection("history")
	userCollection := database.Collection("users")
	taskCollection := database.Collection("tasks")

	channelRepository := entities.NewChannelRepository(channelCollections)
	historyRepository := entities.NewHistoryRepository(historyCollection)
	userRepository := entities.NewUserRepository(userCollection)
	taskRepository := entities.NewTaskRepository(taskCollection)

	inputHandlers := handlers.NewChannelHandlers(channelRepository)

	aiModel := provider.NewOllama()
	engineObj := engine.NewEngine(aiModel, historyRepository, taskRepository)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World from Gin!",
		})
	})
	inputHandlers.AddHandlers(r)

	inputChan := make(chan engine.InputMessage)

	go channels.StartChannels(channelRepository, inputChan)
	go engine.InputHandler(inputChan, historyRepository, userRepository, engineObj)
	r.Run(":8080")
}

func createMongoDb() *mongo.Database {
	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	mongoConn, err := mongo.Connect(opts)

	if err != nil {
		panic(err)
	}

	return mongoConn.Database(os.Getenv("MONGODB_DATABASE"))
}
