import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import Login from './components/auth/Login'
import Register from './components/auth/Register'
import ProductPage from './pages/ProductPage'
//import ProtectedRoute from './components/auth/ProtectedRoute'
import './index.css'

function App() {
  const token = localStorage.getItem('access_token')

  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
            <Link to="/" className="text-2xl font-bold text-pink-600">
              Цветочный маркетплейс
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
                <button
                  onClick={() => {
                    localStorage.removeItem('access_token')
                    window.location.href = '/'
                  }}
                  className="text-gray-600 hover:text-pink-600"
                >
                  Выйти
                </button>
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
            <Route path="/product/:id" element={<ProductPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App