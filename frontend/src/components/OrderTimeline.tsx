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
    const map: Record<string, string> = {
      pending: '📦',
      confirmed: '✅',
      preparing: '🌸',
      packing: '📦',
      delivery: '🚚',
      delivered: '🎉',
      cancelled: '❌',
    }
    return map[status] || '📌'
  }

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'border-yellow-400 bg-yellow-50',
      confirmed: 'border-blue-400 bg-blue-50',
      preparing: 'border-purple-400 bg-purple-50',
      packing: 'border-indigo-400 bg-indigo-50',
      delivery: 'border-orange-400 bg-orange-50',
      delivered: 'border-green-400 bg-green-50',
      cancelled: 'border-red-400 bg-red-50',
    }
    return map[status] || 'border-gray-400 bg-gray-50'
  }

  // Сортируем по времени (от старых к новым)
  const sorted = [...statuses].sort(
    (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  )

  return (
    <div className="relative pl-6">
      {/* Вертикальная линия */}
      <div className="absolute left-2 top-2 bottom-0 w-0.5 bg-gray-300" />

      {sorted.map((item, index) => (
        <div key={item.id} className={`relative mb-6 last:mb-0 ${index === 0 ? 'pt-0' : 'pt-4'}`}>
          {/* Точка на линии */}
          <div className={`absolute -left-[22px] w-4 h-4 rounded-full border-2 ${getStatusColor(item.status)} z-10`} />

          <div className="ml-4">
            <div className="flex items-center gap-2">
              <span className="text-xl">{getStatusIcon(item.status)}</span>
              <span className="font-medium text-gray-800">{getStatusLabel(item.status)}</span>
            </div>
            <div className="text-sm text-gray-500 mt-0.5">
              {new Date(item.created_at).toLocaleString('ru-RU')}
              {item.changed_by === 'seller' && ' (продавец)'}
              {item.changed_by === 'system' && ' (система)'}
              {item.changed_by !== 'system' && item.changed_by !== 'seller' && ` (${item.changed_by})`}
            </div>
            {item.comment && (
              <div className="text-sm text-gray-400 mt-0.5 italic">«{item.comment}»</div>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}

export default OrderTimeline