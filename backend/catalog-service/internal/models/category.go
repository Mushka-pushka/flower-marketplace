package models

import (
	"time"

	"github.com/google/uuid"
)

// Category — категория товаров
type Category struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Slug        string     `json:"slug" db:"slug"`
	Description string     `json:"description,omitempty" db:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	ImageURL    string     `json:"image_url,omitempty" db:"image_url"`
	SortOrder   int        `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// CategoryWithCount — категория с количеством товаров
type CategoryWithCount struct {
	Category
	ProductCount int `json:"product_count"`
}

// АДМИН: УПРАВЛЕНИЕ КАТЕГОРИЯМИ

// CreateCategoryRequest — запрос на создание категории
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Slug        string     `json:"slug" binding:"required"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	ImageURL    string     `json:"image_url"`
	SortOrder   int        `json:"sort_order"`
}

// UpdateCategoryRequest — запрос на обновление категории
type UpdateCategoryRequest struct {
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	ImageURL    string     `json:"image_url"`
	SortOrder   int        `json:"sort_order"`
}