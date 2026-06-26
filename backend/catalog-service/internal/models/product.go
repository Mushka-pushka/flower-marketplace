package models

import (
	"time"

	"github.com/google/uuid"
)

// ОСНОВНЫЕ СУЩНОСТИ

// Product — товар (цветок / букет)
type Product struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	ShopID      uuid.UUID   `json:"shop_id" db:"shop_id"`
	CategoryID  uuid.UUID   `json:"category_id" db:"category_id"`
	Name        string      `json:"name" db:"name"`
	Slug        string      `json:"slug" db:"slug"`
	Description string      `json:"description,omitempty" db:"description"`
	Price       float64     `json:"price" db:"price"`
	OldPrice    *float64    `json:"old_price,omitempty" db:"old_price"`
	Stock       int         `json:"stock" db:"stock"`
	Unit        string      `json:"unit" db:"unit"`               // "шт", "букет"
	Packaging   string      `json:"packaging,omitempty" db:"packaging"`
	Tags        []string    `json:"tags" db:"tags"`               // для семантического поиска!
	IsActive    bool        `json:"is_active" db:"is_active"`
	IsFeatured  bool        `json:"is_featured" db:"is_featured"`
	Rating      float64     `json:"rating" db:"rating"`
	ViewsCount  int         `json:"views_count" db:"views_count"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// ProductImage — фотография товара
type ProductImage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	IsMain    bool      `json:"is_main" db:"is_main"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ЗАПРОСЫ И ОТВЕТЫ

// CreateProductRequest — запрос на создание товара
type CreateProductRequest struct {
	ShopID      uuid.UUID `json:"shop_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	OldPrice    *float64  `json:"old_price"`
	Stock       int       `json:"stock" binding:"gte=0"`
	Unit        string    `json:"unit"`
	Packaging   string    `json:"packaging"`
	Tags        []string  `json:"tags"`
	IsFeatured  bool      `json:"is_featured"`
}

// UpdateProductRequest — запрос на обновление товара
type UpdateProductRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	OldPrice    *float64   `json:"old_price,omitempty"`
	Stock       *int       `json:"stock,omitempty"`
	Unit        *string    `json:"unit,omitempty"`
	Packaging   *string    `json:"packaging,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
	IsFeatured  *bool      `json:"is_featured,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
}

// SearchRequest — запрос на семантический поиск
type SearchRequest struct {
	Query    string   `json:"query"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	SortBy   string   `json:"sort_by"` // price_asc, price_desc, rating, relevance, newest
}

// SearchResponse — ответ на поисковый запрос
type SearchResponse struct {
	Items      []Product `json:"items"`
	Total      int64     `json:"total"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
	Query      string    `json:"query,omitempty"`
	TagsUsed   []string  `json:"tags_used,omitempty"`
	SortBy     string    `json:"sort_by,omitempty"`
	HasMore    bool      `json:"has_more"`
}