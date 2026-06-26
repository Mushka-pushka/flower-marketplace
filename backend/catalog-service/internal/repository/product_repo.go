package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrSlugExists       = errors.New("slug already exists")
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// ============================================================
// CRUD ОПЕРАЦИИ
// ============================================================

// CreateProduct — создание нового товара
func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (
			id, shop_id, category_id, name, slug, description, price, old_price,
			stock, unit, packaging, tags, is_active, is_featured, rating, views_count,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.ShopID,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.OldPrice,
		product.Stock,
		product.Unit,
		product.Packaging,
		product.Tags,
		product.IsActive,
		product.IsFeatured,
		product.Rating,
		product.ViewsCount,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "products_slug_key") {
			return ErrSlugExists
		}
		return err
	}
	return nil
}

// GetProductByID — получение товара по ID
func (r *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, shop_id, category_id, name, slug, description, price, old_price,
			stock, unit, packaging, tags, is_active, is_featured, rating, views_count,
			created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.ShopID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.OldPrice,
		&product.Stock,
		&product.Unit,
		&product.Packaging,
		&product.Tags,
		&product.IsActive,
		&product.IsFeatured,
		&product.Rating,
		&product.ViewsCount,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// GetProductBySlug — получение товара по slug
func (r *ProductRepository) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `
		SELECT id, shop_id, category_id, name, slug, description, price, old_price,
			stock, unit, packaging, tags, is_active, is_featured, rating, views_count,
			created_at, updated_at
		FROM products
		WHERE slug = $1
	`

	var product models.Product
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&product.ID,
		&product.ShopID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.OldPrice,
		&product.Stock,
		&product.Unit,
		&product.Packaging,
		&product.Tags,
		&product.IsActive,
		&product.IsFeatured,
		&product.Rating,
		&product.ViewsCount,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// UpdateProduct — обновление товара
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products SET
			name = $1, slug = $2, description = $3, price = $4, old_price = $5,
			stock = $6, unit = $7, packaging = $8, tags = $9, is_active = $10,
			is_featured = $11, category_id = $12, updated_at = $13
		WHERE id = $14
	`

	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.OldPrice,
		product.Stock,
		product.Unit,
		product.Packaging,
		product.Tags,
		product.IsActive,
		product.IsFeatured,
		product.CategoryID,
		product.UpdatedAt,
		product.ID,
	)

	return err
}

// DeleteProduct — мягкое удаление товара
func (r *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// IncrementViews — увеличивает счётчик просмотров товара
func (r *ProductRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET views_count = views_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ============================================================
// СЕМАНТИЧЕСКИЙ ПОИСК
// ============================================================

// SearchProducts — семантический поиск товаров
func (r *ProductRepository) SearchProducts(ctx context.Context, req *models.SearchRequest) ([]models.Product, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// 1. Текстовый поиск
	if req.Query != "" {
		conditions = append(conditions, fmt.Sprintf(
			"to_tsvector('russian', name || ' ' || COALESCE(description, '')) @@ plainto_tsquery('russian', $%d)",
			argIndex,
		))
		args = append(args, req.Query)
		argIndex++
	}

	// 2. Фильтр по категории
	if req.Category != "" {
		conditions = append(conditions, fmt.Sprintf(
			"category_id IN (SELECT id FROM categories WHERE slug = $%d)",
			argIndex,
		))
		args = append(args, req.Category)
		argIndex++
	}

	// 3. Фильтр по тегам
	if len(req.Tags) > 0 {
		tagConditions := []string{}
		for _, tag := range req.Tags {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM unnest(tags) AS t WHERE LOWER(t) = LOWER($%d::text))", argIndex,
			))
			args = append(args, tag)
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
	}

	// 4. Фильтр по цене
	if req.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
		args = append(args, *req.MinPrice)
		argIndex++
	}
	if req.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
		args = append(args, *req.MaxPrice)
		argIndex++
	}

	// 5. Только активные товары
	conditions = append(conditions, "is_active = true")

	// Собираем WHERE
	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// 6. Сортировка
	sortClause := "ORDER BY rating DESC, views_count DESC"
	switch req.SortBy {
	case "price_asc":
		sortClause = "ORDER BY price ASC"
	case "price_desc":
		sortClause = "ORDER BY price DESC"
	case "rating":
		sortClause = "ORDER BY rating DESC"
	case "newest":
		sortClause = "ORDER BY created_at DESC"
	case "relevance":
		if req.Query != "" {
			sortClause = fmt.Sprintf(`
				ORDER BY ts_rank(to_tsvector('russian', name || ' ' || COALESCE(description, '')), plainto_tsquery('russian', $%d)) DESC,
				rating DESC, views_count DESC
			`, argIndex)
			args = append(args, req.Query)
			argIndex++
		}
	}

	// 7. Пагинация
	limit := req.Limit
	if limit <= 0 {
		limit = 24
	}
	if limit > 100 {
		limit = 100
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Запрос на список
	query := fmt.Sprintf(`
		SELECT id, shop_id, category_id, name, slug, description, price, old_price,
			stock, unit, packaging, tags, is_active, is_featured, rating, views_count,
			created_at, updated_at
		FROM products
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CategoryID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.OldPrice,
			&product.Stock,
			&product.Unit,
			&product.Packaging,
			&product.Tags,
			&product.IsActive,
			&product.IsFeatured,
			&product.Rating,
			&product.ViewsCount,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	// ============================================================
	// ПОДСЧЁТ ОБЩЕГО КОЛИЧЕСТВА 
	// ============================================================

	var countConditions []string
	var countArgsList []interface{}
	countArgIndex := 1

	if req.Query != "" {
		countConditions = append(countConditions, fmt.Sprintf(
			"to_tsvector('russian', name || ' ' || COALESCE(description, '')) @@ plainto_tsquery('russian', $%d)",
			countArgIndex,
		))
		countArgsList = append(countArgsList, req.Query)
		countArgIndex++
	}
	if req.Category != "" {
		countConditions = append(countConditions, fmt.Sprintf(
			"category_id IN (SELECT id FROM categories WHERE slug = $%d)",
			countArgIndex,
		))
		countArgsList = append(countArgsList, req.Category)
		countArgIndex++
	}
	if len(req.Tags) > 0 {
		tagConditions := []string{}
		for _, tag := range req.Tags {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM unnest(tags) AS t WHERE LOWER(t) = LOWER($%d::text))", countArgIndex,
			))
			countArgsList = append(countArgsList, tag)
			countArgIndex++
		}
		countConditions = append(countConditions, "("+strings.Join(tagConditions, " OR ")+")")
	}
	if req.MinPrice != nil {
		countConditions = append(countConditions, fmt.Sprintf("price >= $%d", countArgIndex))
		countArgsList = append(countArgsList, *req.MinPrice)
		countArgIndex++
	}
	if req.MaxPrice != nil {
		countConditions = append(countConditions, fmt.Sprintf("price <= $%d", countArgIndex))
		countArgsList = append(countArgsList, *req.MaxPrice)
		countArgIndex++
	}
	countConditions = append(countConditions, "is_active = true")

	countWhereClause := strings.Join(countConditions, " AND ")
	countQueryFinal := fmt.Sprintf(`SELECT COUNT(*) FROM products WHERE %s`, countWhereClause)

	var total int64
	err = r.db.QueryRow(ctx, countQueryFinal, countArgsList...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ============================================================
// КАТЕГОРИИ
// ============================================================

// GetCategories — получение всех категорий
func (r *ProductRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT id, name, slug, 
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

// GetCategoriesWithCount — получение категорий с количеством товаров
func (r *ProductRepository) GetCategoriesWithCount(ctx context.Context, withCount bool) ([]models.CategoryWithCount, error) {
	query := `
		SELECT c.id, c.name, c.slug, 
			COALESCE(c.description, '') as description, 
			c.parent_id, 
			COALESCE(c.image_url, '') as image_url, 
			c.sort_order, 
			c.created_at
	`

	if withCount {
		query += `, COUNT(p.id) as product_count`
	}

	query += `
		FROM categories c
	`

	if withCount {
		query += `
			LEFT JOIN products p ON p.category_id = c.id AND p.is_active = true
			GROUP BY c.id, c.name, c.slug, c.description, c.parent_id, c.image_url, c.sort_order, c.created_at
		`
	}

	query += ` ORDER BY c.sort_order, c.name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryWithCount
	for rows.Next() {
		var cat models.CategoryWithCount
		var category models.Category

		if withCount {
			err := rows.Scan(
				&category.ID,
				&category.Name,
				&category.Slug,
				&category.Description,
				&category.ParentID,
				&category.ImageURL,
				&category.SortOrder,
				&category.CreatedAt,
				&cat.ProductCount,
			)
			if err != nil {
				return nil, err
			}
			cat.Category = category
			categories = append(categories, cat)
		} else {
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
			cat.Category = category
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

// GetCategoryBySlug — получение категории по slug
func (r *ProductRepository) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `
		SELECT id, name, slug, 
			COALESCE(description, '') as description, 
			parent_id, 
			COALESCE(image_url, '') as image_url, 
			sort_order, 
			created_at
		FROM categories
		WHERE slug = $1
	`

	var category models.Category
	err := r.db.QueryRow(ctx, query, slug).Scan(
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

// ============================================================
// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
// ============================================================

// ExistsCategory — проверяет существование категории по ID
func (r *ProductRepository) ExistsCategory(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}

// ExistsShop — проверяет существование магазина по ID
func (r *ProductRepository) ExistsShop(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE id = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	return exists, err
}