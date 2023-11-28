package storage

import (
	"context"

	"wordwiz/internal/models"
)

type Storage interface {
	AddWord(ctx context.Context, userID string, word models.Word, definitions []models.Definition) (err error)
}
