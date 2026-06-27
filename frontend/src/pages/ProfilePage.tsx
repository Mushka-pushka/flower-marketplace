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
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">👤 Личный кабинет</h1>

      {/* Вкладки */}
      <div className="flex gap-2 mb-6 border-b pb-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as TabType)}
            className={`px-4 py-2 rounded-lg transition ${
              activeTab === tab.id
                ? 'bg-pink-500 text-white'
                : 'text-gray-600 hover:bg-gray-100'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Контент вкладок */}
      <div className="bg-white rounded-lg shadow p-6">
        {activeTab === 'orders' && <OrdersPage />}

        {activeTab === 'favorites' && (
          <div>
            <h2 className="text-xl font-semibold mb-4">❤️ Избранное</h2>
            <FavoritesList />
          </div>
        )}

        {activeTab === 'cart' && (
          <div>
            <h2 className="text-xl font-semibold mb-4">🛒 Корзина</h2>
            {items.length === 0 ? (
              <p className="text-gray-500">Корзина пуста</p>
            ) : (
              <div>
                {items.map((item) => (
                  <div key={item.id} className="flex justify-between py-2 border-b">
                    <span>{item.name}</span>
                    <span>{item.quantity} × {item.price} BYN</span>
                  </div>
                ))}
                <Link to="/cart" className="text-pink-500 hover:underline mt-4 inline-block">
                  Перейти в корзину →
                </Link>
              </div>
            )}
          </div>
        )}

        {activeTab === 'settings' && (
          <div>
            <h2 className="text-xl font-semibold mb-4">⚙️ Настройки профиля</h2>
            <p className="text-gray-500">Здесь будут настройки пользователя</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default ProfilePage