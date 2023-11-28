package handler

import (
	"net/http"
	"wordwiz/internal/models"
	"wordwiz/pkg/logger"

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
	var req models.AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.stg.AddWord(c.Request.Context(), req.UserID, req.Word, req.Definitions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetUserWords godoc
//
//	@Summary		GetUserWords
//	@Description	GetUserWords
//	@Tags			word
//	@Accept			json
//	@Produce		json
//	@Param			user_id	path		string	true	"User ID"
//	@Success		200		{object}	models.UserWords
//	@Router			/user/{user_id}/words [GET]
func (h *Handler) GetUserWords(c *gin.Context) {
	userID := c.Param("user_id")

	userWords, err := h.stg.GetUserWordList(c.Request.Context(), userID)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to fetch user words")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user words"})
		return
	}

	c.JSON(http.StatusOK, userWords)
}
