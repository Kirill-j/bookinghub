package service

import (
	"context"
	"errors"
	"time"

	"bookinghub-backend/internal/repo"
)

var (
	ErrInvalidTime = errors.New("Некорректный интервал времени")
	ErrConflict    = errors.New("Выбранное время уже занято")
)

type BookingService struct {
	repo *repo.BookingRepo
}

func NewBookingService(repo *repo.BookingRepo) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create(ctx context.Context, userID, resourceID uint64, startAt, endAt time.Time) (uint64, error) {
	if userID == 0 || resourceID == 0 {
		return 0, ErrInvalidTime
	}
	if !endAt.After(startAt) {
		return 0, ErrInvalidTime
	}

	if endAt.Sub(startAt) < 30*time.Minute {
		return 0, errors.New("Минимальная длительность бронирования: 30 минут")
	}

	if startAt.Before(time.Now().Add(-1 * time.Minute)) {
		return 0, errors.New("Нельзя бронировать время в прошлом")
	}

	conflict, err := s.repo.HasConflict(ctx, resourceID, startAt, endAt)
	if err != nil {
		return 0, err
	}
	if conflict {
		return 0, ErrConflict
	}

	return s.repo.Create(ctx, resourceID, userID, startAt, endAt)
}
