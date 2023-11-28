package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(ctx context.Context, pgUrl string) (*Storage, error) {
	conn, err := pgx.Connect(ctx, pgUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Storage{
		conn: conn,
	}, nil
}

func (c *Storage) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}
