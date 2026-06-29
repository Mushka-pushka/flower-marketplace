import { Link } from 'react-router-dom'
import { FaSignOutAlt, FaCog } from 'react-icons/fa'
import { useAuth } from '../../context/AuthContext'

const HeaderAdmin = () => {
  const { logout } = useAuth()

  return (
    <header className="bg-white border-b border-gray-100 sticky top-0 z-40 shadow-[0_4px_20px_rgba(0,0,0,0.04)]">
      <div className="max-w-7xl mx-auto px-4 py-5 flex justify-between items-center">
        <Link to="/profile" className="text-3xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaCog className="text-[#8A9A86] text-2xl" />
          Админ-панель
        </Link>
        <nav className="flex gap-8 items-center text-[#1C1C1C]">
          <button onClick={logout} className="hover:text-[#8A9A86] transition flex items-center gap-2 text-base font-medium">
            <FaSignOutAlt /> Выйти
          </button>
        </nav>
      </div>
    </header>
  )
}

export default HeaderAdmin