import { useEffect, useState } from 'react'
import {
  FaChartLine,
  FaShoppingCart,
  FaUsers,
  FaStore,
  FaMoneyBillWave,
  FaClock,
  FaCheckCircle,
  FaTimesCircle,
  FaSpinner,
  FaDollarSign,
  FaLeaf,
  FaBox,
  FaCalendarAlt,
  FaUserCog,
  FaPercent,
  FaRocket,
} from 'react-icons/fa'
import { adminGetStats, type AdminStats } from '../api/admin.api'
import { toast } from 'react-hot-toast'
import { useNavigate } from 'react-router-dom'

const AdminStatsPage = () => {
  const navigate = useNavigate()
  const [stats, setStats] = useState<AdminStats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const data = await adminGetStats()
      console.log('Stats:', data)
      setStats(data)
    } catch (error) {
      console.error('Ошибка:', error)
      toast.error('Не удалось загрузить статистику')
    } finally {
      setLoading(false)
    }
  }

  const formatCurrency = (value: any): string => {
    if (value === undefined || value === null || isNaN(value)) {
      return '0'
    }
    return Number(value).toLocaleString('ru-RU', {
      minimumFractionDigits: 0,
      maximumFractionDigits: 2,
    })
  }

  const safeValue = (value: any) => {
    return value !== undefined && value !== null ? value : 0
  }

  const getPercent = (value: number, total: number) => {
    if (total === 0) return 0
    return Math.round((value / total) * 100)
  }

  // Функция для перехода по вкладкам
  const goToTab = (tab: string) => {
    navigate(`/profile?tab=${tab}`)
    window.location.reload()
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
        <p className="text-gray-400">Нет данных</p>
      </div>
    )
  }

  const inProgressOrders = 
    (stats.orders_by_status?.pending || 0) +
    (stats.orders_by_status?.confirmed || 0) +
    (stats.orders_by_status?.preparing || 0) +
    (stats.orders_by_status?.packing || 0) +
    (stats.orders_by_status?.delivery || 0)

  const averageOrder = stats.total_orders > 0 
    ? Math.round(stats.total_revenue / stats.total_orders) 
    : 0

  const deliveredPercent = stats.total_orders > 0
    ? getPercent(stats.orders_by_status?.delivered || 0, stats.total_orders)
    : 0

  const usersByRole = stats.users_by_role || {}

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
            <FaChartLine className="text-[#8A9A86]" />
            Общая статистика
          </h2>
          <p className="text-gray-400 text-sm flex items-center gap-2">
            <FaCalendarAlt className="text-xs" />
            Актуально на {new Date().toLocaleDateString('ru-RU', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            })}
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 bg-green-50 text-green-700 px-4 py-2 rounded-xl border border-green-200">
            <FaRocket className="text-green-600" />
            <span className="text-sm font-medium">Платформа активна</span>
          </div>
        </div>
      </div>

      {/* Основные KPI — 4 главных метрики */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition group cursor-default">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-50 rounded-xl flex items-center justify-center group-hover:scale-110 transition">
              <FaDollarSign className="text-green-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{formatCurrency(stats.total_revenue)}</p>
              <p className="text-sm text-gray-400">Общая выручка</p>
              <p className="text-xs text-gray-400 mt-0.5">BYN</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition group cursor-default">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-purple-50 rounded-xl flex items-center justify-center group-hover:scale-110 transition">
              <FaMoneyBillWave className="text-purple-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{formatCurrency(stats.platform_revenue)}</p>
              <p className="text-sm text-gray-400">Комиссия платформы</p>
              <p className="text-xs text-gray-400 mt-0.5">10% от выручки</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition group cursor-default">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-50 rounded-xl flex items-center justify-center group-hover:scale-110 transition">
              <FaShoppingCart className="text-blue-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{safeValue(stats.total_orders)}</p>
              <p className="text-sm text-gray-400">Всего заказов</p>
              <p className="text-xs text-gray-400 mt-0.5">
                {deliveredPercent}% доставлено
              </p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition group cursor-default">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-amber-50 rounded-xl flex items-center justify-center group-hover:scale-110 transition">
              <FaClock className="text-amber-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{averageOrder}</p>
              <p className="text-sm text-gray-400">Средний чек</p>
              <p className="text-xs text-gray-400 mt-0.5">BYN</p>
            </div>
          </div>
        </div>
      </div>

      {/* Вторая строка — пользователи, магазины, товары, категории */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-[#8A9A86]/10 rounded-xl flex items-center justify-center">
              <FaUsers className="text-[#8A9A86] text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{safeValue(stats.total_users)}</p>
              <p className="text-sm text-gray-400">Пользователей</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-orange-50 rounded-xl flex items-center justify-center">
              <FaStore className="text-orange-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{safeValue(stats.total_shops)}</p>
              <p className="text-sm text-gray-400">Магазинов</p>
              <p className="text-xs text-gray-400">{safeValue(stats.verified_shops)} верифицировано</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-emerald-50 rounded-xl flex items-center justify-center">
              <FaLeaf className="text-emerald-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{safeValue(stats.total_products)}</p>
              <p className="text-sm text-gray-400">Товаров</p>
              <p className="text-xs text-gray-400">{safeValue(stats.active_products)} активно</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-indigo-50 rounded-xl flex items-center justify-center">
              <FaBox className="text-indigo-600 text-xl" />
            </div>
            <div>
              <p className="text-2xl font-bold text-[#1C1C1C]">{safeValue(stats.total_categories)}</p>
              <p className="text-sm text-gray-400">Категорий</p>
            </div>
          </div>
        </div>
      </div>

      {/* Детали по заказам + прогресс-бар */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 col-span-2">
          <h3 className="text-sm font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
            <FaShoppingCart className="text-[#8A9A86]" />
            Статусы заказов
          </h3>
          <div className="space-y-3">
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 flex items-center gap-1.5">
                  <FaCheckCircle className="text-green-500 text-xs" />
                  Доставлено
                </span>
                <span className="font-medium text-green-600">
                  {safeValue(stats.orders_by_status?.delivered)} ({deliveredPercent}%)
                </span>
              </div>
              <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-green-500 rounded-full transition-all duration-500"
                  style={{ width: `${deliveredPercent}%` }}
                />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 flex items-center gap-1.5">
                  <FaClock className="text-yellow-500 text-xs" />
                  В обработке
                </span>
                <span className="font-medium text-yellow-600">{inProgressOrders}</span>
              </div>
              <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-yellow-500 rounded-full transition-all duration-500"
                  style={{ width: `${getPercent(inProgressOrders, stats.total_orders)}%` }}
                />
              </div>
            </div>
            <div>
              <div className="flex justify-between text-sm mb-1">
                <span className="text-gray-600 flex items-center gap-1.5">
                  <FaTimesCircle className="text-red-500 text-xs" />
                  Отменено
                </span>
                <span className="font-medium text-red-600">{safeValue(stats.orders_by_status?.cancelled)}</span>
              </div>
              <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-red-500 rounded-full transition-all duration-500"
                  style={{ width: `${getPercent(stats.orders_by_status?.cancelled || 0, stats.total_orders)}%` }}
                />
              </div>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100">
          <h3 className="text-sm font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
            <FaPercent className="text-[#8A9A86]" />
            Распределение пользователей
          </h3>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 flex items-center gap-1.5">
                <FaUsers className="text-[#8A9A86] text-xs" />
                Покупатели
              </span>
              <span className="font-medium text-[#1C1C1C]">{safeValue(usersByRole.customer)}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 flex items-center gap-1.5">
                <FaStore className="text-orange-500 text-xs" />
                Продавцы
              </span>
              <span className="font-medium text-[#1C1C1C]">{safeValue(usersByRole.seller)}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-600 flex items-center gap-1.5">
                <FaUserCog className="text-purple-500 text-xs" />
                Админы
              </span>
              <span className="font-medium text-[#1C1C1C]">{safeValue(usersByRole.admin)}</span>
            </div>
            <div className="pt-3 border-t border-gray-100 flex justify-between items-center">
              <span className="text-sm font-medium text-gray-600">Всего</span>
              <span className="font-bold text-[#1C1C1C]">{safeValue(stats.total_users)}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Карточка здоровья платформы — минималистичная */}
      <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-[#8A9A86]/10 rounded-xl flex items-center justify-center">
              <FaChartLine className="text-[#8A9A86] text-lg" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-[#1C1C1C]">Здоровье платформы</h3>
              <p className="text-xs text-gray-400">Ключевые показатели эффективности</p>
            </div>
          </div>
          <div className="flex flex-wrap items-center gap-8">
            <div className="text-center">
              <p className="text-xl font-bold text-[#1C1C1C]">{getPercent(deliveredPercent, 100)}%</p>
              <p className="text-xs text-gray-400">Доставка</p>
            </div>
            <div className="text-center">
              <div className="flex justify-center">
                {stats.total_orders > 0 ? (
                  <FaCheckCircle className="text-xl text-green-500" />
                ) : (
                  <FaClock className="text-xl text-yellow-500" />
                )}
              </div>
              <p className="text-xs text-gray-400 mt-0.5">
                {stats.total_orders > 0 ? 'Активно' : 'Ожидание'}
              </p>
            </div>
            <div className="text-center">
              <p className="text-xl font-bold text-[#1C1C1C]">{formatCurrency(stats.total_revenue)}</p>
              <p className="text-xs text-gray-400">Выручка</p>
            </div>
            <button
              onClick={() => goToTab('admin-sellers')}
              className="flex items-center gap-2 bg-[#8A9A86] text-white px-4 py-2 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium"
            >
              <FaStore className="text-xs" />
              Управление
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default AdminStatsPage