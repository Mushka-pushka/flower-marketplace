export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  CATALOG: '/catalog',
  CART: '/cart',
  PROFILE: '/profile',
  CHECKOUT: '/checkout',
  CHECKOUT_SUCCESS: '/checkout/success',
  FAVORITES: '/favorites',
} as const

export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  REFRESH_TOKEN: 'refresh_token',
  USER: 'user',
} as const

export const CART_STORAGE_KEY = (userId?: string) => 
  userId ? `flower_cart_${userId}` : 'flower_cart_guest'

export const FAVORITES_STORAGE_KEY = (userId?: string) =>
  userId ? `flower_favorites_${userId}` : 'flower_favorites_guest'