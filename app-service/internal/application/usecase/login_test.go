package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type stubUserRepository struct {
	user *entity.User
}

func (s *stubUserRepository) FindByEmail(email string) (*entity.User, error) {
	if s.user == nil {
		return nil, nil
	}
	if s.user.Email != email {
		return nil, nil
	}
	return s.user, nil
}

func (s *stubUserRepository) FindByID(id entity.UserID) (*entity.User, error) {
	if s.user == nil {
		return nil, nil
	}
	if s.user.ID != id {
		return nil, nil
	}
	return s.user, nil
}

func TestLoginUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("plain-text-demo-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	name := "Demo User"
	repo := &stubUserRepository{
		user: &entity.User{
			ID:           entity.UserID("00000000-0000-0000-0000-000000000001"),
			Email:        "demo@example.com",
			PasswordHash: string(hashedPassword),
			Name:         &name,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	uc := NewLoginUsecase(repo)
	out, err := uc.Execute(context.Background(), LoginInput{
		Email:    "demo@example.com",
		Password: "plain-text-demo-password",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out == nil || out.User == nil {
		t.Fatalf("expected user in output")
	}
	if out.User.Email != "demo@example.com" {
		t.Fatalf("unexpected email: %s", out.User.Email)
	}
}

func TestLoginUsecase_Execute_InvalidPassword(t *testing.T) {
	t.Parallel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("plain-text-demo-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &stubUserRepository{
		user: &entity.User{
			ID:           entity.UserID("00000000-0000-0000-0000-000000000001"),
			Email:        "demo@example.com",
			PasswordHash: string(hashedPassword),
		},
	}

	uc := NewLoginUsecase(repo)
	_, err = uc.Execute(context.Background(), LoginInput{
		Email:    "demo@example.com",
		Password: "wrong-password",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

