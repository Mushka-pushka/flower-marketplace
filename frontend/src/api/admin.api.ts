import client from './client'

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

// УПРАВЛЕНИЕ ПОЛЬЗОВАТЕЛЯМИ

// Получение списка пользователей
export const adminGetUsers = async (params?: {
  search?: string
  role?: string
  limit?: number
  offset?: number
}): Promise<User[]> => {
  const response = await client.get('/admin/users/list', { params })
  return response.data.users || []
}

// Обновление статуса пользователя
export const adminUpdateUserStatus = async (userId: string, isActive: boolean): Promise<void> => {
  await client.put('/admin/users/status', { user_id: userId, is_active: isActive })
}

// УПРАВЛЕНИЕ ПРОДАВЦАМИ

// Получение списка продавцов
export const adminGetSellers = async (params?: { verified?: boolean }): Promise<any[]> => {
  const response = await client.get('/admin/sellers', { params })
  return response.data
}

// Верификация продавца
export const adminVerifySeller = async (shopId: string, verify: boolean): Promise<void> => {
  await client.put('/admin/sellers/verify', { shop_id: shopId, verify })
}

// УПРАВЛЕНИЕ МАГАЗИНОМ (ДЛЯ ПРОДАВЦА)

// Обновление названия магазина
export const updateShopName = async (name: string): Promise<{ shop_id: string; shop_name: string }> => {
    const response = await client.put('/admin/shop', { name })
    return response.data
}

// Получение информации о магазине
export const getShopInfo = async (): Promise<{ id: string; name: string; is_verified: boolean; rating: number }> => {
    const response = await client.get('/admin/shop')
    return response.data
}

// СТАТИСТИКА (АДМИН)

// Получение статистики
export const adminGetStats = async (): Promise<any> => {
  const response = await client.get('/admin/stats')
  return response.data
}