package models

import (
    "time"

    "github.com/google/uuid"
)

type Shop struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    Name        string     `json:"name" db:"name"`
    Description string     `json:"description,omitempty" db:"description"`
    SellerID    uuid.UUID  `json:"seller_id" db:"seller_id"`
    IsVerified  bool       `json:"is_verified" db:"is_verified"`
    Rating      float64    `json:"rating" db:"rating"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}