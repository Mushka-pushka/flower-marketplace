import { useState } from 'react'
import {
  FaBox,
  FaCheckCircle,
  FaLeaf,
  FaBoxOpen,
  FaTruck,
  FaGift,
  FaTimesCircle,
  FaArrowRight,
  FaBan,
} from 'react-icons/fa'
import { updateOrderStatus, cancelOrder } from '../api/order.api'
import { useAuth } from '../context/AuthContext'

interface StatusHistory {
  id: string
  status: string
  changed_by: string
  comment: string
  created_at: string
}

interface OrderTimelineProps {
  statuses: StatusHistory[]
  orderId: string
  currentStatus: string
  onStatusUpdate?: () => void
}

const OrderTimeline = ({ statuses, orderId, currentStatus, onStatusUpdate }: OrderTimelineProps) => {
  const { user } = useAuth()
  const [updating, setUpdating] = useState(false)

  // Строгая последовательность статусов
  const statusFlow = ['paid', 'confirmed', 'preparing', 'packing', 'delivery', 'delivered']
  const currentIndex = statusFlow.indexOf(currentStatus)

  // Следующий статус (только один шаг вперёд)
  const nextStatus = currentIndex >= 0 && currentIndex < statusFlow.length - 1
    ? statusFlow[currentIndex + 1]
    : null

  // Продавец может менять статус только если заказ оплачен и не завершён
  const canUpdate = user?.role === 'seller' && currentIndex >= 0 && nextStatus !== null

  // Продавец может отменить только до статуса "preparing"
  const canCancel = user?.role === 'seller' &&
    (currentStatus === 'paid' || currentStatus === 'confirmed')

  // Покупатель может отменить только до подтверждения продавцом
  const customerCanCancel = user?.role === 'customer' &&
    (currentStatus === 'pending' || currentStatus === 'paid')

  const getStatusLabel = (status: string) => {
    const map: Record<string, string> = {
      pending: 'Ожидает оплаты',
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

  const getStatusIcon = (status: string) => {
    const map: Record<string, React.ReactNode> = {
      pending: <FaBox className="text-yellow-500" />,
      paid: <FaCheckCircle className="text-blue-500" />,
      confirmed: <FaCheckCircle className="text-blue-600" />,
      preparing: <FaLeaf className="text-purple-500" />,
      packing: <FaBoxOpen className="text-indigo-500" />,
      delivery: <FaTruck className="text-orange-500" />,
      delivered: <FaGift className="text-green-500" />,
      cancelled: <FaTimesCircle className="text-red-500" />,
    }
    return map[status] || <FaBox className="text-gray-500" />
  }

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'border-yellow-400 bg-yellow-50 text-yellow-700',
      paid: 'border-blue-400 bg-blue-50 text-blue-700',
      confirmed: 'border-blue-500 bg-blue-50 text-blue-700',
      preparing: 'border-purple-400 bg-purple-50 text-purple-700',
      packing: 'border-indigo-400 bg-indigo-50 text-indigo-700',
      delivery: 'border-orange-400 bg-orange-50 text-orange-700',
      delivered: 'border-green-400 bg-green-50 text-green-700',
      cancelled: 'border-red-400 bg-red-50 text-red-700',
    }
    return map[status] || 'border-gray-400 bg-gray-50 text-gray-700'
  }

  const handleMoveToNext = async () => {
    if (!nextStatus) return
    if (!confirm(`Перевести заказ в статус «${getStatusLabel(nextStatus)}»?`)) return

    setUpdating(true)
    try {
      await updateOrderStatus({
        order_id: orderId,
        status: nextStatus,
        comment: `Статус изменён на «${getStatusLabel(nextStatus)}»`,
      })
      onStatusUpdate?.()
    } catch (error) {
      console.error('Ошибка обновления статуса:', error)
    } finally {
      setUpdating(false)
    }
  }

  const handleCancel = async () => {
    if (!confirm('Вы уверены, что хотите отменить заказ?')) return

    setUpdating(true)
    try {
      await cancelOrder(orderId)
      onStatusUpdate?.()
    } catch (error) {
      console.error('Ошибка отмены заказа:', error)
    } finally {
      setUpdating(false)
    }
  }

  const sorted = [...statuses].sort(
    (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  )

  return (
    <div>
      {/* Таймлайн истории статусов */}
      <div className="relative pl-8">
        <div className="absolute left-3 top-2 bottom-0 w-0.5 bg-gray-200" />

        {sorted.map((item, index) => (
          <div key={item.id} className={`relative mb-6 last:mb-0 ${index === 0 ? 'pt-0' : 'pt-4'}`}>
            <div className={`absolute -left-[22px] w-5 h-5 rounded-full border-2 shadow-sm flex items-center justify-center ${getStatusColor(item.status)} z-10`}>
              <span className="text-[10px]">{getStatusIcon(item.status)}</span>
            </div>

            <div className="ml-4">
              <div className="flex items-center gap-2">
                <span className="font-semibold text-[#1C1C1C]">{getStatusLabel(item.status)}</span>
              </div>
              <div className="text-sm text-gray-400 mt-0.5">
                {new Date(item.created_at).toLocaleString('ru-RU', {
                  day: '2-digit',
                  month: '2-digit',
                  year: 'numeric',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
                {item.changed_by === 'seller' && ' продавец'}
                {item.changed_by === 'system' && ' система'}
                {item.changed_by === 'customer' && ' покупатель'}
              </div>
              {item.comment && (
                <div className="text-sm text-gray-400 mt-0.5 italic bg-gray-50 px-3 py-1 rounded-full inline-block">
                  «{item.comment}»
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Панель управления статусом (только для продавца) */}
      {canUpdate && (
        <div className="mt-6 pt-4 border-t border-gray-200">
          <div className="bg-[#8A9A86]/5 rounded-xl p-4 border border-[#8A9A86]/20">
            <p className="text-sm text-gray-600 mb-3">
              Текущий статус: <span className="font-semibold text-[#1C1C1C]">{getStatusLabel(currentStatus)}</span>
            </p>
            <div className="flex flex-wrap gap-3">
              <button
                onClick={handleMoveToNext}
                disabled={updating}
                className="px-5 py-2.5 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50 flex items-center gap-2"
              >
                <FaArrowRight />
                {updating ? 'Обновление...' : `Перевести в «${getStatusLabel(nextStatus)}»`}
              </button>

              {canCancel && (
                <button
                  onClick={handleCancel}
                  disabled={updating}
                  className="px-5 py-2.5 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition text-sm font-medium border border-red-200 disabled:opacity-50 flex items-center gap-2"
                >
                  <FaBan />
                  Отменить заказ
                </button>
              )}
            </div>

            {!canCancel && (
              <p className="text-xs text-gray-400 mt-2 flex items-center gap-1">
                <FaBan className="text-[10px]" />
                Отмена недоступна после начала сборки заказа
              </p>
            )}
          </div>
        </div>
      )}

      {/* Кнопка отмены для покупателя */}
      {customerCanCancel && (
        <div className="mt-6 pt-4 border-t border-gray-200">
          <button
            onClick={handleCancel}
            disabled={updating}
            className="w-full px-5 py-2.5 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition text-sm font-medium border border-red-200 disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <FaBan />
            Отменить заказ
          </button>
        </div>
      )}
    </div>
  )
}

export default OrderTimeline