import client from './client'

// Типы
export interface Category {
  id: string
  name: string
  slug: string
  description?: string
  parent_id?: string
  sort_order: number
  created_at: string
}

export interface Product {
  id: string
  shop_id: string
  category_id: string
  name: string
  slug: string
  description: string
  price: number
  old_price?: number
  stock: number
  unit: string
  packaging?: string
  tags: string[]
  is_active: boolean
  is_featured: boolean
  rating: number
  views_count: number
  images?: string[]
  created_at: string
  updated_at: string
  shop_name?: string
}

export interface SearchResponse {
  items: Product[]
  total: number
  limit: number
  offset: number
  query?: string
  tags_used?: string[]
  sort_by?: string
  has_more: boolean
}

// Интерфейс для автодополнения
export interface AutocompleteSuggestion {
  text: string
  type: 'product' | 'category' | 'tag'
  slug: string
  score: number
}

// Интерфейс для адреса доставки
export interface DeliveryAddress {
  id: string
  user_id: string
  name: string
  address: string
  entrance: string
  floor: string
  intercom: string
  comment: string
  is_default: boolean
  created_at: string
  updated_at: string
}

// Получение категорий
export const getCategories = async (): Promise<Category[]> => {
  const response = await client.get('/catalog/categories')
  return response.data
}

// Поиск товаров
export const searchProducts = async (params: {
  q?: string
  category?: string
  tags?: string
  min_price?: number
  max_price?: number
  sort_by?: string
  limit?: number
  offset?: number
}): Promise<SearchResponse> => {
  const encodedParams = { ...params }
  if (encodedParams.q) {
    encodedParams.q = encodeURIComponent(encodedParams.q)
  }
  
  const response = await client.get('/catalog/search', { params })
  
  console.log('Full response:', response.data)
  console.log('Items with images:', response.data.items?.map((item: any) => ({
    name: item.name,
    images: item.images
  })))
  
  return response.data
}

// Получение товара по ID
export const getProductById = async (id: string): Promise<Product> => {
  const response = await client.get(`/catalog/products/${id}`)
  return response.data
}

// Автодополнение
export const getAutocomplete = async (
  query: string, 
  limit?: number
): Promise<AutocompleteSuggestion[]> => {
  const response = await client.get('/catalog/autocomplete', { 
    params: { q: query, limit } 
  })
  return response.data
}

export interface Review {
  id: string
  product_id: string
  user_id: string
  order_id?: string | null
  rating: number
  comment: string
  is_approved: boolean
  reply?: string | null
  reply_created_at?: string | null
  reply_updated_at?: string | null
  user_name?: string
  user_avatar?: string
  created_at: string
  updated_at?: string
}

// Получение отзывов на товар
export const getProductReviews = async (productId: string): Promise<Review[]> => {
  const response = await client.get('/catalog/reviews', { params: { product_id: productId } })
  return response.data
}

// Создание отзыва
export const createReview = async (data: {
  product_id: string
  rating: number
  comment: string
}): Promise<Review> => {
  const response = await client.post('/catalog/reviews', data)
  return response.data
}

// Создание адреса доставки
export const createAddress = async (data: {  
  name: string
  address: string
  entrance?: string
  floor?: string
  intercom?: string
  comment?: string
  is_default?: boolean
}): Promise<DeliveryAddress> => {
  const response = await client.post('/catalog/addresses', data)
  return response.data
}

// Обновление отзыва
export const updateReview = async (reviewId: string, data: {
  rating: number
  comment: string
}): Promise<Review> => {
  const response = await client.put('/catalog/reviews', data, { params: { id: reviewId } })
  return response.data
}

// Удаление отзыва
export const deleteReview = async (reviewId: string): Promise<void> => {
  await client.delete('/catalog/reviews', { params: { id: reviewId } })
}

// ДЛЯ ПРОДАВЦА
// Получение товаров продавца
export const getSellerProducts = async (): Promise<Product[]> => {
  const response = await client.get('/catalog/seller/products')
  return response.data
}

// Создание товара (для продавца)
export const createSellerProduct = async (data: FormData): Promise<Product> => {
  const response = await client.post('/catalog/seller/products', data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

// Обновление товара (для продавца)
export const updateSellerProduct = async (id: string, data: FormData): Promise<Product> => {
  const response = await client.put(`/catalog/seller/products/${id}`, data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}

// Удаление товара (мягкое)
export const deleteSellerProduct = async (id: string): Promise<void> => {
  await client.delete(`/catalog/seller/products/${id}`)
}


// ОТЗЫВЫ — ДЛЯ ПРОДАВЦА

// Интерфейс отзыва с ответом продавца
export interface ReviewWithReply extends Review {
  reply?: string
  reply_created_at?: string
  reply_updated_at?: string
  product_name?: string
  shop_id?: string
}

// Интерфейс ответа на отзыв
export interface ReplyRequest {
  text: string
}

// Получение отзывов на товары продавца
export const getSellerReviews = async (params?: {
  limit?: number
  offset?: number
}): Promise<{ reviews: ReviewWithReply[]; total: number; limit: number; offset: number; has_more: boolean }> => {
  const response = await client.get('/seller/reviews', { params })
  return response.data
}

// Добавление ответа на отзыв
export const addReplyToReview = async (reviewId: string, text: string): Promise<void> => {
  await client.post('/seller/reviews/reply', { text }, { params: { id: reviewId } })
}

// Обновление ответа на отзыв
export const updateReplyOnReview = async (reviewId: string, text: string): Promise<void> => {
  await client.put('/seller/reviews/reply', { text }, { params: { id: reviewId } })
}

// Удаление ответа на отзыв
export const deleteReplyFromReview = async (reviewId: string): Promise<void> => {
  await client.delete('/seller/reviews/reply', { params: { id: reviewId } })
}