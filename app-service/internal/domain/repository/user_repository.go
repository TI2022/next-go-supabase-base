package repository

import "github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindByID(id entity.UserID) (*entity.User, error)
}

