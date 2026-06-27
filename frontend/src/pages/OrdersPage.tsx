import { useEffect, useState } from 'react'
import { getMyOrders, getOrderDetails } from '../api/order.api'
import type { Order, OrderDetails } from '../api/order.api'
import OrderTimeline from '../components/OrderTimeline'

const OrdersPage = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedOrder, setSelectedOrder] = useState<OrderDetails | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true)
        const data = await getMyOrders('6b75b13b-2b7b-4df1-b700-b39ac0bc1d45')
        setOrders(data)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Ошибка загрузки заказов')
      } finally {
        setLoading(false)
      }
    }

    fetchOrders()
  }, [])

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
      pending: '⏳ Ожидает подтверждения',
      confirmed: '✅ Подтверждён',
      preparing: '🌸 Собирается',
      packing: '📦 Упаковывается',
      delivery: '🚚 В доставке',
      delivered: '🎉 Доставлен',
      cancelled: '❌ Отменён',
    }
    return map[status] || status
  }

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'text-yellow-600 bg-yellow-50',
      confirmed: 'text-blue-600 bg-blue-50',
      preparing: 'text-purple-600 bg-purple-50',
      packing: 'text-indigo-600 bg-indigo-50',
      delivery: 'text-orange-600 bg-orange-50',
      delivered: 'text-green-600 bg-green-50',
      cancelled: 'text-red-600 bg-red-50',
    }
    return map[status] || 'text-gray-600 bg-gray-50'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-500">Загрузка заказов...</div>
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">{error}</div>
  }

  if (orders.length === 0) {
    return (
      <div>
        <h2 className="text-xl font-semibold mb-4">📦 Мои заказы</h2>
        <p className="text-gray-500">У вас пока нет заказов</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">📦 Мои заказы</h2>
      <div className="space-y-4">
        {orders.map((order) => (
          <div
            key={order.id}
            onClick={() => handleOrderClick(order.id)}
            className="border rounded-lg p-4 bg-white shadow-sm cursor-pointer hover:shadow-md transition"
          >
            <div className="flex justify-between items-start">
              <div>
                <p className="text-sm text-gray-500">Заказ #{order.id.slice(0, 8)}</p>
                <p className="text-sm text-gray-500">
                  {new Date(order.created_at).toLocaleDateString('ru-RU')}
                </p>
              </div>
              <div className="text-right">
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(order.current_status)}`}>
                  {getStatusLabel(order.current_status)}
                </span>
                <p className="text-lg font-bold text-pink-600 mt-1">{order.total_amount} BYN</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Модальное окно с деталями заказа */}
      {isModalOpen && selectedOrder && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
          onClick={() => setIsModalOpen(false)}
        >
          <div
            className="bg-white rounded-2xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-2xl font-bold text-gray-800">
                Заказ #{selectedOrder.order.id.slice(0, 8)}
              </h2>
              <button
                onClick={() => setIsModalOpen(false)}
                className="text-gray-400 hover:text-gray-600 text-3xl leading-none"
              >
                ×
              </button>
            </div>

            <div className="mb-4">
              <p className="text-sm text-gray-500">
                Сумма: <span className="font-bold text-pink-600">{selectedOrder.order.total_amount} BYN</span>
              </p>
              <p className="text-sm text-gray-500">
                Статус: <span className={`px-2 py-0.5 rounded-full text-sm ${getStatusColor(selectedOrder.order.current_status)}`}>
                  {getStatusLabel(selectedOrder.order.current_status)}
                </span>
              </p>
            </div>

            <div className="border-t pt-4">
              <h3 className="font-semibold text-gray-700 mb-3">📜 История статусов</h3>
              <OrderTimeline statuses={selectedOrder.statuses} />
            </div>

            <div className="border-t pt-4 mt-4">
              <h3 className="font-semibold text-gray-700 mb-2">🛍️ Товары</h3>
              <div className="space-y-2">
                {selectedOrder.items.map((item) => (
                  <div key={item.id} className="flex justify-between text-sm border-b pb-2">
                    <span className="text-gray-600">
                      {item.quantity} × {item.price} BYN
                    </span>
                    <span className="text-gray-800 font-medium">{item.total} BYN</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default OrdersPage