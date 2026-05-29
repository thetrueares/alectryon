package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/channels"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/engine/provider"
	"go.iain.rocks/alectryon/api/entities"
	"go.iain.rocks/alectryon/api/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	r := gin.Default()

	// Enable CORS for frontend access
	r.Use(cors.Default())

	database := createMongoDb()
	inputCollections := database.Collection("inputs")
	historyCollection := database.Collection("history")
	userCollection := database.Collection("users")
	taskCollection := database.Collection("tasks")

	channelRepository := entities.NewChannelRepository(inputCollections)
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
	go channels.StartChannels(channelRepository, historyRepository, userRepository, engineObj)

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
