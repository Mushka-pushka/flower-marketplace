import { Link } from 'react-router-dom'
import { FaSignOutAlt, FaStore, FaBox, FaLeaf, FaChartBar } from 'react-icons/fa'
import { useAuth } from '../../context/AuthContext'

const HeaderSeller = () => {
  const { logout } = useAuth()

  return (
    <header className="bg-white border-b border-gray-100 sticky top-0 z-40 shadow-[0_4px_20px_rgba(0,0,0,0.04)]">
      <div className="max-w-7xl mx-auto px-4 py-5 flex justify-between items-center">
        <div className="flex items-center gap-8">
          <Link to="/profile" className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
            <FaStore className="text-[#8A9A86] text-2xl" />
            Панель продавца
          </Link>
          <nav className="flex gap-6 items-center text-[#1C1C1C]">
            <Link to="/profile" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-sm font-medium">
              <FaBox /> Заказы
            </Link>
            <Link to="/profile?tab=seller-products" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-sm font-medium">
              <FaLeaf /> Товары
            </Link>
            <Link to="/profile?tab=seller-analytics" className="hover:text-[#8A9A86] transition flex items-center gap-2 text-sm font-medium">
              <FaChartBar /> Аналитика
            </Link>
          </nav>
        </div>
        <button onClick={logout} className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
          <FaSignOutAlt /> Выйти
        </button>
      </div>
    </header>
  )
}

export default HeaderSeller