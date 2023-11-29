package postgres

import (
	"context"
)

func (s *Storage) SignInUp(ctx context.Context, email, fullName string) (userID string, err error) {
	// Use a single query to either update the existing user or insert a new one
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		ON CONFLICT (email)
		DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`

	err = s.conn.QueryRow(ctx, query, fullName, email).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
