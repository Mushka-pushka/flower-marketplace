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
  created_at: string
  updated_at: string
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
  const response = await client.get('/catalog/search', { params })
  return response.data
}

// Получение товара по ID
export const getProductById = async (id: string): Promise<Product> => {
  const response = await client.get('/catalog/products', { params: { id } })
  return response.data
}

// Автодополнение
export const getAutocomplete = async (query: string, limit?: number): Promise<{ text: string; type: string; slug: string; score: number }[]> => {
  const response = await client.get('/catalog/autocomplete', { params: { q: query, limit } })
  return response.data
}