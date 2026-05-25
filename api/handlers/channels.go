package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChannelListResponse struct {
	Channels []ChannelResponse `json:"channels"`
}

type ChannelResponse struct {
	Id        string         `json:"id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Active    bool           `json:"active"`
	Options   map[string]any `json:"options"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type ChannelCreateRequest struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Active  bool           `json:"active"`
	Options map[string]any `json:"options"`
}

type ChannelUpdateRequest struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Active  bool           `json:"active"`
	Options map[string]any `json:"options"`
}

func NewChannelHandlers(repository *models.ChannelRepository) *ChannelHandlers {
	return &ChannelHandlers{repository: repository}
}

type ChannelHandlers struct {
	repository *models.ChannelRepository
}

func (lh ChannelHandlers) AddHandlers(r *gin.Engine) {

	r.GET("/channels", lh.ListChannelHandler)
	r.POST("/channels", lh.CreateChannelHandler)
	r.POST("/channels/:id/toggle", lh.ToogleChannelHandler)
	r.GET("/channels/:id", lh.FetchChannelHandler)
	r.POST("/channels/:id", lh.UpdateChannelHandler)
	r.DELETE("/channels/:id", lh.DeleteChannelHandler)
}

func (lh ChannelHandlers) ListChannelHandler(c *gin.Context) {

	inputs, err := lh.repository.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})

		return
	}

	var inputsList []ChannelResponse

	for _, inputModel := range inputs {
		inputResponse := ConvertModelToResponse(inputModel)
		inputsList = append(inputsList, inputResponse)
	}

	c.JSON(200, ChannelListResponse{Channels: inputsList})
}

func (lh ChannelHandlers) CreateChannelHandler(c *gin.Context) {

	var createBody ChannelCreateRequest

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

func (lh ChannelHandlers) FetchChannelHandler(c *gin.Context) {
	id := c.Param("id")
	input, err := lh.repository.GetById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error(), "id": id})

		return
	}

	c.JSON(http.StatusOK, ConvertModelToResponse(input))
}

func (lh ChannelHandlers) UpdateChannelHandler(c *gin.Context) {

	id := c.Param("id")
	input, err := lh.repository.GetById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error(), "id": id})

		return
	}

	var createBody ChannelUpdateRequest

	if err := c.ShouldBindJSON(&createBody); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Invalid body"})
		return
	}

	input = UpdateInputFromUpdateRequest(input, createBody)
	err = lh.repository.Save(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Input updated successfully",
	})
}

func (lh ChannelHandlers) ToogleChannelHandler(c *gin.Context) {
	id := c.Param("id")
	input, err := lh.repository.GetById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error(), "id": id})

		return
	}
	input.Active = !input.Active
	err = lh.repository.Save(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error.", "message": err.Error()})

		return
	}

	c.JSON(http.StatusAccepted, ConvertModelToResponse(input))
}

func (lh ChannelHandlers) DeleteChannelHandler(c *gin.Context) {

	id := c.Param("id")
	err := lh.repository.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "message": err.Error(), "id": id})

		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "success"})
}

func UpdateInputFromUpdateRequest(original models.ChannelEntity, update ChannelUpdateRequest) models.ChannelEntity {
	original.Name = update.Name
	original.Type = models.ChannelType(update.Type)
	original.Active = update.Active
	original.UpdatedAt = time.Now()

	return original
}

func ConvertCreateRequestToModel(channel ChannelCreateRequest) models.ChannelEntity {
	return models.ChannelEntity{
		ID:        bson.NewObjectID(),
		Type:      models.ChannelType(channel.Type),
		Name:      channel.Name,
		Active:    channel.Active,
		Options:   channel.Options,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ConvertModelToResponse(channel models.ChannelEntity) ChannelResponse {

	return ChannelResponse{
		Id:        channel.ID.Hex(),
		Type:      string(channel.Type),
		Name:      channel.Name,
		Active:    channel.Active,
		Options:   channel.Options,
		CreatedAt: channel.CreatedAt.Format(time.RFC3339),
		UpdatedAt: channel.UpdatedAt.Format(time.RFC3339),
	}
}
