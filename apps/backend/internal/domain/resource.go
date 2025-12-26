package domain

import "time"

type Resource struct {
	ID           uint64    `json:"id" db:"id"`
	OwnerUserID  uint64    `json:"ownerUserId" db:"owner_user_id"`
	CategoryID   uint64    `json:"categoryId" db:"category_id"`
	Title        string    `json:"title" db:"title"`
	Description  *string   `json:"description" db:"description"`
	Location     *string   `json:"location" db:"location"`
	PricePerHour int       `json:"pricePerHour" db:"price_per_hour"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

type Category struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
