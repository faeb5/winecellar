// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: wines.sql

package database

import (
	"context"
)

const createWine = `-- name: CreateWine :one
INSERT INTO wines (
    id,
    name,
    color,
    wine_maker,
    country,
    vintage
) VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, color, name, wine_maker, country, vintage, created_at, updated_at
`

type CreateWineParams struct {
	ID        string
	Name      string
	Color     string
	WineMaker string
	Country   string
	Vintage   int64
}

func (q *Queries) CreateWine(ctx context.Context, arg CreateWineParams) (Wine, error) {
	row := q.db.QueryRowContext(ctx, createWine,
		arg.ID,
		arg.Name,
		arg.Color,
		arg.WineMaker,
		arg.Country,
		arg.Vintage,
	)
	var i Wine
	err := row.Scan(
		&i.ID,
		&i.Color,
		&i.Name,
		&i.WineMaker,
		&i.Country,
		&i.Vintage,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllWines = `-- name: DeleteAllWines :exec
DELETE FROM wines
`

func (q *Queries) DeleteAllWines(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllWines)
	return err
}
