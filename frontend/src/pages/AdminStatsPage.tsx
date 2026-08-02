import { useEffect, useState } from 'react'
import {
  FaChartLine,
  FaShoppingCart,
  FaUsers,
  FaStore,
  FaMoneyBillWave,
  FaTrophy,
  FaLeaf,
  FaClock,
  FaCheckCircle,
  FaTimesCircle,
  FaSpinner,
} from 'react-icons/fa'
import { adminGetStats } from '../api/admin.api'
import { toast } from 'react-hot-toast'

interface Stats {
  total_revenue: number
  total_orders: number
  total_users: number
  total_sellers: number
  average_order: number
  platform_commission: number
  pending_orders: number
  delivered_orders: number
  cancelled_orders: number
  top_sellers: Array<{
    shop_name: string
    revenue: number
    orders_count: number
  }>
  top_products: Array<{
    name: string
    total_sold: number
    revenue: number
  }>
  recent_orders: Array<{
    id: string
    customer_name: string
    total_amount: number
    status: string
    created_at: string
  }>
}

const AdminStatsPage = () => {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)
  const [timeRange, setTimeRange] = useState<'today' | 'week' | 'month' | 'all'>('all')

  useEffect(() => {
    fetchStats()
  }, [timeRange])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const data = await adminGetStats({ period: timeRange })
      setStats(data)
    } catch (error) {
      console.error('Ошибка загрузки статистики:', error)
      toast.error('Не удалось загрузить статистику')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <FaSpinner className="text-3xl text-[#8A9A86] animate-spin mx-auto" />
        <p className="text-gray-400 mt-3">Загрузка статистики...</p>
      </div>
    )
  }

  if (!stats) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-400">Нет данных для отображения</p>
      </div>
    )
  }

  // Форматирование валюты
  const formatCurrency = (amount: number) => {
    return amount.toLocaleString('ru-RU', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaChartLine className="text-[#8A9A86]" />
          Общая статистика
        </h2>

        {/* Фильтр по времени */}
        <div className="flex gap-2">
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as typeof timeRange)}
            className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          >
            <option value="today">Сегодня</option>
            <option value="week">Неделя</option>
            <option value="month">Месяц</option>
            <option value="all">Всё время</option>
          </select>
        </div>
      </div>

      {/* Основные метрики */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-4 mb-6">
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaMoneyBillWave className="text-[#8A9A86] text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">
              {formatCurrency(stats.total_revenue)} BYN
            </p>
          </div>
          <p className="text-sm text-gray-400">Общая выручка</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaShoppingCart className="text-blue-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.total_orders}</p>
          </div>
          <p className="text-sm text-gray-400">Всего заказов</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaUsers className="text-purple-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.total_users}</p>
          </div>
          <p className="text-sm text-gray-400">Пользователей</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaStore className="text-green-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.total_sellers}</p>
          </div>
          <p className="text-sm text-gray-400">Продавцов</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaClock className="text-orange-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">
              {formatCurrency(stats.average_order)} BYN
            </p>
          </div>
          <p className="text-sm text-gray-400">Средний чек</p>
        </div>
      </div>

      {/* Статусы заказов и комиссия */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaCheckCircle className="text-green-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.delivered_orders}</p>
          </div>
          <p className="text-sm text-gray-400">Выполнено</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaClock className="text-yellow-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.pending_orders}</p>
          </div>
          <p className="text-sm text-gray-400">В обработке</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaTimesCircle className="text-red-500 text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">{stats.cancelled_orders}</p>
          </div>
          <p className="text-sm text-gray-400">Отменено</p>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
          <div className="flex items-center gap-2">
            <FaMoneyBillWave className="text-[#8A9A86] text-lg" />
            <p className="text-2xl font-bold text-[#1C1C1C]">
              {formatCurrency(stats.platform_commission)} BYN
            </p>
          </div>
          <p className="text-sm text-gray-400">Комиссия платформы</p>
        </div>
      </div>

      {/* Топ продавцов и топ товаров */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        {/* Топ продавцов */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
            <FaTrophy className="text-yellow-500" />
            Топ продавцов
          </h3>
          {stats.top_sellers && stats.top_sellers.length > 0 ? (
            <div className="space-y-3">
              {stats.top_sellers.slice(0, 5).map((seller, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className={`text-sm font-bold ${
                      index === 0 ? 'text-yellow-500' :
                      index === 1 ? 'text-gray-400' :
                      index === 2 ? 'text-amber-600' :
                      'text-gray-300'
                    }`}>
                      #{index + 1}
                    </span>
                    <span className="font-medium text-[#1C1C1C]">{seller.shop_name}</span>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-bold text-[#8A9A86]">
                      {formatCurrency(seller.revenue)} BYN
                    </p>
                    <p className="text-xs text-gray-400">{seller.orders_count} заказов</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">Нет данных</p>
          )}
        </div>

        {/* Топ товаров */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
            <FaLeaf className="text-[#8A9A86]" />
            Популярные товары
          </h3>
          {stats.top_products && stats.top_products.length > 0 ? (
            <div className="space-y-3">
              {stats.top_products.slice(0, 5).map((product, index) => (
                <div key={index} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-gray-400">#{index + 1}</span>
                    <span className="font-medium text-[#1C1C1C]">{product.name}</span>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-bold text-[#8A9A86]">
                      {formatCurrency(product.revenue)} BYN
                    </p>
                    <p className="text-xs text-gray-400">{product.total_sold} шт.</p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-center py-4">Нет данных</p>
          )}
        </div>
      </div>

      {/* Недавние заказы */}
      <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
        <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
          <FaClock className="text-[#8A9A86]" />
          Последние заказы
        </h3>
        {stats.recent_orders && stats.recent_orders.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-100">
                <tr>
                  <th className="text-left p-2 font-medium text-[#1C1C1C]">Заказ</th>
                  <th className="text-left p-2 font-medium text-[#1C1C1C]">Покупатель</th>
                  <th className="text-left p-2 font-medium text-[#1C1C1C]">Сумма</th>
                  <th className="text-left p-2 font-medium text-[#1C1C1C]">Статус</th>
                  <th className="text-left p-2 font-medium text-[#1C1C1C]">Дата</th>
                </tr>
              </thead>
              <tbody>
                {stats.recent_orders.map((order) => (
                  <tr key={order.id} className="border-b border-gray-50 hover:bg-gray-50/50 transition">
                    <td className="p-2 font-mono text-xs text-gray-500">
                      #{order.id.slice(0, 8)}
                    </td>
                    <td className="p-2 text-[#1C1C1C]">{order.customer_name}</td>
                    <td className="p-2 font-bold text-[#8A9A86]">
                      {formatCurrency(order.total_amount)} BYN
                    </td>
                    <td className="p-2">
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                        order.status === 'delivered' ? 'bg-green-50 text-green-600' :
                        order.status === 'cancelled' ? 'bg-red-50 text-red-600' :
                        order.status === 'pending' ? 'bg-yellow-50 text-yellow-600' :
                        'bg-blue-50 text-blue-600'
                      }`}>
                        {order.status === 'delivered' ? 'Доставлен' :
                         order.status === 'cancelled' ? 'Отменён' :
                         order.status === 'pending' ? 'Ожидает' :
                         order.status}
                      </span>
                    </td>
                    <td className="p-2 text-gray-500 text-xs">
                      {new Date(order.created_at).toLocaleDateString('ru-RU', {
                        day: '2-digit',
                        month: '2-digit',
                        year: 'numeric',
                      })}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="text-gray-400 text-center py-4">Нет недавних заказов</p>
        )}
      </div>
    </div>
  )
}

export default AdminStatsPage