package handler

import (
	"net/http"
	"text/template"
	"wordwiz/internal/models"
	"wordwiz/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	uid, exists := c.Get("uid")
	if !exists {
		logger.Get().Error().Msg("error on accessing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on accessing uid"})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		logger.Get().Error().Msg("error on parsing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on parsing uid"})
		return
	}

	var req models.AddWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Get().Err(err).Msg("Failed to BindJSON on AddWord")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word, _ := h.stg.GetWordByLanguageAndWord(c.Request.Context(), req.Word.Lang, req.Word.Word)

	if word != nil {
		req.Word = *word
	} else {
		req.Word.ID = uuid.New().String()
	}

	err := h.stg.AddWord(c.Request.Context(), uidStr, req.Word, req.Definitions)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to AddWord")
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
//	@Success		200		{object}	models.UserWords
//	@Router			/user/words [GET]
func (h *Handler) GetUserWords(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		logger.Get().Error().Msg("error on accessing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on accessing uid"})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		logger.Get().Error().Msg("error on parsing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on parsing uid"})
		return
	}

	userWords, err := h.stg.GetUserWordList(c.Request.Context(), uidStr)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to GetUserWordList")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user words"})
		return
	}

	c.JSON(http.StatusOK, userWords)
}

type PageData struct {
	Words []struct {
		Word       string
		Definition string
	}
}

func (h *Handler) GetUserWordsPage(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		logger.Get().Error().Msg("error on accessing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on accessing uid"})
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		logger.Get().Error().Msg("error on parsing uid")
		c.JSON(http.StatusForbidden, gin.H{"error": "error on parsing uid"})
		return
	}

	userWords, err := h.stg.GetUserWordList(c.Request.Context(), uidStr)
	if err != nil {
		logger.Get().Err(err).Msg("Failed to GetUserWordList")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user words"})
		return
	}

	pagedata := PageData{}

	for i := range userWords.Words {
		pagedata.Words = append(pagedata.Words, struct {
			Word       string
			Definition string
		}{Word: userWords.Words[i].Word.Word, Definition: userWords.Words[i].Definitions[0].Definition})
	}

	tmpl, err := template.ParseFiles("internal/server/handler/templates/words_template.html")
	if err != nil {
		logger.Get().Err(err).Msg("error in parsing template file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if err := tmpl.Execute(c.Writer, pagedata); err != nil {
		logger.Get().Err(err).Msg("error in parsing template file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})

		return
	}
}
