import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'

interface User {
  id: string
  email: string
  role: 'customer' | 'seller' | 'admin'
  first_name: string
  last_name: string
  phone: string
}

// Добавляем тип для сохранения в localStorage
interface StoredUser {
  id: string
  email: string
  role: 'customer' | 'seller' | 'admin'
  first_name: string
  last_name: string
  phone: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  login: (user: User, token: string) => void
  logout: () => void
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    // Загружаем данные из localStorage при старте
    const storedUser = localStorage.getItem('user')
    const storedToken = localStorage.getItem('access_token')
    
    if (storedUser && storedToken) {
      try {
        const parsed: StoredUser = JSON.parse(storedUser)
        setUser(parsed)
        setToken(storedToken)
      } catch {
        localStorage.removeItem('user')
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
      }
    }
    setIsLoading(false)
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

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth must be used within AuthProvider')
  return context
}