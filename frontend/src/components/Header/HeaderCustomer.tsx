import { useState } from 'react'
import { Link } from 'react-router-dom'
import { 
  FaHome, FaBook, FaUser, FaShoppingCart, FaSignOutAlt, FaLeaf, 
  FaBars, FaTimes, FaHeart 
} from 'react-icons/fa'
import { useCart } from '../../context/CartContext'
import { useAuth } from '../../context/AuthContext'

const HeaderCustomer = () => {
  const { user, logout } = useAuth()
  const { totalItems } = useCart()
  const [isMenuOpen, setIsMenuOpen] = useState(false)

  const toggleMenu = () => setIsMenuOpen(!isMenuOpen)
  const closeMenu = () => setIsMenuOpen(false)

  return (
    <header className="bg-white border-b border-gray-100 sticky top-0 z-40 shadow-[0_4px_20px_rgba(0,0,0,0.04)]">
      <div className="max-w-7xl mx-auto px-4 py-5 flex justify-between items-center">
        {/* Логотип */}
        <Link to="/" className="text-2xl md:text-3xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaLeaf className="text-[#8A9A86] text-2xl" />
          <span className="hidden sm:inline">Цветочный маркетплейс</span>
          <span className="sm:hidden">Цветы</span>
        </Link>

        {/* Десктопное меню */}
        <nav className="hidden md:flex gap-8 items-center text-[#1C1C1C]">
          <Link to="/" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
            <FaHome /> Главная
          </Link>
          <Link to="/catalog" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
            <FaBook /> Каталог
          </Link>

          {user ? (
            <>
              <Link to="/favorites" className="hover:text-[#8A9A86] transition text-xl" title="Избранное">
                <FaHeart />
              </Link>
              <Link to="/profile" className="hover:text-[#8A9A86] transition text-xl" title="Профиль">
                <FaUser />
              </Link>
              <Link to="/cart" className="relative hover:text-[#8A9A86] transition text-xl" title="Корзина">
                <FaShoppingCart />
                {totalItems > 0 && (
                  <span className="absolute -top-2 -right-3 bg-[#8A9A86] text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                    {totalItems}
                  </span>
                )}
              </Link>
              <button onClick={logout} className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
                <FaSignOutAlt /> Выйти
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="hover:text-[#8A9A86] transition text-base font-medium">Вход</Link>
              <Link to="/register" className="hover:text-[#8A9A86] transition text-base font-medium">Регистрация</Link>
            </>
          )}
        </nav>

        {/* Бургер-меню для мобильных */}
        <button 
          onClick={toggleMenu}
          className="md:hidden text-2xl text-[#1C1C1C] hover:text-[#8A9A86] transition"
        >
          {isMenuOpen ? <FaTimes /> : <FaBars />}
        </button>
      </div>

      {/* Мобильное меню */}
      {isMenuOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white p-4 space-y-3 animate-fade-in-up">
          <Link to="/" onClick={closeMenu} className="flex items-center gap-3 text-[#1C1C1C] hover:text-[#8A9A86] transition py-2">
            <FaHome /> Главная
          </Link>
          <Link to="/catalog" onClick={closeMenu} className="flex items-center gap-3 text-[#1C1C1C] hover:text-[#8A9A86] transition py-2">
            <FaBook /> Каталог
          </Link>
          {user ? (
            <>
              <Link to="/favorites" onClick={closeMenu} className="flex items-center gap-3 text-[#1C1C1C] hover:text-[#8A9A86] transition py-2">
                <FaHeart /> Избранное
              </Link>
              <Link to="/profile" onClick={closeMenu} className="flex items-center gap-3 text-[#1C1C1C] hover:text-[#8A9A86] transition py-2">
                <FaUser /> Профиль
              </Link>
              <Link to="/cart" onClick={closeMenu} className="flex items-center gap-3 text-[#1C1C1C] hover:text-[#8A9A86] transition py-2">
                <FaShoppingCart /> Корзина {totalItems > 0 && `(${totalItems})`}
              </Link>
              <button 
                onClick={() => { logout(); closeMenu(); }} 
                className="flex items-center gap-3 text-red-500 hover:text-red-700 transition py-2 w-full text-left"
              >
                <FaSignOutAlt /> Выйти
              </button>
            </>
          ) : (
            <>
              <Link to="/login" onClick={closeMenu} className="flex items-center gap-3 text-[#8A9A86] hover:underline transition py-2">
                Войти
              </Link>
              <Link to="/register" onClick={closeMenu} className="flex items-center gap-3 text-[#8A9A86] hover:underline transition py-2">
                Зарегистрироваться
              </Link>
            </>
          )}
        </div>
      )}
    </header>
  )
}

export default HeaderCustomer