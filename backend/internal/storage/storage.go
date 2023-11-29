package storage

import (
	"context"

	"wordwiz/internal/models"
)

type Storage interface {
	GetWordByLanguageAndWord(ctx context.Context, language string, word string) (*models.Word, error)
	AddWord(ctx context.Context, userID string, word models.Word, definitions []models.Definition) (err error)
	GetUserWordList(ctx context.Context, userID string) (resp models.UserWords, err error)
	SignInUp(ctx context.Context, email, fullName string) (userID string, err error)
}
