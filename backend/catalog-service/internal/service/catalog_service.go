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
	cartRepo     *repository.CartRepository
	favoriteRepo  *repository.FavoriteRepository
	reviewRepo    *repository.ReviewRepository
	autocompleteRepo  *repository.AutocompleteRepository
	addressRepo       *repository.AddressRepository
	categoryAdminRepo  *repository.CategoryAdminRepository
	cfg          *config.Config
	valkeyClient *redis.Client
}

func NewCatalogService(
	productRepo *repository.ProductRepository,
	cartRepo *repository.CartRepository,
	favoriteRepo *repository.FavoriteRepository,
	reviewRepo *repository.ReviewRepository,
	autocompleteRepo *repository.AutocompleteRepository,
	addressRepo *repository.AddressRepository,
	categoryAdminRepo *repository.CategoryAdminRepository,
	cfg *config.Config,
	valkeyClient *redis.Client,
) *CatalogService {
	return &CatalogService{
		productRepo:  productRepo,
		cartRepo:     cartRepo,
		favoriteRepo:  favoriteRepo,
		reviewRepo:    reviewRepo,
		autocompleteRepo:  autocompleteRepo,
		addressRepo:       addressRepo,
		categoryAdminRepo:  categoryAdminRepo,
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

	cached, err := s.valkeyClient.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp models.SearchResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			log.Printf("Cache hit for: %s", cacheKey)
			return &resp, nil
		}
	}

	products, total, err := s.productRepo.SearchProducts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

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

	data, _ := json.Marshal(resp)
	s.valkeyClient.Set(ctx, cacheKey, data, 5*time.Minute)
	log.Printf("Saved to cache: %s", cacheKey)

	return resp, nil
}

// clearSearchCache — очищает все кэши поиска
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

// DecreaseStock — уменьшает количество товара на складе
func (s *CatalogService) DecreaseStock(ctx context.Context, productID uuid.UUID, quantity int) error {
    return s.productRepo.DecreaseStock(ctx, productID, quantity)
}

// КОРЗИНА (CART)

// AddToCart — добавляет товар в корзину
func (s *CatalogService) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	return s.cartRepo.AddToCart(ctx, userID, productID, quantity)
}

// GetCart — получает корзину пользователя
func (s *CatalogService) GetCart(ctx context.Context, userID uuid.UUID) ([]models.CartItemWithProduct, error) {
	return s.cartRepo.GetCartByUserID(ctx, userID)
}

// UpdateCartItem — обновляет количество товара в корзине
func (s *CatalogService) UpdateCartItem(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return s.cartRepo.UpdateCartItemQuantity(ctx, cartItemID, quantity)
}

// RemoveFromCart — удаляет товар из корзины
func (s *CatalogService) RemoveFromCart(ctx context.Context, cartItemID uuid.UUID) error {
	return s.cartRepo.RemoveFromCart(ctx, cartItemID)
}

// ИЗБРАННОЕ (FAVORITES)

// AddFavorite — добавляет товар в избранное
func (s *CatalogService) AddFavorite(ctx context.Context, userID, productID uuid.UUID) error {
	return s.favoriteRepo.AddFavorite(ctx, userID, productID)
}

// GetFavorites — получает все избранные товары пользователя
func (s *CatalogService) GetFavorites(ctx context.Context, userID uuid.UUID) ([]models.FavoriteWithProduct, error) {
	return s.favoriteRepo.GetFavoritesByUserID(ctx, userID)
}

// RemoveFavorite — удаляет товар из избранного
func (s *CatalogService) RemoveFavorite(ctx context.Context, userID, productID uuid.UUID) error {
	return s.favoriteRepo.RemoveFavorite(ctx, userID, productID)
}

// IsFavorite — проверяет, находится ли товар в избранном
func (s *CatalogService) IsFavorite(ctx context.Context, userID, productID uuid.UUID) (bool, error) {
	return s.favoriteRepo.IsFavorite(ctx, userID, productID)
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

// ОТЗЫВЫ (REVIEWS)

// CreateReview — создаёт отзыв
func (s *CatalogService) CreateReview(ctx context.Context, req *models.CreateReviewRequest, userID uuid.UUID) (*models.Review, error) {
	now := time.Now()
	review := &models.Review{
		ID:         uuid.New(),
		ProductID:  req.ProductID,
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsApproved: true, 
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if req.OrderID != uuid.Nil {
		review.OrderID = &req.OrderID
	}

	err := s.reviewRepo.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	// Обновляем рейтинг товара
	go s.updateProductRating(ctx, req.ProductID)

	return review, nil
}

// GetProductReviews — получает отзывы на товар
func (s *CatalogService) GetProductReviews(ctx context.Context, productID uuid.UUID) ([]models.ReviewWithUser, error) {
	return s.reviewRepo.GetReviewsByProductID(ctx, productID)
}

// GetMyReviews — получает отзывы пользователя
func (s *CatalogService) GetMyReviews(ctx context.Context, userID uuid.UUID) ([]models.ReviewWithUser, error) {
	return s.reviewRepo.GetReviewsByUserID(ctx, userID)
}

// UpdateReview — обновляет отзыв
func (s *CatalogService) UpdateReview(ctx context.Context, reviewID uuid.UUID, req *models.UpdateReviewRequest, userID uuid.UUID) (*models.Review, error) {
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	// Проверяем, что отзыв принадлежит пользователю
	if review.UserID != userID {
		return nil, errors.New("you can only update your own reviews")
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}
	review.UpdatedAt = time.Now()

	err = s.reviewRepo.UpdateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	go s.updateProductRating(ctx, review.ProductID)

	return review, nil
}

// DeleteReview — удаляет отзыв
func (s *CatalogService) DeleteReview(ctx context.Context, reviewID uuid.UUID, userID uuid.UUID) error {
	review, err := s.reviewRepo.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return errors.New("you can only delete your own reviews")
	}

	err = s.reviewRepo.DeleteReview(ctx, reviewID)
	if err == nil {
		go s.updateProductRating(ctx, review.ProductID)
	}
	return err
}

// updateProductRating — обновляет рейтинг товара
func (s *CatalogService) updateProductRating(ctx context.Context, productID uuid.UUID) {
	avgRating, count, err := s.reviewRepo.GetAverageRating(ctx, productID)
	if err != nil {
		return
	}

	// Обновляем рейтинг в товаре
	// Здесь можно добавить метод в productRepo для обновления рейтинга
	log.Printf("Product %s: avg rating %.2f, %d reviews", productID, avgRating, count)
}

// ApproveReview — одобряет отзыв (для админа)
func (s *CatalogService) ApproveReview(ctx context.Context, reviewID uuid.UUID) error {
	return s.reviewRepo.ApproveReview(ctx, reviewID)
}

// АВТОДОПОЛНЕНИЕ

// GetAutocompleteSuggestions — получает подсказки для поиска
func (s *CatalogService) GetAutocompleteSuggestions(ctx context.Context, query string, limit int) ([]models.AutocompleteSuggestion, error) {
	return s.autocompleteRepo.GetSuggestions(ctx, query, limit)
}

// АДРЕСА ДОСТАВКИ

// CreateAddress — создаёт адрес доставки
func (s *CatalogService) CreateAddress(ctx context.Context, userID uuid.UUID, req *models.CreateAddressRequest) (*models.DeliveryAddress, error) {
	now := time.Now()
	address := &models.DeliveryAddress{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       req.Name,
		Address:    req.Address,
		Entrance:   req.Entrance,
		Floor:      req.Floor,
		Intercom:   req.Intercom,
		Comment:    req.Comment,
		IsDefault:  req.IsDefault,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.addressRepo.CreateAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}

// GetAddresses — получает все адреса пользователя
func (s *CatalogService) GetAddresses(ctx context.Context, userID uuid.UUID) ([]models.DeliveryAddress, error) {
	return s.addressRepo.GetAddressesByUserID(ctx, userID)
}

// UpdateAddress — обновляет адрес
func (s *CatalogService) UpdateAddress(ctx context.Context, addressID, userID uuid.UUID, req *models.UpdateAddressRequest) (*models.DeliveryAddress, error) {
	address, err := s.addressRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	if address.UserID != userID {
		return nil, errors.New("you can only update your own addresses")
	}

	if req.Name != "" {
		address.Name = req.Name
	}
	if req.Address != "" {
		address.Address = req.Address
	}
	address.Entrance = req.Entrance
	address.Floor = req.Floor
	address.Intercom = req.Intercom
	address.Comment = req.Comment
	address.IsDefault = req.IsDefault
	address.UpdatedAt = time.Now()

	err = s.addressRepo.UpdateAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	return address, nil
}

// DeleteAddress — удаляет адрес
func (s *CatalogService) DeleteAddress(ctx context.Context, addressID, userID uuid.UUID) error {
	address, err := s.addressRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return errors.New("you can only delete your own addresses")
	}

	return s.addressRepo.DeleteAddress(ctx, addressID)
}

// SetDefaultAddress — устанавливает адрес по умолчанию
func (s *CatalogService) SetDefaultAddress(ctx context.Context, addressID, userID uuid.UUID) error {
	address, err := s.addressRepo.GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return errors.New("you can only set your own addresses as default")
	}

	return s.addressRepo.SetDefaultAddress(ctx, userID, addressID)
}

// ============================================================
// АДМИН: УПРАВЛЕНИЕ КАТЕГОРИЯМИ
// ============================================================

// AdminCreateCategory — создаёт категорию (админ)
func (s *CatalogService) AdminCreateCategory(ctx context.Context, req *models.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		CreatedAt:   time.Now(),
	}

	err := s.categoryAdminRepo.CreateCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// AdminGetAllCategories — получает все категории (админ)
func (s *CatalogService) AdminGetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryAdminRepo.GetAllCategoriesAdmin(ctx)
}

// AdminGetCategoryByID — получает категорию по ID (админ)
func (s *CatalogService) AdminGetCategoryByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	return s.categoryAdminRepo.GetCategoryByID(ctx, id)
}

// AdminUpdateCategory — обновляет категорию (админ)
func (s *CatalogService) AdminUpdateCategory(ctx context.Context, id uuid.UUID, req *models.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.categoryAdminRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.ImageURL != "" {
		category.ImageURL = req.ImageURL
	}
	if req.SortOrder > 0 {
		category.SortOrder = req.SortOrder
	}

	err = s.categoryAdminRepo.UpdateCategory(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// AdminDeleteCategory — удаляет категорию (админ)
func (s *CatalogService) AdminDeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.categoryAdminRepo.DeleteCategory(ctx, id)
}