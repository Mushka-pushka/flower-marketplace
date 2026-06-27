import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useCart } from '../context/CartContext'
import OrdersPage from './OrdersPage'
import FavoritesList from '../components/FavoritesList'

type TabType = 'orders' | 'favorites' | 'cart' | 'settings'

const ProfilePage = () => {
  const [activeTab, setActiveTab] = useState<TabType>('orders')
  const { items } = useCart()

  const tabs = [
    { id: 'orders', label: '📦 Заказы' },
    { id: 'favorites', label: '❤️ Избранное' },
    { id: 'cart', label: `🛒 Корзина (${items.length})` },
    { id: 'settings', label: '⚙️ Настройки' },
  ]

  return (
    <div className="max-w-4xl mx-auto animate-fade-in-up">
      <h1 className="text-4xl font-bold gradient-text mb-6">👤 Личный кабинет</h1>

      {/* Вкладки */}
      <div className="flex flex-wrap gap-2 mb-6 border-b border-pink-100 pb-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as TabType)}
            className={`px-5 py-2.5 rounded-full transition-all duration-300 font-medium ${
              activeTab === tab.id
                ? 'btn-primary shadow-lg'
                : 'text-gray-600 hover:bg-pink-50 hover:text-pink-600'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Контент вкладок */}
      <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-pink-50/50">
        {activeTab === 'orders' && <OrdersPage />}

        {activeTab === 'favorites' && <FavoritesList />}

        {activeTab === 'cart' && (
          <div>
            <h2 className="text-2xl font-semibold gradient-text mb-4">🛒 Корзина</h2>
            {items.length === 0 ? (
              <p className="text-gray-400 text-lg">Корзина пуста</p>
            ) : (
              <div className="space-y-2">
                {items.map((item) => (
                  <div key={item.id} className="flex justify-between py-2 border-b border-gray-100">
                    <span className="text-gray-700">{item.name}</span>
                    <span className="text-pink-600 font-medium">
                      {item.quantity} × {item.price} BYN
                    </span>
                  </div>
                ))}
                <Link
                  to="/cart"
                  className="text-pink-500 hover:text-pink-700 font-medium mt-4 inline-block transition"
                >
                  Перейти в корзину →
                </Link>
              </div>
            )}
          </div>
        )}

        {activeTab === 'settings' && (
          <div>
            <h2 className="text-2xl font-semibold gradient-text mb-4">⚙️ Настройки профиля</h2>
            <p className="text-gray-400 text-lg">Здесь будут настройки пользователя</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default ProfilePage