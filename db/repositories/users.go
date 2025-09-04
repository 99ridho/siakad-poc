package repositories

import (
	"context"
	"errors"
	"siakad-poc/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (generated.User, error)
	GetUserByEmail(ctx context.Context, email string) (generated.User, error)
}

type DefaultUserRepository struct {
	query *generated.Queries
	pool  *pgxpool.Pool
}

func NewDefaultUserRepository(pool *pgxpool.Pool) *DefaultUserRepository {
	return &DefaultUserRepository{
		query: generated.New(pool),
		pool:  pool,
	}
}

func (r *DefaultUserRepository) GetUser(ctx context.Context, id string) (generated.User, error) {
	var uuidID pgtype.UUID
	err := uuidID.Scan(id)
	if err != nil {
		return generated.User{}, errors.New("can't parse id as uuid")
	}

	return r.query.GetUser(ctx, uuidID)
}

func (r *DefaultUserRepository) GetUserByEmail(ctx context.Context, email string) (generated.User, error) {
	return r.query.GetUserByEmail(ctx, email)
}
