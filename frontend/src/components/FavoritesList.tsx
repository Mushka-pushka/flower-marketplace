import { useState } from 'react'
import { Link } from 'react-router-dom'
import { FaHeart, FaTrash, FaLeaf, FaSpinner } from 'react-icons/fa'
import { useFavorites } from '../context/FavoritesContext'
import ProductModal from './ProductModal'

const FavoritesList = () => {
  const { items, removeFavorite, loading: favoritesLoading } = useFavorites()
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)
  const [visibleItems, setVisibleItems] = useState(6)
  const [isLoading, setIsLoading] = useState(false)

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  const loadMore = () => {
    setIsLoading(true)
    setTimeout(() => {
      setVisibleItems(prev => Math.min(prev + 6, items.length))
      setIsLoading(false)
    }, 500)
  }

  if (favoritesLoading) {
    return (
      <div className="text-center py-12">
        <FaSpinner className="text-3xl text-[#8A9A86] animate-spin mx-auto" />
        <p className="text-gray-400 mt-3">Загрузка избранного...</p>
      </div>
    )
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

  const displayedItems = items.slice(0, visibleItems)

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        {displayedItems.map((item) => {
          console.log('Rendering favorite:', item.name, 'Image:', item.image)
          return (
            <div
              key={item.id}
              onClick={() => openModal(item.product_id)}
              className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-5 cursor-pointer flex flex-col border border-gray-100"
            >
              <div className="aspect-square bg-gray-50 rounded-xl flex items-center justify-center overflow-hidden">
                {item.image ? (
                  <img 
                    src={`http://localhost:8082${item.image}`} 
                    alt={item.name} 
                    className="w-full h-full object-cover"
                    loading="lazy"
                  />
                ) : (
                  <FaLeaf className="text-gray-300 text-5xl" />
                )}
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
          )
        })}
      </div>

      {visibleItems < items.length && (
        <div className="text-center mt-6">
          <button
            onClick={loadMore}
            disabled={isLoading}
            className="px-6 py-2 border border-[#8A9A86] text-[#8A9A86] rounded-xl hover:bg-[#8A9A86] hover:text-white transition text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? (
              <>
                <FaSpinner className="inline animate-spin mr-2" />
                Загрузка...
              </>
            ) : (
              `Показать еще (${items.length - visibleItems} осталось)`
            )}
          </button>
        </div>
      )}

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </>
  )
}

export default FavoritesList