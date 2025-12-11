package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, id int64, user *domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	query := `
		INSERT INTO users (email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, created_at, updated_at
	`

	user := &domain.User{}
	now := time.Now()

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Email,
		req.Name,
		now,
		now,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	query := `
		UPDATE users
		SET email = COALESCE(NULLIF($1, ''), email),
		    name = COALESCE(NULLIF($2, ''), name),
		    updated_at = $3
		WHERE id = $4
		RETURNING id, email, name, created_at, updated_at
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		req.Email,
		req.Name,
		time.Now(),
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
