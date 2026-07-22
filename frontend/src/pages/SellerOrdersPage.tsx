import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import {
  FaBox,
  FaTimes,
  FaShoppingBag,
  FaFilter,
  FaSearch,
} from 'react-icons/fa'
import { getShopOrders, updateOrderStatus } from '../api/order.api'
import type { Order, OrderDetails } from '../api/order.api'
import { getOrderDetails } from '../api/order.api'
import OrderTimeline from '../components/OrderTimeline'
import { useAuth } from '../context/AuthContext'
import { toast } from 'react-hot-toast'

const SellerOrdersPage = () => {
  const { user } = useAuth()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedOrder, setSelectedOrder] = useState<OrderDetails | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [shopId, setShopId] = useState<string | null>(null)

  useEffect(() => {
    if (user?.shop_id) {
      setShopId(user.shop_id)
    }
  }, [user])

  const fetchOrders = async () => {
    if (!shopId) return

    try {
      setLoading(true)
      const data = await getShopOrders(shopId)
      setOrders(Array.isArray(data) ? data : [])
    } catch (err: any) {
      console.error('Ошибка загрузки заказов магазина:', err)
      setError(err.response?.data?.error || 'Ошибка загрузки заказов')
      setOrders([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchOrders()
  }, [shopId])

  const handleOrderClick = async (orderId: string) => {
    try {
      const details = await getOrderDetails(orderId)
      setSelectedOrder(details)
      setIsModalOpen(true)
    } catch (err) {
      console.error('Ошибка загрузки деталей заказа:', err)
      toast.error('Не удалось загрузить детали заказа')
    }
  }

  const handleUpdateStatus = async (orderId: string, newStatus: string) => {
    if (!confirm(`Изменить статус заказа на "${getStatusLabel(newStatus)}"?`)) return

    try {
      await updateOrderStatus({
        order_id: orderId,
        status: newStatus,
        comment: `Статус изменён на ${getStatusLabel(newStatus)}`,
      })
      toast.success('Статус заказа обновлён')
      await fetchOrders()
      if (selectedOrder && selectedOrder.order.id === orderId) {
        const updated = await getOrderDetails(orderId)
        setSelectedOrder(updated)
      }
    } catch (error) {
      console.error('Ошибка обновления статуса:', error)
      toast.error('Не удалось обновить статус')
    }
  }

  const getStatusLabel = (status: string) => {
    const map: Record<string, string> = {
      pending: 'Ожидает подтверждения',
      paid: 'Оплачен',
      confirmed: 'Подтверждён',
      preparing: 'Собирается',
      packing: 'Упаковывается',
      delivery: 'В доставке',
      delivered: 'Доставлен',
      cancelled: 'Отменён',
    }
    return map[status] || status
  }

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'text-yellow-600 bg-yellow-50 border-yellow-200',
      paid: 'text-blue-600 bg-blue-50 border-blue-200',
      confirmed: 'text-blue-600 bg-blue-50 border-blue-200',
      preparing: 'text-purple-600 bg-purple-50 border-purple-200',
      packing: 'text-indigo-600 bg-indigo-50 border-indigo-200',
      delivery: 'text-orange-600 bg-orange-50 border-orange-200',
      delivered: 'text-green-600 bg-green-50 border-green-200',
      cancelled: 'text-red-600 bg-red-50 border-red-200',
    }
    return map[status] || 'text-gray-600 bg-gray-50 border-gray-200'
  }

  const getStatusFlow = (status: string): string[] => {
    const flow: Record<string, string[]> = {
      pending: ['confirmed', 'preparing', 'packing', 'delivery', 'delivered'],
      confirmed: ['preparing', 'packing', 'delivery', 'delivered'],
      preparing: ['packing', 'delivery', 'delivered'],
      packing: ['delivery', 'delivered'],
      delivery: ['delivered'],
      delivered: [],
      cancelled: [],
      paid: ['confirmed', 'preparing', 'packing', 'delivery', 'delivered'],
    }
    return flow[status] || []
  }

  // Фильтрация с поиском
  const filteredOrders = orders.filter((order) => {
    // Фильтр по статусу
    if (statusFilter !== 'all' && order.current_status !== statusFilter) return false
    
    // Поиск
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      const orderId = order.id.toLowerCase()
      const customerName = `${order.customer_first_name || ''} ${order.customer_last_name || ''}`.toLowerCase()
      const customerEmail = (order.customer_email || '').toLowerCase()
      const productNames = (order.product_names || '').toLowerCase() 
      
      return orderId.includes(query) || 
             customerName.includes(query) || 
             customerEmail.includes(query) ||
             productNames.includes(query) 
    }
    
    return true
  })

  const statusOptions = [
    { value: 'all', label: 'Все' },
    { value: 'pending', label: 'Ожидают' },
    { value: 'paid', label: 'Оплачены' },
    { value: 'confirmed', label: 'Подтверждены' },
    { value: 'preparing', label: 'Готовятся' },
    { value: 'packing', label: 'Упаковываются' },
    { value: 'delivery', label: 'В доставке' },
    { value: 'delivered', label: 'Доставлены' },
    { value: 'cancelled', label: 'Отменены' },
  ]

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка заказов...</div>
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">{error}</div>
  }

  if (!orders || orders.length === 0) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
          <FaBox className="text-[#8A9A86]" />
          Заказы магазина
        </h2>
        <p className="text-gray-400 text-base">У вашего магазина пока нет заказов</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
        <FaBox className="text-[#8A9A86]" />
        Заказы магазина
      </h2>

      {/* Фильтр и поиск */}
      <div className="flex flex-wrap items-center gap-3 mb-4">
        <div className="flex items-center gap-2">
          <FaFilter className="text-[#8A9A86]" />
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          >
            {statusOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>

        {/* Поиск */}
        <div className="flex-1 min-w-[200px] relative">
          <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Поиск по номеру заказа или покупателю..."
            className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          />
        </div>

        <span className="text-sm text-gray-400 ml-auto whitespace-nowrap">
          Найдено: {filteredOrders.length} из {orders.length} заказов
        </span>
      </div>

      {/* Список заказов */}
      <div className="space-y-3">
        {filteredOrders.map((order) => {
          const nextStatuses = getStatusFlow(order.current_status)

          return (
            <div
              key={order.id}
              onClick={() => handleOrderClick(order.id)}
              className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-4 cursor-pointer border border-gray-100"
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-3 flex-wrap">
                    <p className="text-sm font-medium text-[#1C1C1C]">
                      Заказ #{order.id.slice(0, 8)}
                    </p>
                    <span className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(order.current_status)}`}>
                      {getStatusLabel(order.current_status)}
                    </span>
                  </div>
                  <p className="text-sm text-gray-400 mt-1">
                    {order.customer_first_name} {order.customer_last_name} • {order.customer_email}
                  </p>
                  <p className="text-sm text-gray-400">
                    {new Date(order.created_at).toLocaleDateString('ru-RU', {
                      day: '2-digit',
                      month: '2-digit',
                      year: 'numeric',
                    })} • Сумма: <span className="font-bold text-[#8A9A86]">{order.total_amount} BYN</span>
                  </p>
                  <p className="text-sm text-gray-400">
                    Товары: <span className="text-[#1C1C1C] font-medium">
                      {order.product_names ? order.product_names.split(', ').slice(0, 3).join(', ') + (order.product_names.split(',').length > 3 ? '...' : '') : '—'}
                    </span>
                  </p>
                </div>

                <div className="flex flex-col items-end gap-2">
                  {/* Кнопки смены статуса */}
                  {nextStatuses.length > 0 && (
                    <div className="flex flex-wrap gap-1 justify-end">
                      {nextStatuses.slice(0, 3).map((status) => (
                        <button
                          key={status}
                          onClick={(e) => {
                            e.stopPropagation()
                            handleUpdateStatus(order.id, status)
                          }}
                          className="px-3 py-1 text-xs bg-[#8A9A86] text-white rounded-lg hover:bg-[#7A8A76] transition whitespace-nowrap"
                        >
                          {getStatusLabel(status)}
                        </button>
                      ))}
                      {nextStatuses.length > 3 && (
                        <span className="text-xs text-gray-400 self-center">
                          +{nextStatuses.length - 3}
                        </span>
                      )}
                    </div>
                  )}
                  {order.current_status === 'delivered' && (
                    <span className="text-xs text-green-600 font-medium">
                      Завершён
                    </span>
                  )}
                  {order.current_status === 'cancelled' && (
                    <span className="text-xs text-red-600 font-medium">
                      Отменён
                    </span>
                  )}
                </div>
              </div>
            </div>
          )
        })}

        {filteredOrders.length === 0 && (
          <p className="text-center text-gray-400 py-8">
            Заказы не найдены. Попробуйте изменить фильтры.
          </p>
        )}
      </div>

      {/* Модалка деталей заказа */}
      {isModalOpen &&
        selectedOrder &&
        createPortal(
          <div
            className="fixed inset-0 z-[9999] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
            onClick={() => setIsModalOpen(false)}
          >
            <div
              className="bg-white rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 shadow-[0_8px_40px_rgba(0,0,0,0.08)] border border-gray-100 animate-fade-in-up"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="flex justify-between items-center mb-4 pb-3 border-b border-gray-100">
                <h2 className="text-xl font-bold text-[#1C1C1C] flex items-center gap-2">
                  <FaShoppingBag className="text-[#8A9A86]" />
                  Заказ #{selectedOrder.order.id.slice(0, 8)}
                </h2>
                <button
                  onClick={() => setIsModalOpen(false)}
                  className="text-gray-400 hover:text-gray-600 text-2xl leading-none transition flex-shrink-0"
                >
                  <FaTimes />
                </button>
              </div>

              <div className="mb-4">
                <p className="text-sm text-gray-400">
                  Покупатель: <span className="font-medium text-[#1C1C1C]">
                    {selectedOrder.customer_first_name} {selectedOrder.customer_last_name}
                  </span>
                </p>
                <p className="text-sm text-gray-400">
                  Email: <span className="font-medium text-[#1C1C1C]">{selectedOrder.customer_email}</span>
                </p>
                <p className="text-sm text-gray-400">
                  Сумма: <span className="font-bold text-[#8A9A86]">{selectedOrder.order.total_amount} BYN</span>
                </p>
                <p className="text-sm text-gray-400">
                  Статус:{' '}
                  <span
                    className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(
                      selectedOrder.order.current_status
                    )}`}
                  >
                    {getStatusLabel(selectedOrder.order.current_status)}
                  </span>
                </p>
              </div>

              <div className="border-t border-gray-100 pt-4">
                <h3 className="font-semibold text-[#1C1C1C] mb-3">История статусов</h3>
                <OrderTimeline
                  statuses={selectedOrder.statuses}
                  orderId={selectedOrder.order.id}
                  currentStatus={selectedOrder.order.current_status}
                  onStatusUpdate={fetchOrders}
                />
              </div>

              <div className="border-t border-gray-100 pt-4 mt-4">
                <h3 className="font-semibold text-[#1C1C1C] mb-2">Товары</h3>
                <div className="space-y-2">
                  {selectedOrder.items.map((item) => (
                    <div
                      key={item.id}
                      className="flex justify-between text-sm border-b border-gray-100 pb-2 last:border-0"
                    >
                      <span className="text-gray-400">
                        {item.quantity} × {item.price} BYN
                      </span>
                      <span className="text-[#1C1C1C] font-medium">{item.total} BYN</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Быстрые кнопки смены статуса в модалке */}
              {selectedOrder.order.current_status !== 'delivered' &&
                selectedOrder.order.current_status !== 'cancelled' && (
                  <div className="border-t border-gray-100 pt-4 mt-4">
                    <h4 className="text-sm font-medium text-[#1C1C1C] mb-2">Изменить статус:</h4>
                    <div className="flex flex-wrap gap-2">
                      {getStatusFlow(selectedOrder.order.current_status).map((status) => (
                        <button
                          key={status}
                          onClick={() => {
                            handleUpdateStatus(selectedOrder.order.id, status)
                            setIsModalOpen(false)
                          }}
                          className="px-4 py-2 text-sm bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition"
                        >
                          {getStatusLabel(status)}
                        </button>
                      ))}
                    </div>
                  </div>
                )}
            </div>
          </div>,
          document.body
        )}
    </div>
  )
}

export default SellerOrdersPage