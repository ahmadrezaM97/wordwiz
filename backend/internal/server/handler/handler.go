package handler

import (
	"net/http"
	"wordwiz/config"
	"wordwiz/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg config.Config
	stg storage.Storage
}

func New(cfg config.Config, stg storage.Storage) *Handler {
	return &Handler{
		cfg: cfg,
		stg: stg,
	}
}

// Health godoc
//	@Summary		show http server health
//	@Description	to check http server health
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Router			/health [GET]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app":    h.cfg.App,
		"http":   h.cfg.HTTP,
		"status": "OK",
	})
}
