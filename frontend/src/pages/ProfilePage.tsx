import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import {
  FaBox,
  FaHeart,
  FaShoppingCart,
  FaCog,
  FaClipboardList,
  FaLeaf,
  FaChartBar,
  FaStore,
  FaUsers,
  FaChartLine,
  FaFolder,
  FaUser,
} from 'react-icons/fa'
import { useCart } from '../context/CartContext'
import { useAuth } from '../context/AuthContext'
import OrdersPage from './OrdersPage'
import FavoritesList from '../components/FavoritesList'
import SellerOrdersPage from './SellerOrdersPage'
import SellerProductsPage from './SellerProductsPage'
import SellerAnalyticsPage from './SellerAnalyticsPage'
import AdminSellersPage from './AdminSellersPage'
import AdminUsersPage from './AdminUsersPage'
import AdminStatsPage from './AdminStatsPage'
import AdminCategoriesPage from './AdminCategoriesPage'
import ProfileSettings from './ProfileSettings'

type TabType =
  | 'orders'
  | 'favorites'
  | 'cart'
  | 'settings'
  | 'seller-orders'
  | 'seller-products'
  | 'seller-analytics'
  | 'admin-sellers'
  | 'admin-users'
  | 'admin-stats'
  | 'admin-categories'

const ProfilePage = () => {
  const { user } = useAuth()
  const { items } = useCart()

  const getDefaultTab = (): TabType => {
    if (user?.role === 'admin') return 'admin-sellers'
    if (user?.role === 'seller') return 'seller-orders'
    return 'orders'
  }

  const [activeTab, setActiveTab] = useState<TabType>(getDefaultTab())

  useEffect(() => {
    // ✅ Если активная вкладка — "settings", НЕ сбрасываем её
    if (activeTab !== 'settings') {
      setActiveTab(getDefaultTab())
    }
  }, [user])

  const getTabs = () => {
    if (user?.role === 'seller') {
      return [
        { id: 'seller-orders', label: 'Заказы магазина', icon: <FaClipboardList /> },
        { id: 'seller-products', label: 'Товары', icon: <FaLeaf /> },
        { id: 'seller-analytics', label: 'Аналитика', icon: <FaChartBar /> },
        { id: 'settings', label: 'Настройки', icon: <FaCog /> },
      ]
    }

    if (user?.role === 'admin') {
      return [
        { id: 'admin-sellers', label: 'Продавцы', icon: <FaStore /> },
        { id: 'admin-users', label: 'Пользователи', icon: <FaUsers /> },
        { id: 'admin-stats', label: 'Статистика', icon: <FaChartLine /> },
        { id: 'admin-categories', label: 'Категории', icon: <FaFolder /> },
        { id: 'settings', label: 'Настройки', icon: <FaCog /> },
      ]
    }

    return [
      { id: 'orders', label: 'Заказы', icon: <FaBox /> },
      { id: 'favorites', label: 'Избранное', icon: <FaHeart /> },
      { id: 'cart', label: `Корзина (${items.length})`, icon: <FaShoppingCart /> },
      { id: 'settings', label: 'Настройки', icon: <FaCog /> },
    ]
  }

  const tabs = getTabs()

  const renderContent = () => {
    switch (activeTab) {
      case 'orders':
        return <OrdersPage />
      case 'favorites':
        return <FavoritesList />
      case 'cart':
        return (
          <div>
            <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
              <FaShoppingCart className="text-[#8A9A86]" />
              Корзина
            </h2>
            {items.length === 0 ? (
              <p className="text-gray-400 text-base">Корзина пуста</p>
            ) : (
              <>
                {items.map((item) => (
                  <div key={item.id} className="flex justify-between py-2 border-b border-gray-100">
                    <span className="text-[#1C1C1C]">{item.name}</span>
                    <span className="text-[#8A9A86] font-medium">
                      {item.quantity} × {item.price} BYN
                    </span>
                  </div>
                ))}
                <Link to="/cart" className="text-[#8A9A86] hover:text-[#7A8A76] font-medium mt-4 inline-block transition">
                  Перейти в корзину →
                </Link>
              </>
            )}
          </div>
        )
      case 'settings':
        return <ProfileSettings />
      case 'seller-orders':
        return <SellerOrdersPage />
      case 'seller-products':
        return <SellerProductsPage />
      case 'seller-analytics':
        return <SellerAnalyticsPage />
      case 'admin-sellers':
        return <AdminSellersPage />
      case 'admin-users':
        return <AdminUsersPage />
      case 'admin-stats':
        return <AdminStatsPage />
      case 'admin-categories':
        return <AdminCategoriesPage />
      default:
        return null
    }
  }

  return (
    <div className="max-w-4xl mx-auto animate-fade-in-up">
      <h1 className="text-3xl font-bold text-[#1C1C1C] mb-6 flex items-center gap-2">
        <FaUser className="text-[#8A9A86]" />
        Личный кабинет
      </h1>

      <div className="flex flex-wrap gap-2 mb-6 border-b border-gray-100 pb-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as TabType)}
            className={`px-5 py-2.5 rounded-xl transition-all duration-300 font-medium flex items-center gap-2 text-sm ${
              activeTab === tab.id
                ? 'bg-[#8A9A86] text-white shadow-[0_4px_20px_rgba(138,154,134,0.3)]'
                : 'text-[#1C1C1C] hover:bg-gray-50'
            }`}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
        {renderContent()}
      </div>
    </div>
  )
}

export default ProfilePage