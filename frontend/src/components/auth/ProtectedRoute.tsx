import { Navigate } from 'react-router-dom'
import { useEffect, useState } from 'react'

interface ProtectedRouteProps {
  children: React.ReactNode
  allowedRoles?: string[]
}

const ProtectedRoute = ({ children, allowedRoles }: ProtectedRouteProps) => {
  const [isValid, setIsValid] = useState<boolean | null>(null)

  useEffect(() => {
    const checkAuth = () => {
      try {
        const token = localStorage.getItem('access_token')
        if (!token) {
          setIsValid(false)
          return
        }

        // Проверяем срок действия токена
        try {
          const payload = JSON.parse(atob(token.split('.')[1]))
          const exp = payload.exp
          if (exp && Date.now() >= exp * 1000) {
            // Токен истек
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            localStorage.removeItem('user')
            setIsValid(false)
            return
          }
        } catch {
          setIsValid(false)
          return
        }

        // Проверяем роль
        if (allowedRoles && allowedRoles.length > 0) {
          const userStr = localStorage.getItem('user')
          if (!userStr) {
            setIsValid(false)
            return
          }
          try {
            const user = JSON.parse(userStr)
            if (!user.role || !allowedRoles.includes(user.role)) {
              setIsValid(false)
              return
            }
          } catch {
            setIsValid(false)
            return
          }
        }

        setIsValid(true)
      } catch {
        setIsValid(false)
      }
    }

    checkAuth()
  }, [allowedRoles])

  if (isValid === null) {
    return <div className="flex items-center justify-center min-h-screen">
      <div className="text-gray-400">Загрузка...</div>
    </div>
  }

  if (!isValid) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export default ProtectedRoute