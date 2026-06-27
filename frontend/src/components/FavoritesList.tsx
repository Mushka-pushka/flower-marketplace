import { useState } from 'react'
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
    return <p className="text-gray-500">У вас пока нет избранных товаров</p>
  }

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {items.map((item) => (
          <div
            key={item.id}
            onClick={() => openModal(item.product_id)}
            className="border rounded-lg p-4 bg-white shadow-sm hover:shadow-md transition cursor-pointer flex flex-col"
          >
            <div className="aspect-square bg-gray-100 rounded-lg flex items-center justify-center text-4xl text-gray-300">
              🌸
            </div>
            <h3 className="font-semibold text-gray-800 mt-2 truncate">{item.name}</h3>
            <p className="text-pink-600 font-bold text-lg">{item.price} BYN</p>
            <button
              onClick={(e) => {
                e.stopPropagation()
                removeFavorite(item.product_id)
              }}
              className="mt-2 text-red-400 hover:text-red-600 text-sm self-start"
            >
              🗑️ Удалить
            </button>
          </div>
        ))}
      </div>

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </>
  )
}

export default FavoritesList