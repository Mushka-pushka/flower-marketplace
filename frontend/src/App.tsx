import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import Catalog from './pages/Catalog'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import './index.css'

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        {/* Шапка с навигацией */}
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
            <Link to="/" className="text-2xl font-bold text-pink-600">
              Цветочный маркетплейс
            </Link>
            <nav className="flex gap-6">
              <Link to="/" className="text-gray-600 hover:text-pink-600">Главная</Link>
              <Link to="/catalog" className="text-gray-600 hover:text-pink-600">Каталог</Link>
              <Link to="/login" className="text-gray-600 hover:text-pink-600">Вход</Link>
              <Link to="/register" className="text-gray-600 hover:text-pink-600">Регистрация</Link>
            </nav>
          </div>
        </header>

        {/* Основной контент */}
        <main className="max-w-7xl mx-auto px-4 py-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/catalog" element={<Catalog />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  )
}

export default App