import { Navigate } from 'react-router-dom'

interface ProtectedRouteProps {
  children: React.ReactNode
  requireRole?: 'customer' | 'seller' | 'admin'
}

const ProtectedRoute = ({ children, requireRole }: ProtectedRouteProps) => {
  const token = localStorage.getItem('access_token')
  const user = JSON.parse(localStorage.getItem('user') || 'null')

  if (!token) {
    return <Navigate to="/login" replace />
  }

  if (requireRole && user?.role !== requireRole) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}

export default ProtectedRoute