import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import Login from './components/auth/Login'
import Register from './components/auth/Register'
import CartPage from './pages/CartPage'
import ProfilePage from './pages/ProfilePage'
import CheckoutPage from './pages/CheckoutPage'
import CheckoutSuccess from './pages/CheckoutSuccess'
import { CartProvider } from './context/CartContext'
import { FavoritesProvider } from './context/FavoritesContext'
import { useCart } from './context/CartContext'
import './index.css'

function AppContent() {
  const token = localStorage.getItem('access_token')
  const { totalItems } = useCart()

  return (
    <div className="min-h-screen bg-gradient-to-br from-pink-50 via-rose-50 to-purple-50">
      <header className="bg-white/80 backdrop-blur-md shadow-sm border-b border-pink-100 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <Link to="/" className="text-3xl font-bold gradient-text">
            🌸 Цветочный маркетплейс
          </Link>
          <nav className="flex gap-6 items-center text-gray-700">
            <Link to="/" className="hover:text-pink-500 transition">Главная</Link>
            <Link to="/catalog" className="hover:text-pink-500 transition">Каталог</Link>
            {!token ? (
              <>
                <Link to="/login" className="hover:text-pink-500 transition">Вход</Link>
                <Link to="/register" className="hover:text-pink-500 transition">Регистрация</Link>
              </>
            ) : (
              <>
                <Link to="/profile" className="hover:text-pink-500 transition text-xl" title="Профиль">👤</Link>
                <Link to="/cart" className="relative hover:text-pink-500 transition text-xl" title="Корзина">
                  🛒
                  {totalItems > 0 && (
                    <span className="absolute -top-2 -right-3 bg-pink-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                      {totalItems}
                    </span>
                  )}
                </Link>
                <button
                  onClick={() => {
                    localStorage.removeItem('access_token')
                    window.location.href = '/'
                  }}
                  className="hover:text-pink-500 transition"
                >
                  Выйти
                </button>
              </>
            )}
          </nav>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/catalog" element={<Catalog />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/cart" element={<CartPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/checkout" element={<CheckoutPage />} />
          <Route path="/checkout/success" element={<CheckoutSuccess />} />
        </Routes>
      </main>
    </div>
  )
}

function App() {
  return (
    <BrowserRouter>
      <CartProvider>
        <FavoritesProvider>
          <AppContent />
        </FavoritesProvider>
      </CartProvider>
    </BrowserRouter>
  )
}

export default App