import client from './client'

// Типы для ответов
export interface User {
  id: string
  email: string
  phone: string
  first_name: string
  last_name: string
  role: 'customer' | 'seller' | 'admin'
  is_active: boolean
  created_at: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  phone?: string
  role?: 'customer' | 'seller'
}

// Регистрация
export const register = async (data: RegisterRequest): Promise<User> => {
  const response = await client.post('/auth/register', data)
  return response.data
}

// Вход
export const login = async (email: string, password: string): Promise<LoginResponse> => {
  const response = await client.post('/auth/login', { email, password })
  if (response.data.access_token) {
    localStorage.setItem('access_token', response.data.access_token)
    localStorage.setItem('refresh_token', response.data.refresh_token)
  }
  return response.data
}

// Выход
export const logout = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}

// Получение профиля
export const getProfile = async (): Promise<User> => {
  const response = await client.get('/auth/me')
  return response.data
}

// Обновление профиля
export const updateProfile = async (data: Partial<User>): Promise<User> => {
  const response = await client.put('/auth/profile', data)
  return response.data
}

// Смена пароля
export const changePassword = async (oldPassword: string, newPassword: string): Promise<void> => {
  await client.put('/auth/password', { old_password: oldPassword, new_password: newPassword })
}