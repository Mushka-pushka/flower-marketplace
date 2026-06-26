package repository

import (
	"context"
	"errors"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategorySlugExists = errors.New("category with this slug already exists")
)

type CategoryAdminRepository struct {
	db *pgxpool.Pool
}

func NewCategoryAdminRepository(db *pgxpool.Pool) *CategoryAdminRepository {
	return &CategoryAdminRepository{db: db}
}

// CreateCategory — создаёт категорию
func (r *CategoryAdminRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (id, name, slug, description, parent_id, image_url, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		category.ID,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.ImageURL,
		category.SortOrder,
		category.CreatedAt,
	)
	if err != nil {
		if pgx.ErrNoRows != nil {
			// Проверка на уникальность slug
			if err.Error() == "duplicate key value violates unique constraint \"categories_slug_key\"" {
				return ErrCategorySlugExists
			}
		}
		return err
	}
	return nil
}

// GetCategoryByID — получает категорию по ID
func (r *CategoryAdminRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	query := `
		SELECT 
			id, 
			name, 
			slug, 
			COALESCE(description, '') as description,
			parent_id, 
			COALESCE(image_url, '') as image_url,
			sort_order, 
			created_at
		FROM categories
		WHERE id = $1
	`

	var category models.Category
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ParentID,
		&category.ImageURL,
		&category.SortOrder,
		&category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// UpdateCategory — обновляет категорию
func (r *CategoryAdminRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories SET
			name = $1, slug = $2, description = $3, parent_id = $4,
			image_url = $5, sort_order = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.ParentID,
		category.ImageURL,
		category.SortOrder,
		category.ID,
	)
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint \"categories_slug_key\"" {
			return ErrCategorySlugExists
		}
		return err
	}
	return nil
}

// DeleteCategory — удаляет категорию (проверяет, что нет товаров)
func (r *CategoryAdminRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	// Проверяем, есть ли товары в этой категории
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE category_id = $1`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category with existing products")
	}

	// Проверяем, есть ли дочерние категории
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM categories WHERE parent_id = $1`, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category with subcategories")
	}

	query := `DELETE FROM categories WHERE id = $1`
	_, err = r.db.Exec(ctx, query, id)
	return err
}

// GetAllCategoriesAdmin — получает все категории (для админа)
func (r *CategoryAdminRepository) GetAllCategoriesAdmin(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT 
			id, 
			name, 
			slug, 
			COALESCE(description, '') as description,
			parent_id, 
			COALESCE(image_url, '') as image_url,
			sort_order, 
			created_at
		FROM categories
		ORDER BY sort_order, name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ParentID,
			&category.ImageURL,
			&category.SortOrder,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}