import { useState } from 'react'
import { Link } from 'react-router-dom'
import { FaHeart, FaTrash, FaLeaf } from 'react-icons/fa'
import { useFavorites } from '../context/FavoritesContext'
import ProductModal from './ProductModal'

const FavoritesList = () => {
  const { items, removeFavorite } = useFavorites()
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <FaHeart className="text-4xl text-gray-300 mx-auto mb-4" />
        <p className="text-gray-400 text-lg">У вас пока нет избранных товаров</p>
        <Link to="/catalog" className="text-[#8A9A86] hover:underline mt-2 inline-block font-medium">
          Перейти в каталог
        </Link>
      </div>
    )
  }

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        {items.map((item) => (
          <div
            key={item.id}
            onClick={() => openModal(item.product_id)}
            className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-5 cursor-pointer flex flex-col border border-gray-100"
          >
            <div className="aspect-square bg-gray-50 rounded-xl flex items-center justify-center text-4xl overflow-hidden">
              <FaLeaf className="text-gray-300 text-5xl" />
            </div>
            <h3 className="font-semibold text-[#1C1C1C] mt-3 truncate text-lg">{item.name}</h3>
            <p className="text-[#8A9A86] font-bold text-xl">{item.price} BYN</p>
            <button
              onClick={(e) => {
                e.stopPropagation()
                removeFavorite(item.product_id)
              }}
              className="mt-3 text-gray-400 hover:text-red-500 text-sm self-start transition flex items-center gap-1.5"
            >
              <FaTrash /> Удалить
            </button>
          </div>
        ))}
      </div>

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </>
  )
}

export default FavoritesList