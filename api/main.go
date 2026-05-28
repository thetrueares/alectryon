package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/channels"
	"go.iain.rocks/alectryon/api/engine"
	"go.iain.rocks/alectryon/api/engine/vendor"
	"go.iain.rocks/alectryon/api/entities"
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
	historyCollection := database.Collection("history")
	userCollection := database.Collection("users")

	repository := entities.NewChannelRepository(inputCollections)
	inputHandlers := handlers.NewChannelHandlers(repository)

	historyRepository := entities.NewHistoryRepository(historyCollection)
	userRepository := entities.NewUserRepository(userCollection)

	aiModel := vendor.NewOllama()
	engine := engine.NewEngine(aiModel, *historyRepository)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World from Gin!",
		})
	})
	inputHandlers.AddHandlers(r)
	go startChannels(repository, historyRepository, userRepository, engine)

	r.Run(":8080")
}

func createMongoDb() (*mongo.Client, error) {
	uri := os.Getenv("MONGODB_URI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	return mongo.Connect(opts)
}

func startChannels(
	repository *entities.ChannelRepository,
	historyRepository *entities.HistoryRepository,
	userRepository *entities.UserRepository,
	engine engine.EngineInterface,
) {
	inputModels, err := repository.GetAll()

	if err != nil {
		panic(err.Error())
	}

	for _, input := range inputModels {
		if input.Type == entities.ChannelTypeTelegramBot {
			go channels.StartTelegramBot(input, historyRepository, userRepository, engine)
		}
	}
}
