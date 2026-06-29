import { FaClipboardList } from 'react-icons/fa'

const SellerOrdersPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaClipboardList className="text-[#8A9A86]" />
        Заказы магазина
      </h2>
      <p className="text-gray-400 text-base">Список заказов вашего магазина</p>
    </div>
  )
}

export default SellerOrdersPage