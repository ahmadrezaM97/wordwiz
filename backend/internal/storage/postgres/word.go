package postgres

import (
	"context"
	"fmt"

	"wordwiz/internal/models"

	"github.com/google/uuid"
)

func (s *Storage) AddWord(ctx context.Context, userID string, word models.Word, definitions []models.Definition) (err error) {
	// Begin a transaction
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx) // Rollback if any error occurs
		}
	}()

	wordID := uuid.New().String()

	fmt.Println(word.ImageURL)
	// Insert word into "words" table
	_, err = tx.Exec(ctx, `
		INSERT INTO "words" ("id", "language", "word", "example", "image_url", "link")
		VALUES ($1, $2, $3, $4, $5, $6)
	`, wordID, word.Lang, word.Word, word.Example, word.ImageURL, word.Link)
	if err != nil {
		return err
	}

	// Insert definitions into "definitions" table
	for _, definition := range definitions {
		_, err = tx.Exec(ctx, `
			INSERT INTO "definitions" ("language", "definition", "word_fk")
			VALUES ($1, $2, $3)
		`, definition.Lang, definition.Definition, wordID)
		if err != nil {
			return err
		}
	}

	// Insert data into "userwords" table
	_, err = tx.Exec(ctx, `
		INSERT INTO "userwords" ("user_fk", "word_fk", "status", "note")
		VALUES ($1, $2, $3, $4)
	`, userID, wordID, 1, "Some note about the word")
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
