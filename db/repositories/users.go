package repositories

import (
	"context"
	"errors"
	"math/big"
	"siakad-poc/db/generated"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (generated.User, error)
	GetUserByEmail(ctx context.Context, email string) (generated.User, error)
	CreateUser(ctx context.Context, email, password string, role int64) (generated.User, error)
}

type DefaultUserRepository struct {
	query *generated.Queries
	pool  *pgxpool.Pool
}

// Compile time interface conformance check
var _ UserRepository = (*DefaultUserRepository)(nil)

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

func (r *DefaultUserRepository) CreateUser(ctx context.Context, email, password string, role int64) (generated.User, error) {
	params := generated.CreateUserParams{
		Email:    email,
		Password: password,
		Role: pgtype.Numeric{
			Int:   big.NewInt(role),
			Valid: true,
		},
	}

	return r.query.CreateUser(ctx, params)
}
