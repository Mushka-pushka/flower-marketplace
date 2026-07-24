import { useEffect, useState } from 'react'
import {
  FaBox,
  FaDollarSign,
  FaChartBar,
  FaShoppingCart,
  FaCalendarAlt,
} from 'react-icons/fa'
import {
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { useAuth } from '../context/AuthContext'
import { toast } from 'react-hot-toast'

// Типы данных
interface SellerAnalytics {
  total_orders: number
  total_revenue: number
  completed_orders: number
  cancelled_orders: number
  average_order_sum: number
}

interface StatusStats {
  status: string
  count: number
}

interface PopularProduct {
  product_id: string
  product_name: string
  total_sold: number
  total_revenue: number
  orders_count: number
}

interface SalesDay {
  date: string
  orders_count: number
  revenue: number
}

const SellerAnalyticsPage = () => {
  const { user } = useAuth()
  const [shopId, setShopId] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [analytics, setAnalytics] = useState<SellerAnalytics | null>(null)
  const [statusStats, setStatusStats] = useState<StatusStats[]>([])
  const [popularProducts, setPopularProducts] = useState<PopularProduct[]>([])
  const [salesDynamics, setSalesDynamics] = useState<SalesDay[]>([])
  const [period, setPeriod] = useState(30)

  // Цвета для диаграммы статусов
  const STATUS_COLORS: Record<string, string> = {
    pending: '#FBBF24',
    confirmed: '#3B82F6',
    preparing: '#8B5CF6',
    packing: '#6366F1',
    delivery: '#F97316',
    delivered: '#22C55E',
    cancelled: '#EF4444',
    paid: '#06B6D4',
  }

  const STATUS_LABELS: Record<string, string> = {
    pending: 'Ожидает',
    confirmed: 'Подтверждён',
    preparing: 'Готовится',
    packing: 'Упаковывается',
    delivery: 'В доставке',
    delivered: 'Доставлен',
    cancelled: 'Отменён',
    paid: 'Оплачен',
  }

  useEffect(() => {
    if (user?.shop_id) {
      setShopId(user.shop_id)
    } else {
      setShopId('66e67740-0bca-4634-8bf2-f8da3042c0dc')
    }
  }, [user])

  useEffect(() => {
    if (shopId) {
      fetchAnalytics()
    }
  }, [shopId, period])

  const fetchAnalytics = async () => {
    if (!shopId) return

    try {
      setLoading(true)

      const [analyticsRes, statusRes, popularRes, dynamicsRes] = await Promise.all([
        fetch(`http://localhost:8080/api/v1/analytics/seller?shop_id=${shopId}`, {
          headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
        }),
        fetch(`http://localhost:8080/api/v1/analytics/statuses?shop_id=${shopId}`, {
          headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
        }),
        fetch(`http://localhost:8080/api/v1/analytics/popular?shop_id=${shopId}&limit=5`, {
          headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
        }),
        fetch(`http://localhost:8080/api/v1/analytics/dynamics?shop_id=${shopId}&days=${period}`, {
          headers: { Authorization: `Bearer ${localStorage.getItem('access_token')}` },
        }),
      ])

      if (!analyticsRes.ok) throw new Error('Failed to fetch analytics')
      if (!statusRes.ok) throw new Error('Failed to fetch status stats')
      if (!popularRes.ok) throw new Error('Failed to fetch popular products')
      if (!dynamicsRes.ok) throw new Error('Failed to fetch sales dynamics')

      const analyticsData = await analyticsRes.json()
      const statusData = await statusRes.json()
      const popularData = await popularRes.json()
      const dynamicsData = await dynamicsRes.json()

      setAnalytics(analyticsData)
      setStatusStats(Array.isArray(statusData) ? statusData : [])
      setPopularProducts(Array.isArray(popularData) ? popularData : [])
      setSalesDynamics(Array.isArray(dynamicsData) ? dynamicsData : [])
    } catch (error) {
      console.error('Ошибка загрузки аналитики:', error)
      toast.error('Не удалось загрузить данные аналитики')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка аналитики...</div>
  }

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
        <FaChartBar className="text-[#8A9A86]" />
        Аналитика магазина
      </h2>

      {/* Период */}
      <div className="flex items-center gap-4">
        <span className="text-sm font-medium text-[#1C1C1C] flex items-center gap-1">
          <FaCalendarAlt className="text-[#8A9A86]" />
          Период:
        </span>
        <select
          value={period}
          onChange={(e) => setPeriod(Number(e.target.value))}
          className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-sm"
        >
          <option value={7}>7 дней</option>
          <option value={14}>14 дней</option>
          <option value={30}>30 дней</option>
          <option value={90}>90 дней</option>
        </select>
      </div>

      {/* Карточки статистики */}
      {analytics && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-[#8A9A86]/10 rounded-xl flex items-center justify-center">
                <FaShoppingCart className="text-[#8A9A86]" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Всего заказов</p>
                <p className="text-2xl font-bold text-[#1C1C1C]">{analytics.total_orders}</p>
              </div>
            </div>
          </div>
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-green-50 rounded-xl flex items-center justify-center">
                <FaDollarSign className="text-green-600" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Выручка</p>
                <p className="text-2xl font-bold text-[#1C1C1C]">{analytics.total_revenue} BYN</p>
              </div>
            </div>
          </div>
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-blue-50 rounded-xl flex items-center justify-center">
                <FaBox className="text-blue-600" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Доставлено</p>
                <p className="text-2xl font-bold text-[#1C1C1C]">{analytics.completed_orders}</p>
              </div>
            </div>
          </div>
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-purple-50 rounded-xl flex items-center justify-center">
                <FaChartBar className="text-purple-600" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Средний чек</p>
                <p className="text-2xl font-bold text-[#1C1C1C]">{analytics.average_order_sum} BYN</p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Графики */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Круговая диаграмма статусов */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4">Статусы заказов</h3>
          {statusStats.length > 0 ? (
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={statusStats}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={2}
                    dataKey="count"
                    nameKey="status"
                  >
                    {statusStats.map((entry) => (
                      <Cell
                        key={entry.status}
                        fill={STATUS_COLORS[entry.status] || '#9CA3AF'}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value, name) => [`${value} заказов`, STATUS_LABELS[name as string] || name]}
                  />
                </PieChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <p className="text-gray-400 text-sm">Нет данных по статусам</p>
          )}
        </div>

        {/* Популярные товары */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4">Популярные товары</h3>
          {popularProducts.length > 0 ? (
            <div className="space-y-3">
              {popularProducts.map((product, index) => (
                <div
                  key={product.product_id}
                  className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition"
                >
                  <span className="text-lg font-bold text-[#8A9A86] w-8">
                    #{index + 1}
                  </span>
                  <div className="flex-1">
                    <p className="font-medium text-[#1C1C1C] text-sm">
                      {product.product_name}
                    </p>
                    <p className="text-xs text-gray-400">
                      Продано: {product.total_sold} шт. • {product.orders_count} заказов
                    </p>
                  </div>
                  <span className="text-sm font-bold text-[#8A9A86]">
                    {product.total_revenue} BYN
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400 text-sm">Нет данных по популярным товарам</p>
          )}
        </div>
      </div>

      {/* График динамики продаж */}
      <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
        <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4">Динамика продаж</h3>
        {salesDynamics.length > 0 ? (
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={salesDynamics}>
                <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                <XAxis
                  dataKey="date"
                  tick={{ fontSize: 12 }}
                  tickFormatter={(date: string) => new Date(date).toLocaleDateString('ru-RU')}
                />
                <YAxis
                  yAxisId="left"
                  tick={{ fontSize: 12 }}
                  tickFormatter={(value: number) => `${value} BYN`}
                />
                <YAxis
                  yAxisId="right"
                  orientation="right"
                  tick={{ fontSize: 12 }}
                  tickFormatter={(value: number) => `${value} шт.`}
                />
                <Tooltip
                  formatter={(value, name) => {
                    if (name === 'revenue') return [`${value} BYN`, 'Выручка']
                    if (name === 'orders_count') return [`${value} заказов`, 'Заказы']
                    return [value, name]
                  }}
                  labelFormatter={(label) => {
                    if (typeof label === 'string') {
                     return new Date(label).toLocaleDateString('ru-RU')
                    }
                    return label
                   }}   
                />
                <Legend />
                <Line
                  yAxisId="left"
                  type="monotone"
                  dataKey="revenue"
                  stroke="#8A9A86"
                  strokeWidth={2}
                  name="revenue"
                  dot={{ r: 3 }}
                />
                <Line
                  yAxisId="right"
                  type="monotone"
                  dataKey="orders_count"
                  stroke="#3B82F6"
                  strokeWidth={2}
                  name="orders_count"
                  dot={{ r: 3 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        ) : (
          <p className="text-gray-400 text-sm">Нет данных по динамике продаж</p>
        )}
      </div>
    </div>
  )
}

export default SellerAnalyticsPage