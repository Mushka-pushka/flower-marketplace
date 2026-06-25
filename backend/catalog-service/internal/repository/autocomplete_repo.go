package repository

import (
	"context"
	"strings"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AutocompleteRepository struct {
	db *pgxpool.Pool
}

func NewAutocompleteRepository(db *pgxpool.Pool) *AutocompleteRepository {
	return &AutocompleteRepository{db: db}
}

// GetSuggestions — получает подсказки по запросу
func (r *AutocompleteRepository) GetSuggestions(ctx context.Context, query string, limit int) ([]models.AutocompleteSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	if len(query) < 2 {
		return []models.AutocompleteSuggestion{}, nil
	}

	searchTerm := strings.ToLower(query)
	var suggestions []models.AutocompleteSuggestion

	// 1. Подсказки по названиям товаров
	productQuery := `
		SELECT name, slug, 'product' as type,
			CASE 
				WHEN LOWER(name) = $1 THEN 100
				WHEN LOWER(name) LIKE $1 || '%' THEN 50
				WHEN LOWER(name) LIKE '%' || $1 || '%' THEN 10
				ELSE 0
			END as score
		FROM products
		WHERE is_active = true AND LOWER(name) LIKE '%' || $1 || '%'
		ORDER BY score DESC, name
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, productQuery, searchTerm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.AutocompleteSuggestion
		err := rows.Scan(&s.Text, &s.Slug, &s.Type, &s.Score)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, s)
	}

	// 2. Подсказки по категориям (если нужно больше)
	if len(suggestions) < limit {
		remaining := limit - len(suggestions)
		categoryQuery := `
			SELECT name, slug, 'category' as type,
				CASE 
					WHEN LOWER(name) = $1 THEN 100
					WHEN LOWER(name) LIKE $1 || '%' THEN 50
					ELSE 10
				END as score
			FROM categories
			WHERE LOWER(name) LIKE '%' || $1 || '%'
			ORDER BY score DESC, name
			LIMIT $2
		`

		rows, err := r.db.Query(ctx, categoryQuery, searchTerm, remaining)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var s models.AutocompleteSuggestion
			err := rows.Scan(&s.Text, &s.Slug, &s.Type, &s.Score)
			if err != nil {
				return nil, err
			}
			suggestions = append(suggestions, s)
		}
	}

	// 3. Подсказки по тегам (исправленная версия)
	if len(suggestions) < limit {
		remaining := limit - len(suggestions)
		tagQuery := `
			WITH tag_data AS (
				SELECT DISTINCT unnest(tags) as tag
				FROM products
				WHERE is_active = true
			)
			SELECT tag, 'tag' as type,
				CASE 
					WHEN LOWER(tag) = $1 THEN 100
					WHEN LOWER(tag) LIKE $1 || '%' THEN 50
					ELSE 10
				END as score
			FROM tag_data
			WHERE LOWER(tag) LIKE '%' || $1 || '%'
			ORDER BY score DESC, tag
			LIMIT $2
		`

		rows, err := r.db.Query(ctx, tagQuery, searchTerm, remaining)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var s models.AutocompleteSuggestion
			err := rows.Scan(&s.Text, &s.Type, &s.Score)
			if err != nil {
				return nil, err
			}
			s.Slug = s.Text
			suggestions = append(suggestions, s)
		}
	}

	return suggestions, nil
}