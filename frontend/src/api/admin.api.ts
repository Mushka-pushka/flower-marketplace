import client from './client'

// Типы
export interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  phone: string
  role: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface SellerWithShop {
  shop_id: string
  shop_name: string
  shop_description?: string
  is_verified: boolean
  rating: number
  user_id: string
  email: string
  phone: string
  first_name: string
  last_name: string
  is_active: boolean
  created_at: string
}

export interface UsersListResponse {
  users: User[]
  total: number
  limit: number
  offset: number
}

// УПРАВЛЕНИЕ ПОЛЬЗОВАТЕЛЯМИ
export const adminGetUsers = async (params?: {
  search?: string
  role?: string
  is_active?: boolean
  limit?: number
  offset?: number
}): Promise<UsersListResponse> => {
  const response = await client.get('/admin/users/list', { params })
  // Если бэкенд возвращает объект с users, total, limit, offset
  return response.data
  // Если бэкенд возвращает { users: [], total }, то используйте:
  // return {
  //   users: response.data.users || [],
  //   total: response.data.total || 0,
  //   limit: params?.limit || 10,
  //   offset: params?.offset || 0
  // }
}

export const adminUpdateUserStatus = async (userId: string, isActive: boolean): Promise<void> => {
  await client.put('/admin/users/status', { user_id: userId, is_active: isActive })
}

// УПРАВЛЕНИЕ ПРОДАВЦАМИ
export const adminGetSellers = async (params?: { verified?: boolean }): Promise<SellerWithShop[]> => {
  const response = await client.get('/admin/sellers', { params })
  return response.data
}

export const adminVerifySeller = async (shopId: string, verify: boolean): Promise<void> => {
  await client.put('/admin/sellers/verify', { shop_id: shopId, verify })
}

// УПРАВЛЕНИЕ МАГАЗИНОМ (ДЛЯ ПРОДАВЦА)
export const updateShopName = async (name: string): Promise<{ shop_id: string; shop_name: string }> => {
  const response = await client.put('/admin/shop', { name })
  return response.data
}

export const getShopInfo = async (): Promise<{ id: string; name: string; is_verified: boolean; rating: number }> => {
  const response = await client.get('/admin/shop')
  return response.data
}

// СТАТИСТИКА (АДМИН)
export const adminGetStats = async (params?: { period?: string }): Promise<any> => {
  const response = await client.get('/admin/stats', { params })
  return response.data
}