import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { getProductById } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'
import { useCart } from '../context/CartContext'
import { useFavorites } from '../context/FavoritesContext'

interface ProductModalProps {
  productId: string | null
  onClose: () => void
}

const ProductModal = ({ productId, onClose }: ProductModalProps) => {
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showNotification, setShowNotification] = useState(false)

  const { addToCart } = useCart()
  const { toggleFavorite, isFavorite } = useFavorites()

  useEffect(() => {
    if (!productId) return

    const fetchProduct = async () => {
      try {
        setLoading(true)
        const data = await getProductById(productId)
        setProduct(data)
      } catch (err) {
        setError('Не удалось загрузить товар')
      } finally {
        setLoading(false)
      }
    }

    fetchProduct()
  }, [productId])

  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleEsc)
    return () => window.removeEventListener('keydown', handleEsc)
  }, [onClose])

  useEffect(() => {
    if (showNotification) {
      const timer = setTimeout(() => {
        setShowNotification(false)
      }, 2000)
      return () => clearTimeout(timer)
    }
  }, [showNotification])

  // Блокируем прокрутку страницы
  useEffect(() => {
    if (productId) {
      document.body.style.overflow = 'hidden'
    }
    return () => {
      document.body.style.overflow = 'auto'
    }
  }, [productId])

  if (!productId) return null

  const handleAddToCart = () => {
    if (product) {
      addToCart({
        id: product.id,
        name: product.name,
        price: product.price,
      })
      setShowNotification(true)
    }
  }

  const handleToggleFavorite = () => {
    if (product) {
      toggleFavorite({
        id: product.id,
        name: product.name,
        price: product.price,
      })
    }
  }

  const modalContent = (
    <div 
      className="fixed inset-0 z-[9999] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
      onClick={onClose}
    >
      <div 
        className="bg-white/90 backdrop-blur-md rounded-3xl w-full max-w-4xl max-h-[90vh] overflow-y-auto relative shadow-2xl border border-pink-50/50 animate-fade-in-up"
        onClick={(e) => e.stopPropagation()}
      >
        {showNotification && (
          <div className="sticky top-4 z-20 flex justify-center pointer-events-none">
            <div className="bg-green-500 text-white px-6 py-3 rounded-full shadow-lg animate-bounce font-medium pointer-events-auto">
              ✅ Товар добавлен в корзину!
            </div>
          </div>
        )}

        <div className="flex justify-between items-center p-4 border-b border-pink-100 sticky top-0 bg-white/90 backdrop-blur-sm z-10 rounded-t-3xl">
          <h2 className="text-2xl font-bold gradient-text">Детали товара</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-4xl leading-none transition"
          >
            ×
          </button>
        </div>

        <div className="p-6">
          {loading && (
            <div className="text-center py-12 text-gray-500">Загрузка...</div>
          )}

          {error && (
            <div className="text-center py-12 text-red-500">{error}</div>
          )}

          {product && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="aspect-square bg-gradient-to-br from-pink-50 to-purple-50 rounded-2xl flex items-center justify-center text-7xl text-gray-300 shadow-inner">
                🌸
              </div>

              <div>
                <h3 className="text-3xl font-bold text-gray-800">{product.name}</h3>

                <div className="flex items-baseline gap-1 mt-2">
                  <span className="text-4xl font-bold text-pink-600">{product.price}</span>
                  <span className="text-gray-500 text-sm font-medium">BYN</span>
                </div>
                {product.old_price && (
                  <p className="text-gray-400 text-sm line-through">{product.old_price} BYN</p>
                )}

                <p className="text-gray-600 mt-4 leading-relaxed">
                  {product.description || 'Описание отсутствует'}
                </p>

                {product.tags && product.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2 mt-4">
                    {product.tags.map((tag) => (
                      <span key={tag} className="bg-pink-50 text-pink-600 text-sm px-3 py-1 rounded-full border border-pink-100">
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <div className="mt-6 space-y-3">
                  <button
                    onClick={handleAddToCart}
                    className="btn-primary w-full py-3 rounded-full text-lg font-medium"
                  >
                    🛒 Добавить в корзину
                  </button>

                  <button
                    onClick={handleToggleFavorite}
                    className={`w-full py-3 rounded-full text-lg font-medium transition border-2 ${
                      isFavorite(product.id)
                        ? 'bg-pink-50 border-pink-300 text-pink-600'
                        : 'border-gray-300 text-gray-700 hover:border-pink-300 hover:bg-pink-50'
                    }`}
                  >
                    {isFavorite(product.id) ? '❤️ В избранном' : '🤍 В избранное'}
                  </button>
                </div>

                <div className="mt-6 border-t border-pink-100 pt-4">
                  <h4 className="font-semibold text-gray-700">⭐ Отзывы</h4>
                  <p className="text-gray-400 text-sm">Отзывы появятся позже</p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )

  return createPortal(modalContent, document.body)
}

export default ProductModal