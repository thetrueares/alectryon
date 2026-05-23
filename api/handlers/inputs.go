package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InputListResponse struct {
	Inputs []InputResponse
}

type InputResponse struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type InputCreateRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Active bool   `json:"active"`
}

type InputHandlers struct {
	repository models.InputRepository
}

func (lh InputHandlers) ListInputHandler(c *gin.Context) {

	inputs, err := lh.repository.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})

		return
	}

	var inputsList []InputResponse

	for _, inputModel := range inputs {
		inputResponse := ConvertModelToResponse(inputModel)
		inputsList = append(inputsList, inputResponse)
	}

	c.JSON(200, InputListResponse{Inputs: inputsList})
}

func (lh InputHandlers) CreateInputHandler(c *gin.Context) {

	var createBody InputCreateRequest

	if err := c.ShouldBindJSON(&createBody); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid body"})
		return
	}

	model := ConvertCreateRequestToModel(createBody)
	err := lh.repository.Save(model)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Input created successfully",
	})
}

func ConvertCreateRequestToModel(input InputCreateRequest) models.InputModel {
	return models.InputModel{
		// ID:        bson.NewObjectID(),
		Type:      models.InputType(input.Type),
		Name:      input.Name,
		Active:    input.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ConvertModelToResponse(input models.InputModel) InputResponse {

	return InputResponse{
		Id:        input.ID.Hex(),
		Type:      string(input.Type),
		Name:      input.Name,
		Active:    input.Active,
		CreatedAt: input.CreatedAt.Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.Format(time.RFC3339),
	}
}

func NewInputHandlers(collection *mongo.Collection) *InputHandlers {
	repository := models.NewInputRepository(collection)
	return &InputHandlers{repository: *repository}
}
