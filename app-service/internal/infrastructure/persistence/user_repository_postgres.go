package persistence

import (
	"database/sql"
	"errors"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/repository"
)

var _ repository.UserRepository = (*UserRepositoryPostgres)(nil)

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(db *sql.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) FindByEmail(email string) (*entity.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)

	var u entity.User
	var name sql.NullString
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if name.Valid {
		u.Name = &name.String
	}
	return &u, nil
}

func (r *UserRepositoryPostgres) FindByID(id entity.UserID) (*entity.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)

	var u entity.User
	var name sql.NullString
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &name, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if name.Valid {
		u.Name = &name.String
	}
	return &u, nil
}

