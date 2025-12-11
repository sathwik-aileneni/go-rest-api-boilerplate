package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/domain"
	"github.com/sathwik-aileneni/go-rest-api-boilerplate/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	if req.Email == "" || req.Name == "" {
		return nil, errors.New("email and name are required")
	}

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return nil, err
	}

	s.logger.Info("user created successfully", "user_id", user.ID)
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		s.logger.Error("failed to get user", "user_id", id, "error", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all users", "error", err)
		return nil, err
	}

	return users, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := s.repo.Update(ctx, id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		s.logger.Error("failed to update user", "user_id", id, "error", err)
		return nil, err
	}

	s.logger.Info("user updated successfully", "user_id", user.ID)
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		s.logger.Error("failed to delete user", "user_id", id, "error", err)
		return err
	}

	s.logger.Info("user deleted successfully", "user_id", id)
	return nil
}
