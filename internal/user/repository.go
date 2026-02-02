package user

import (
	"context"
	"github/mbpaiba/my-api/internal/db/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(queries *sqlc.Queries) *Repository {
	return &Repository{queries: queries}
}

func (r *Repository) FindAll(ctx context.Context) ([]User, error) {
	rows, err := r.queries.FindUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = User{
			ID:        row.ID.String(),
			ClerkID:   row.ClerkID,
			Email:     row.Email,
			Username:  row.Username,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}

	return users, nil
}
