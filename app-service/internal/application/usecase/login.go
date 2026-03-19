package usecase

import (
	"context"
	"errors"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User *entity.User
}

type LoginUsecase struct {
	users repository.UserRepository
}

func NewLoginUsecase(users repository.UserRepository) *LoginUsecase {
	return &LoginUsecase{users: users}
}

func (u *LoginUsecase) Execute(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	user, err := u.users.FindByEmail(in.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LoginOutput{User: user}, nil
}

