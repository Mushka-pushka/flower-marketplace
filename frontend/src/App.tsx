import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import { CartProvider } from './context/CartContext'
import { FavoritesProvider } from './context/FavoritesContext'
import { useAuth } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'

// Шапки
import HeaderCustomer from './components/Header/HeaderCustomer'
import HeaderSeller from './components/Header/HeaderSeller'
import HeaderAdmin from './components/Header/HeaderAdmin'

// Страницы
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import Login from './components/auth/Login'
import Register from './components/auth/Register'
import CartPage from './pages/CartPage'
import ProfilePage from './pages/ProfilePage'
import CheckoutPage from './pages/CheckoutPage'
import CheckoutSuccess from './pages/CheckoutSuccess'
import './index.css'

function AppContent() {
  const { user } = useAuth()

  const renderHeader = () => {
    if (user?.role === 'admin') return <HeaderAdmin />
    if (user?.role === 'seller') return <HeaderSeller />
    return <HeaderCustomer />
  }

  const renderHome = () => {
    if (user?.role === 'admin') return <Navigate to="/profile" replace />
    if (user?.role === 'seller') return <Navigate to="/profile" replace />
    return <Home />
  }

  return (
    <div className="min-h-screen bg-white">
      {renderHeader()}

      <main className="max-w-7xl mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={renderHome()} />
          <Route path="/catalog" element={<Catalog />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          <Route path="/cart" element={
            <ProtectedRoute>
              <CartPage />
            </ProtectedRoute>
          } />
          <Route path="/profile" element={
            <ProtectedRoute>
              <ProfilePage />
            </ProtectedRoute>
          } />
          <Route path="/checkout" element={
            <ProtectedRoute>
              <CheckoutPage />
            </ProtectedRoute>
          } />
          <Route path="/checkout/success" element={
            <ProtectedRoute>
              <CheckoutSuccess />
            </ProtectedRoute>
          } />
        </Routes>
      </main>
    </div>
  )
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <CartProvider>
          <FavoritesProvider>
            <AppContent />
          </FavoritesProvider>
        </CartProvider>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App