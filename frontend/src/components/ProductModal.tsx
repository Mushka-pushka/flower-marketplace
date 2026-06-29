import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { useNavigate } from 'react-router-dom'
import Reviews from './Reviews'
import {
  FaTimes,
  FaShoppingCart,
  FaHeart,
  FaRegHeart,
  FaLeaf,
  FaCheckCircle,
  FaExclamationCircle,
  FaTrash,
} from 'react-icons/fa'
import { getProductById } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'
import { useCart } from '../context/CartContext'
import { useFavorites } from '../context/FavoritesContext'
import { useAuth } from '../context/AuthContext'

interface ProductModalProps {
  productId: string | null
  onClose: () => void
}

const ProductModal = ({ productId, onClose }: ProductModalProps) => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { items, addToCart, removeFromCart } = useCart()
  const { toggleFavorite, isFavorite } = useFavorites()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showNotification, setShowNotification] = useState(false)
  const [authNotification, setAuthNotification] = useState(false)
  const [isInCart, setIsInCart] = useState(false)

  useEffect(() => {
    if (!productId) return

    const fetchProduct = async () => {
      try {
        setLoading(true)
        const data = await getProductById(productId)
        setProduct(data)
        // Проверяем, есть ли товар в корзине
        setIsInCart(items.some(item => item.product_id === productId))
      } catch (err) {
        setError('Не удалось загрузить товар')
      } finally {
        setLoading(false)
      }
    }

    fetchProduct()
  }, [productId, items])

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

  useEffect(() => {
    if (authNotification) {
      const timer = setTimeout(() => {
        setAuthNotification(false)
      }, 3000)
      return () => clearTimeout(timer)
    }
  }, [authNotification])

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
    if (!user) {
      setAuthNotification(true)
      return
    }
    if (product) {
      addToCart({
        id: product.id,
        name: product.name,
        price: product.price,
        shop_id: product.shop_id,
      })
      setIsInCart(true)
      setShowNotification(true)
    }
  }

  const handleRemoveFromCart = () => {
    if (product) {
      removeFromCart(product.id)
      setIsInCart(false)
    }
  }

  const handleToggleFavorite = () => {
    if (!user) {
      setAuthNotification(true)
      return
    }
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
        className="bg-white rounded-2xl w-full max-w-4xl max-h-[90vh] overflow-y-auto relative shadow-[0_8px_40px_rgba(0,0,0,0.08)] animate-fade-in-up border border-gray-100"
        onClick={(e) => e.stopPropagation()}
      >
        {showNotification && (
          <div className="sticky top-4 z-20 flex justify-center pointer-events-none">
            <div className="bg-[#8A9A86] text-white px-6 py-3 rounded-xl shadow-lg animate-bounce font-medium pointer-events-auto flex items-center gap-2">
              <FaCheckCircle /> Товар добавлен в корзину!
            </div>
          </div>
        )}

        {authNotification && (
          <div className="sticky top-4 z-20 flex justify-center pointer-events-none">
            <div
              className="bg-amber-50 text-amber-700 px-6 py-3 rounded-xl shadow-lg animate-bounce font-medium pointer-events-auto flex items-center gap-2 border border-amber-200 cursor-pointer"
              onClick={() => navigate('/login')}
            >
              <FaExclamationCircle /> Войдите в аккаунт, чтобы добавить товар
            </div>
          </div>
        )}

        <div className="flex justify-between items-center p-5 border-b border-gray-100 sticky top-0 bg-white z-10 rounded-t-2xl">
          <h2 className="text-2xl font-bold text-[#1C1C1C]">Детали товара</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-3xl leading-none transition"
          >
            <FaTimes />
          </button>
        </div>

        <div className="p-6">
          {loading && (
            <div className="text-center py-12 text-gray-400">Загрузка...</div>
          )}

          {error && (
            <div className="text-center py-12 text-red-500">{error}</div>
          )}

          {product && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div className="aspect-square bg-gray-50 rounded-xl flex items-center justify-center text-7xl overflow-hidden">
                <FaLeaf className="text-gray-300 text-6xl" />
              </div>

              <div>
                <h3 className="text-3xl font-bold text-[#1C1C1C]">{product.name}</h3>

                <div className="flex items-baseline gap-1 mt-2">
                  <span className="text-4xl font-bold text-[#8A9A86]">{product.price}</span>
                  <span className="text-gray-400 text-sm font-medium">BYN</span>
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
                      <span key={tag} className="bg-gray-50 text-[#1C1C1C] text-sm px-3 py-1 rounded-full border border-gray-200">
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <div className="mt-6 space-y-3">
                  {isInCart ? (
                    <button
                      onClick={handleRemoveFromCart}
                      className="w-full bg-red-50 text-red-600 border-2 border-red-200 py-3.5 rounded-xl hover:bg-red-100 transition flex items-center justify-center gap-2 text-base font-medium"
                    >
                      <FaTrash /> Удалить из корзины
                    </button>
                  ) : (
                    <button
                      onClick={handleAddToCart}
                      className="w-full bg-[#8A9A86] text-white py-3.5 rounded-xl hover:bg-[#7A8A76] transition flex items-center justify-center gap-2 text-base font-medium"
                    >
                      <FaShoppingCart /> Добавить в корзину
                    </button>
                  )}

                  <button
                    onClick={handleToggleFavorite}
                    className={`w-full py-3.5 rounded-xl text-base font-medium transition border-2 flex items-center justify-center gap-2 ${
                      isFavorite(product.id)
                        ? 'bg-gray-50 border-[#8A9A86] text-[#8A9A86]'
                        : 'border-gray-200 text-[#1C1C1C] hover:border-[#8A9A86] hover:bg-gray-50'
                    }`}
                  >
                    {isFavorite(product.id) ? (
                      <>
                        <FaHeart className="text-[#8A9A86]" /> В избранном
                      </>
                    ) : (
                      <>
                        <FaRegHeart /> В избранное
                      </>
                    )}
                  </button>
                </div>

                <div className="mt-6 border-t border-gray-100 pt-4">
                  <Reviews productId={product.id} />
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