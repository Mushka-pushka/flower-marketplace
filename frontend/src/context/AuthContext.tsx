import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { getProfile } from '../api/auth.api'  

interface User {
  id: string
  email: string
  role: 'customer' | 'seller' | 'admin'
  first_name: string
  last_name: string
  phone: string
  avatar_url?: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (user: User, token: string) => void
  logout: () => void
  isLoading: boolean
  updateUser: (userData: Partial<User>) => void  
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const loadUser = async () => {
      const storedToken = localStorage.getItem('access_token')
      const storedUser = localStorage.getItem('user')
      
      if (storedToken && storedUser) {
        try {
          // ✅ Загружаем свежие данные с сервера
          const freshUser = await getProfile()
          setUser(freshUser)
          setToken(storedToken)
          // Обновляем localStorage свежими данными
          localStorage.setItem('user', JSON.stringify(freshUser))
        } catch (error) {
          console.error('Failed to load user profile:', error)
          // Если ошибка — используем данные из localStorage
          try {
            setUser(JSON.parse(storedUser))
            setToken(storedToken)
          } catch {
            localStorage.removeItem('user')
            localStorage.removeItem('access_token')
          }
        }
      }
      setIsLoading(false)
    }

    loadUser()
  }, [])

  const login = (userData: User, tokenData: string) => {
    setUser(userData)
    setToken(tokenData)
    localStorage.setItem('user', JSON.stringify(userData))
    localStorage.setItem('access_token', tokenData)
  }

  const logout = () => {
    setUser(null)
    setToken(null)
    localStorage.removeItem('user')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  const updateUser = (userData: Partial<User>) => {
    if (user) {
      const updatedUser = { ...user, ...userData }
      setUser(updatedUser)
      localStorage.setItem('user', JSON.stringify(updatedUser))
    }
  }

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isLoading, updateUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}