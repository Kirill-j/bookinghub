package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeBookingRepo struct {
	hasConflictFn func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error)
	createFn      func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error)
}

func (f *fakeBookingRepo) HasConflict(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
	return f.hasConflictFn(ctx, resourceID, startAt, endAt)
}

func (f *fakeBookingRepo) Create(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
	return f.createFn(ctx, resourceID, userID, startAt, endAt)
}

func TestBookingService_Create_InvalidIDs(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			t.Fatal("should not call HasConflict")
			return false, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	now := time.Now().Add(1 * time.Hour)
	_, err := s.Create(context.Background(), 0, 1, now, now.Add(time.Hour))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBookingService_Create_InvalidTime_EndNotAfterStart(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			t.Fatal("should not call HasConflict")
			return false, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	now := time.Now().Add(1 * time.Hour)
	_, err := s.Create(context.Background(), 1, 1, now, now)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBookingService_Create_MinDuration(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			t.Fatal("should not call HasConflict")
			return false, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	start := time.Now().Add(2 * time.Hour)
	end := start.Add(10 * time.Minute)
	_, err := s.Create(context.Background(), 1, 1, start, end)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBookingService_Create_PastStart(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			t.Fatal("should not call HasConflict")
			return false, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	start := time.Now().Add(-10 * time.Minute)
	end := time.Now().Add(1 * time.Hour)
	_, err := s.Create(context.Background(), 1, 1, start, end)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBookingService_Create_Conflict(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			return true, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create when conflict")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	start := time.Now().Add(2 * time.Hour)
	end := start.Add(1 * time.Hour)
	_, err := s.Create(context.Background(), 10, 20, start, end)
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict, got: %v", err)
	}
}

func TestBookingService_Create_RepoError(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			return false, errors.New("db down")
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			t.Fatal("should not call Create when hasConflict returns error")
			return 0, nil
		},
	}
	s := NewBookingService(repo)

	start := time.Now().Add(2 * time.Hour)
	end := start.Add(1 * time.Hour)
	_, err := s.Create(context.Background(), 1, 1, start, end)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestBookingService_Create_OK(t *testing.T) {
	repo := &fakeBookingRepo{
		hasConflictFn: func(ctx context.Context, resourceID uint64, startAt, endAt time.Time) (bool, error) {
			return false, nil
		},
		createFn: func(ctx context.Context, resourceID, userID uint64, startAt, endAt time.Time) (uint64, error) {
			if resourceID != 11 || userID != 22 {
				t.Fatalf("unexpected ids")
			}
			return 777, nil
		},
	}
	s := NewBookingService(repo)

	start := time.Now().Add(2 * time.Hour)
	end := start.Add(1 * time.Hour)
	id, err := s.Create(context.Background(), 22, 11, start, end)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 777 {
		t.Fatalf("expected id=777, got %d", id)
	}
}
