package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.iain.rocks/alectryon/backend/channels"
	"go.iain.rocks/alectryon/backend/engine"
	"go.iain.rocks/alectryon/backend/engine/provider"
	"go.iain.rocks/alectryon/backend/entities"
	"go.iain.rocks/alectryon/backend/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.New()

	logger, _ := zap.NewProduction()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	r.Use(ginzap.RecoveryWithZap(logger, true))
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

	channelHandlers := handlers.NewChannelHandlers(channelRepository)
	channelHandlers.AddHandlers(r)

	aiModel := provider.NewOllama(logger)
	engineObj := engine.NewEngine(aiModel, historyRepository, taskRepository, logger)

	inputChan := make(chan engine.InputMessage)

	go channels.StartChannels(channelRepository, inputChan, logger)
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
