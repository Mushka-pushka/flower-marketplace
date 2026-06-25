package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/config"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/models"
	"github.com/Mushka-pushka/flower-marketplace/backend/catalog-service/internal/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CatalogService struct {
	productRepo  *repository.ProductRepository
	cfg          *config.Config
	valkeyClient *redis.Client
}

func NewCatalogService(
	productRepo *repository.ProductRepository,
	cfg *config.Config,
	valkeyClient *redis.Client,
) *CatalogService {
	return &CatalogService{
		productRepo:  productRepo,
		cfg:          cfg,
		valkeyClient: valkeyClient,
	}
}

// CRUD ОПЕРАЦИИ

// CreateProduct — создание товара
func (s *CatalogService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}
	if req.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if req.ShopID == uuid.Nil {
		return nil, errors.New("shop_id is required")
	}
	if req.CategoryID == uuid.Nil {
		return nil, errors.New("category_id is required")
	}

	slug := generateSlug(req.Name)

	now := time.Now()
	product := &models.Product{
		ID:          uuid.New(),
		ShopID:      req.ShopID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Price:       req.Price,
		OldPrice:    req.OldPrice,
		Stock:       req.Stock,
		Unit:        req.Unit,
		Packaging:   req.Packaging,
		Tags:        normalizeTags(req.Tags),
		IsActive:    true,
		IsFeatured:  req.IsFeatured,
		Rating:      0,
		ViewsCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Очищаем кэш при создании нового товара
	go s.clearSearchCache()

	return product, nil
}

// GetProductByID — получение товара по ID
func (s *CatalogService) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid product id")
	}

	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = s.productRepo.IncrementViews(context.Background(), id)
	}()

	return product, nil
}

// GetProductBySlug — получение товара по slug
func (s *CatalogService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	if slug == "" {
		return nil, errors.New("slug is required")
	}
	return s.productRepo.GetProductBySlug(ctx, slug)
}

// СЕМАНТИЧЕСКИЙ ПОИСК С КЭШИРОВАНИЕМ

// SearchProducts — расширенный семантический поиск с кэшированием в Valkey
func (s *CatalogService) SearchProducts(ctx context.Context, req *models.SearchRequest) (*models.SearchResponse, error) {
	// Нормализуем запрос
	if req.Query != "" {
		tagsFromQuery := extractTagsFromQuery(req.Query)
		if len(req.Tags) == 0 {
			req.Tags = tagsFromQuery
		} else {
			req.Tags = append(req.Tags, tagsFromQuery...)
		}
	}
	req.Tags = uniqueStrings(req.Tags)

	if req.Limit <= 0 {
		req.Limit = 24
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if req.SortBy == "" {
		req.SortBy = "relevance"
	}

	// Формируем ключ для кэша
	cacheKey := fmt.Sprintf("search:%s:%s:%v:%v:%v:%d:%d:%s",
		req.Query,
		req.Category,
		req.Tags,
		req.MinPrice,
		req.MaxPrice,
		req.Limit,
		req.Offset,
		req.SortBy,
	)

	// Пробуем получить из кэша
	cached, err := s.valkeyClient.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp models.SearchResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			log.Printf("Cache hit for: %s", cacheKey)
			return &resp, nil
		}
	}

	// Кэша нет — выполняем поиск
	products, total, err := s.productRepo.SearchProducts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Формируем ответ
	resp := &models.SearchResponse{
		Items:    products,
		Total:    total,
		Limit:    req.Limit,
		Offset:   req.Offset,
		Query:    req.Query,
		TagsUsed: req.Tags,
		SortBy:   req.SortBy,
		HasMore:  int64(req.Offset+req.Limit) < total,
	}

	// Сохраняем в кэш на 5 минут
	data, _ := json.Marshal(resp)
	s.valkeyClient.Set(ctx, cacheKey, data, 5*time.Minute)
	log.Printf("Saved to cache: %s", cacheKey)

	return resp, nil
}

// clearSearchCache — очищает все кэши поиска (используется при создании/обновлении товаров)
func (s *CatalogService) clearSearchCache() {
	ctx := context.Background()
	iter := s.valkeyClient.Scan(ctx, 0, "search:*", 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		s.valkeyClient.Del(ctx, keys...)
		log.Printf("Cleared %d search cache entries", len(keys))
	}
}

// КАТЕГОРИИ

// GetCategories — получение всех категорий
func (s *CatalogService) GetCategories(ctx context.Context, withProducts bool) ([]models.CategoryWithCount, error) {
	return s.productRepo.GetCategoriesWithCount(ctx, withProducts)
}

// ОБНОВЛЕНИЕ И УДАЛЕНИЕ
// UpdateProduct — обновление товара
func (s *CatalogService) UpdateProduct(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid product id")
	}

	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != "" {
		product.Name = *req.Name
		product.Slug = generateSlug(*req.Name)
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil && *req.Price > 0 {
		product.Price = *req.Price
	}
	if req.OldPrice != nil {
		product.OldPrice = req.OldPrice
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.Packaging != nil {
		product.Packaging = *req.Packaging
	}
	if req.Tags != nil {
		product.Tags = normalizeTags(req.Tags)
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.CategoryID != nil && *req.CategoryID != uuid.Nil {
		product.CategoryID = *req.CategoryID
	}

	product.UpdatedAt = time.Now()

	err = s.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Очищаем кэш при обновлении товара
	go s.clearSearchCache()

	return product, nil
}

// DeleteProduct — мягкое удаление товара
func (s *CatalogService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product id")
	}
	err := s.productRepo.DeleteProduct(ctx, id)
	if err == nil {
		go s.clearSearchCache()
	}
	return err
}

// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = transliterate(slug)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 100 {
		slug = slug[:100]
	}
	return slug
}

func transliterate(text string) string {
	translitMap := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "е": "e",
		"ё": "e", "ж": "zh", "з": "z", "и": "i", "й": "y", "к": "k",
		"л": "l", "м": "m", "н": "n", "о": "o", "п": "p", "р": "r",
		"с": "s", "т": "t", "у": "u", "ф": "f", "х": "h", "ц": "ts",
		"ч": "ch", "ш": "sh", "щ": "shch", "ъ": "", "ы": "y", "ь": "",
		"э": "e", "ю": "yu", "я": "ya",
	}

	result := strings.Builder{}
	for _, char := range text {
		charStr := string(char)
		if translit, ok := translitMap[charStr]; ok {
			result.WriteString(translit)
		} else if unicode.IsLetter(char) {
			result.WriteString(charStr)
		}
	}
	return result.String()
}

func extractTagsFromQuery(query string) []string {
	query = strings.ToLower(query)

	synonymsMap := map[string][]string{
		"романтик":     {"романтика", "любовь", "свидание", "валентинка"},
		"свадьб":       {"свадьба", "венчание", "невеста", "жених"},
		"роз":          {"розы", "роза", "красные розы"},
		"тюльпан":      {"тюльпаны", "весна"},
		"пион":         {"пионы", "пышные"},
		"гортензи":     {"гортензии", "голубые"},
		"хризантем":    {"хризантемы", "осень"},
		"мам":          {"мама", "матери", "8 марта"},
		"дружб":        {"дружба", "подруга", "коллега"},
		"корпоратив":   {"корпоратив", "офис", "бизнес"},
		"извинен":      {"извинение", "прости", "сорри"},
		"благодарност": {"благодарность", "спасибо"},
		"нежн":         {"нежные", "пастельные", "кремовые"},
		"ярк":          {"яркие", "красные", "жёлтые"},
		"классик":      {"классические", "элегантные"},
		"экзотик":      {"экзотические", "необычные"},
	}

	tags := []string{}
	words := strings.Fields(query)

	for _, word := range words {
		if len(word) > 2 {
			tags = append(tags, word)
		}
		for key, synonyms := range synonymsMap {
			if strings.Contains(word, key) || strings.Contains(query, key) {
				tags = append(tags, synonyms...)
			}
		}
	}

	return uniqueStrings(tags)
}

func normalizeTags(tags []string) []string {
	normalized := []string{}
	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" {
			normalized = append(normalized, tag)
		}
	}
	return uniqueStrings(normalized)
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}