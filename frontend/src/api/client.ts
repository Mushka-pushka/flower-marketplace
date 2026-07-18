import axios from 'axios'
import { toast } from 'react-hot-toast'
import { API_BASE_URL, STORAGE_KEYS, ROUTES } from '../constants'

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Флаг для предотвращения множественных запросов refresh
let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

const subscribeTokenRefresh = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb)
}

const onTokenRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

// Глобальная обработка ошибок
const showError = (message: string) => {
  toast.error(message)
}

// Интерсептор для добавления JWT-токена
client.interceptors.request.use((config) => {
  const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Интерсептор для обработки ошибок и refresh токена
client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // Если ошибка 401 и это не запрос на refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Если уже идет обновление, добавляем в очередь
        return new Promise((resolve) => {
          subscribeTokenRefresh((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            resolve(client(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const refreshToken = localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN)
        if (!refreshToken) {
          throw new Error('No refresh token')
        }

        // Запрос на обновление токена
        const response = await axios.post(
          `${API_BASE_URL}/auth/refresh`,
          { refresh_token: refreshToken }
        )

        const { access_token, refresh_token } = response.data
        localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, access_token)
        localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refresh_token)

        // Обновляем все ожидающие запросы
        onTokenRefreshed(access_token)

        // Повторяем оригинальный запрос
        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return client(originalRequest)
      } catch (refreshError) {
        // Если refresh не удался - logout
        localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
        localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
        localStorage.removeItem(STORAGE_KEYS.USER)
        showError('Сессия истекла. Пожалуйста, войдите заново.')
        window.location.href = ROUTES.LOGIN
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // Показываем уведомление для не-401 ошибок
    if (error.response?.status !== 401) {
      const message = error.response?.data?.error || 
                      error.response?.data?.message || 
                      'Произошла ошибка. Пожалуйста, попробуйте позже.'
      showError(message)
    }

    // Для других 401 ошибок (не связанных с refresh)
    if (error.response?.status === 401) {
      localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN)
      localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN)
      localStorage.removeItem(STORAGE_KEYS.USER)
      window.location.href = ROUTES.LOGIN
    }

    return Promise.reject(error)
  }
)

export default client