package handler

import (
	"net/http"
	"wordwiz/internal/models"

	"github.com/gin-gonic/gin"
)

// AddWord godoc
//
//	@Summary		Add word
//	@Description	Add word
//	@Tags			word
//	@Accept			json
//	@Produce		json
//	@Param			AddWordRequest	body	models.AddWordRequest	true	"AddWordRequest body"
//	@Router			/words/add [POST]
func (h *Handler) AddWord(c *gin.Context) {
	// Read data from JSON post
	var req models.AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add word using storage
	err := h.stg.AddWord(c.Request.Context(), req.UserID, req.Word, req.Definitions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success status to clients
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
