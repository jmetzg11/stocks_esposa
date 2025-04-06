package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StartSimulation(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"Message": "Simulation started"})
}
