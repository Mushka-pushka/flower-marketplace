import {
  FaBox,
  FaCheckCircle,
  FaLeaf,
  FaBoxOpen,
  FaTruck,
  FaGift,
  FaTimesCircle,
} from 'react-icons/fa'

interface StatusHistory {
  id: string
  status: string
  changed_by: string
  comment: string
  created_at: string
}

interface OrderTimelineProps {
  statuses: StatusHistory[]
}

const OrderTimeline = ({ statuses }: OrderTimelineProps) => {
  const getStatusLabel = (status: string) => {
    const map: Record<string, string> = {
      pending: 'Заказ создан',
      confirmed: 'Продавец подтвердил',
      preparing: 'Букет собирается',
      packing: 'Заказ упаковывается',
      delivery: 'Передан курьеру',
      delivered: 'Доставлен получателю',
      cancelled: 'Заказ отменён',
    }
    return map[status] || status
  }

  const getStatusIcon = (status: string) => {
    const map: Record<string, React.ReactNode> = {
      pending: <FaBox className="text-yellow-500" />,
      confirmed: <FaCheckCircle className="text-blue-500" />,
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
      confirmed: 'border-blue-400 bg-blue-50 text-blue-700',
      preparing: 'border-purple-400 bg-purple-50 text-purple-700',
      packing: 'border-indigo-400 bg-indigo-50 text-indigo-700',
      delivery: 'border-orange-400 bg-orange-50 text-orange-700',
      delivered: 'border-green-400 bg-green-50 text-green-700',
      cancelled: 'border-red-400 bg-red-50 text-red-700',
    }
    return map[status] || 'border-gray-400 bg-gray-50 text-gray-700'
  }

  const sorted = [...statuses].sort(
    (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  )

  return (
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
  )
}

export default OrderTimeline