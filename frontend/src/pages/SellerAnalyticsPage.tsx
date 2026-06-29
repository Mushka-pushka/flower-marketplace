import { FaChartBar } from 'react-icons/fa'

const SellerAnalyticsPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaChartBar className="text-[#8A9A86]" />
        Аналитика
      </h2>
      <p className="text-gray-400 text-base">Количество заказов, выручка, популярные товары</p>
    </div>
  )
}

export default SellerAnalyticsPage