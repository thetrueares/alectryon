package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.iain.rocks/alectryon/api/models"
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

func InputListHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World from Gin!!!",
	})
}

func ConvertModelToResponse(input models.InputModel) InputResponse {

	return InputResponse{
		Id:        input.ID.String(),
		Type:      string(input.Type),
		Active:    input.Active,
		CreatedAt: input.CreatedAt.Format(time.RFC3339),
	}
}
