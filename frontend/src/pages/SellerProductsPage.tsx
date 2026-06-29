import { FaLeaf } from 'react-icons/fa'

const SellerProductsPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaLeaf className="text-[#8A9A86]" />
        Управление товарами
      </h2>
      <p className="text-gray-400 text-base">Добавление, редактирование, удаление товаров</p>
    </div>
  )
}

export default SellerProductsPage