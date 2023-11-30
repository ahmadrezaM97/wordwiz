package postgres

import (
	"context"

	"wordwiz/internal/models"

	"github.com/jackc/pgx"
)

func (s *Storage) GetWordByLanguageAndWord(ctx context.Context, language string, word string) (resp *models.Word, err error) {
	var result models.Word

	err = s.conn.QueryRow(ctx, `
        SELECT id, language, word, example, image_url, link FROM "words"
        WHERE "language" = $1 AND "word" = $2
    `, language, word).Scan(
		&result.ID,
		&result.Lang,
		&result.Word,
		&result.Example,
		&result.ImageURL,
		&result.Link,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

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

	// Insert word into "words" table
	_, err = tx.Exec(ctx, `
    INSERT INTO "words" ("id", "language", "word", "example", "image_url", "link")
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT ("language", "word") DO NOTHING
	`, word.ID, word.Lang, word.Word, word.Example, word.ImageURL, word.Link)
	if err != nil {
		return err
	}

	// Insert definitions into "definitions" table
	for _, definition := range definitions {
		_, err = tx.Exec(ctx, `
			INSERT INTO "definitions" ("language", "definition", "word_fk")
			VALUES ($1, $2, $3)
			ON CONFLICT ("language", "word_fk") DO NOTHING
		`, definition.Lang, definition.Definition, word.ID)
		if err != nil {
			return err
		}
	}

	// Insert data into "userwords" table
	_, err = tx.Exec(ctx, `
		INSERT INTO "userwords" ("user_fk", "word_fk", "status", "note")
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ("user_fk", "word_fk") DO NOTHING
	`, userID, word.ID, 1, "Some note about the word")
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

func (s *Storage) GetUserWordList(ctx context.Context, userID string) (resp models.UserWords, err error) {
	rows, err := s.conn.Query(ctx, `
		SELECT w.id, w.language, w.word, w.example, w.image_url, w.link,
			d.id, d.language, d.definition, d.word_fk
		FROM userwords uw
		INNER JOIN words w ON uw.word_fk = w.id
		LEFT JOIN definitions d ON w.id = d.word_fk
		WHERE uw.user_fk = $1
		ORDER BY w.created_at DESC
	`, userID)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	userWordsMap := make(map[string]models.UserWord)
	for rows.Next() {
		var wordID, wordLanguage, wordText, wordExample, wordImageURL, wordLink string
		var definitionID, definitionLanguage, definitionText, definitionWordID string

		err := rows.Scan(
			&wordID, &wordLanguage, &wordText, &wordExample, &wordImageURL, &wordLink,
			&definitionID, &definitionLanguage, &definitionText, &definitionWordID,
		)

		if err != nil {
			return resp, err
		}

		userWord, exists := userWordsMap[wordID]
		if !exists {
			userWord = models.UserWord{
				Word: models.Word{
					Lang:     wordLanguage,
					Word:     wordText,
					Example:  wordExample,
					ImageURL: wordImageURL,
					Link:     wordLink,
				},
			}
		}

		userWord.Definitions = append(userWord.Definitions, models.Definition{
			Lang:       definitionLanguage,
			Definition: definitionText,
		})

		userWordsMap[wordID] = userWord
	}

	for _, userWord := range userWordsMap {
		resp.Words = append(resp.Words, userWord)
	}

	return resp, nil
}
