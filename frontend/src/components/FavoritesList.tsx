import { useState } from 'react'
import { Link } from 'react-router-dom' 
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
      <div className="text-center py-8">
        <p className="text-gray-400 text-lg">💔 У вас пока нет избранных товаров</p>
        <Link to="/catalog" className="text-pink-500 hover:underline mt-2 inline-block">
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
            className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 p-5 cursor-pointer flex flex-col card-hover border border-pink-50/50"
          >
            <div className="aspect-square bg-gradient-to-br from-pink-50 to-purple-50 rounded-xl flex items-center justify-center text-4xl">
              🌸
            </div>
            <h3 className="font-semibold text-gray-800 mt-3 truncate text-lg">{item.name}</h3>
            <p className="text-pink-600 font-bold text-xl">{item.price} BYN</p>
            <button
              onClick={(e) => {
                e.stopPropagation()
                removeFavorite(item.product_id)
              }}
              className="mt-3 text-gray-400 hover:text-red-500 text-sm self-start transition"
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