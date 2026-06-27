import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import Login from './components/auth/Login'
import Register from './components/auth/Register'
import CartPage from './pages/CartPage'
import ProfilePage from './pages/ProfilePage'
import { CartProvider } from './context/CartContext'
import { useCart } from './context/CartContext'
import './index.css'

function AppContent() {
  const token = localStorage.getItem('access_token')
  const { totalItems } = useCart()

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <Link to="/" className="text-2xl font-bold text-pink-600">
            🌸 Цветочный маркетплейс
          </Link>
          <nav className="flex gap-6 items-center">
            <Link to="/" className="text-gray-600 hover:text-pink-600">Главная</Link>
            <Link to="/catalog" className="text-gray-600 hover:text-pink-600">Каталог</Link>
            {!token ? (
              <>
                <Link to="/login" className="text-gray-600 hover:text-pink-600">Вход</Link>
                <Link to="/register" className="text-gray-600 hover:text-pink-600">Регистрация</Link>
              </>
            ) : (
              <>
                <Link to="/profile" className="text-gray-600 hover:text-pink-600" title="Профиль">
                  👤
                </Link>
                <Link to="/cart" className="relative text-gray-600 hover:text-pink-600" title="Корзина">
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
                  className="text-gray-600 hover:text-pink-600"
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
        </Routes>
      </main>
    </div>
  )
}

function App() {
  return (
    <BrowserRouter>
      <CartProvider>
        <AppContent />
      </CartProvider>
    </BrowserRouter>
  )
}

export default App