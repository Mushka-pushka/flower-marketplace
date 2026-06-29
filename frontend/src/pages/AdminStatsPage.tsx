import { FaChartLine } from 'react-icons/fa'

const AdminStatsPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaChartLine className="text-[#8A9A86]" />
        Общая статистика
      </h2>
      <p className="text-gray-400 text-base">Заказы, пользователи, выручка платформы</p>
    </div>
  )
}

export default AdminStatsPage