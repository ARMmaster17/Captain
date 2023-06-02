package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pingRequest struct {
	Value string `json:"value"`
}

type pingResponse struct {
	Value string `json:"value"`
}

func HandlePingPost(c *gin.Context) {
	var request pingRequest
	if c.ShouldBindJSON(&request) == nil {
		// TODO: Remove
		fmt.Printf("Got the following value: %s\n", request.Value)
		response := pingResponse{Value: "POOONG"}
		c.JSON(http.StatusOK, response)
		return
	}
	c.String(http.StatusBadRequest, "")
}
