import { Link } from 'react-router-dom'
import { FaHome, FaBook, FaUser, FaShoppingCart, FaSignOutAlt, FaLeaf } from 'react-icons/fa'
import { useCart } from '../../context/CartContext'
import { useAuth } from '../../context/AuthContext'

const HeaderCustomer = () => {
  const { user, logout } = useAuth()
  const { totalItems } = useCart()

  return (
    <header className="bg-white border-b border-gray-100 sticky top-0 z-40 shadow-[0_4px_20px_rgba(0,0,0,0.04)]">
      <div className="max-w-7xl mx-auto px-4 py-5 flex justify-between items-center">
        <Link to="/" className="text-3xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaLeaf className="text-[#8A9A86] text-2xl" />
          Цветочный маркетплейс
        </Link>
        <nav className="flex gap-8 items-center text-[#1C1C1C]">
          <Link to="/" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
            <FaHome /> Главная
          </Link>
          <Link to="/catalog" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
            <FaBook /> Каталог
          </Link>

          {user ? (
            <>
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
      </div>
    </header>
  )
}

export default HeaderCustomer