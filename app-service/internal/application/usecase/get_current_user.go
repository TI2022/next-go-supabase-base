package usecase

import (
	"context"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/repository"
)

type GetCurrentUserInput struct {
	UserID entity.UserID
}

type GetCurrentUserOutput struct {
	User *entity.User
}

type GetCurrentUserUsecase struct {
	users repository.UserRepository
}

func NewGetCurrentUserUsecase(users repository.UserRepository) *GetCurrentUserUsecase {
	return &GetCurrentUserUsecase{users: users}
}

func (u *GetCurrentUserUsecase) Execute(ctx context.Context, in GetCurrentUserInput) (*GetCurrentUserOutput, error) {
	user, err := u.users.FindByID(in.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &GetCurrentUserOutput{User: nil}, nil
	}
	return &GetCurrentUserOutput{User: user}, nil
}

