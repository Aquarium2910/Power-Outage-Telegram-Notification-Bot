package user

import (
	"NEW_LOE_bot/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"NEW_LOE_bot/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) Register(ctx context.Context, telegramID int64, username string, firstName string) (entity.User, error) {

	now := time.Now().UTC()
	UUID := uuid.New().String()

	user := entity.User{
		ID:         UUID,
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		IsActive:   true,
		State:      entity.Idle,
		CreatedAt:  now,
	}

	err := uc.repo.Store(ctx, &user)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - Register - uc.repo.Store: %w", err)
	}

	return user, nil
}

func (uc *UseCase) GetByTelegramID(ctx context.Context, telegramID int64) (entity.User, error) {

	user, err := uc.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetByID: %w", err)
	}

	return user, nil
}

func (uc *UseCase) UpdateIsActive(ctx context.Context, telegramID int64, isActive bool) error {
	user, err := uc.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return nil
		}
		return err
	}

	if user.IsActive != isActive {
		user.IsActive = isActive
		err = uc.repo.Update(ctx, &user)
		if err != nil {
			return fmt.Errorf("UserUseCase - UpdateIsActive - uc.repo.Store: %w", err)
		}
	}

	return nil
}
