import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import {
  FaBox,
  FaTimes,
  FaShoppingBag,
} from 'react-icons/fa'
import { getMyOrders, getOrderDetails } from '../api/order.api'
import type { Order, OrderDetails } from '../api/order.api'
import OrderTimeline from '../components/OrderTimeline'
import { useAuth } from '../context/AuthContext'

const OrdersPage = () => {
  const { user } = useAuth()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedOrder, setSelectedOrder] = useState<OrderDetails | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)

  useEffect(() => {
    const fetchOrders = async () => {
      if (!user) {
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        const data = await getMyOrders(user.id)
        setOrders(Array.isArray(data) ? data : [])
      } catch (err: any) {
        console.error('Ошибка загрузки заказов:', err)
        setError(err.response?.data?.error || 'Ошибка загрузки заказов')
        setOrders([])
      } finally {
        setLoading(false)
      }
    }

    fetchOrders()
  }, [user])

  useEffect(() => {
    if (isModalOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'auto'
    }
    return () => {
      document.body.style.overflow = 'auto'
    }
  }, [isModalOpen])

  const handleOrderClick = async (orderId: string) => {
    try {
      const details = await getOrderDetails(orderId)
      setSelectedOrder(details)
      setIsModalOpen(true)
    } catch (err) {
      console.error('Ошибка загрузки деталей заказа:', err)
    }
  }

  const getStatusLabel = (status: string) => {
    const map: Record<string, string> = {
      pending: 'Ожидает подтверждения',
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
      confirmed: 'text-blue-600 bg-blue-50 border-blue-200',
      preparing: 'text-purple-600 bg-purple-50 border-purple-200',
      packing: 'text-indigo-600 bg-indigo-50 border-indigo-200',
      delivery: 'text-orange-600 bg-orange-50 border-orange-200',
      delivered: 'text-green-600 bg-green-50 border-green-200',
      cancelled: 'text-red-600 bg-red-50 border-red-200',
    }
    return map[status] || 'text-gray-600 bg-gray-50 border-gray-200'
  }

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
          Мои заказы
        </h2>
        <p className="text-gray-400 text-base">У вас пока нет заказов</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
        <FaBox className="text-[#8A9A86]" />
        Мои заказы
      </h2>
      <div className="space-y-3">
        {orders.map((order) => (
          <div
            key={order.id}
            onClick={() => handleOrderClick(order.id)}
            className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-4 cursor-pointer border border-gray-100"
          >
            <div className="flex justify-between items-start">
              <div>
                <p className="text-sm text-gray-400">Заказ #{order.id.slice(0, 8)}</p>
                <p className="text-sm text-gray-400">
                  {new Date(order.created_at).toLocaleDateString('ru-RU', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                  })}
                </p>
              </div>
              <div className="text-right">
                <span className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(order.current_status)}`}>
                  {getStatusLabel(order.current_status)}
                </span>
                <p className="text-lg font-bold text-[#8A9A86] mt-1">{order.total_amount} BYN</p>
              </div>
            </div>
          </div>
        ))}
      </div>

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
                <OrderTimeline statuses={selectedOrder.statuses} />
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
            </div>
          </div>,
          document.body
        )}
    </div>
  )
}

export default OrdersPage