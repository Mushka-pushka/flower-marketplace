import { FaStore } from 'react-icons/fa'

const AdminSellersPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaStore className="text-[#8A9A86]" />
        Модерация продавцов
      </h2>
      <p className="text-gray-400 text-base">Подтверждение и отклонение продавцов</p>
    </div>
  )
}

export default AdminSellersPage