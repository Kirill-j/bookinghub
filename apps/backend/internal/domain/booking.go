package domain

import "time"

type BookingStatus string

const (
	BookingPending  BookingStatus = "PENDING"
	BookingApproved BookingStatus = "APPROVED"
	BookingRejected BookingStatus = "REJECTED"
	BookingCanceled BookingStatus = "CANCELED"
)

type Booking struct {
	ID             uint64        `json:"id" db:"id"`
	ResourceID     uint64        `json:"resourceId" db:"resource_id"`
	UserID         uint64        `json:"userId" db:"user_id"`
	StartAt        time.Time     `json:"startAt" db:"start_at"`
	EndAt          time.Time     `json:"endAt" db:"end_at"`
	Status         BookingStatus `json:"status" db:"status"`
	ManagerComment *string       `json:"managerComment" db:"manager_comment"`
	CreatedAt      time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt      *time.Time    `json:"updatedAt" db:"updated_at"`
}
