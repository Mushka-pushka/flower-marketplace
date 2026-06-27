import { useEffect, useState } from 'react'
import { getProductById } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'

interface ProductModalProps {
  productId: string | null
  onClose: () => void
}

const ProductModal = ({ productId, onClose }: ProductModalProps) => {
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

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

  // Закрытие по Escape
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleEsc)
    return () => window.removeEventListener('keydown', handleEsc)
  }, [onClose])

  if (!productId) return null

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
      onClick={onClose}
    >
      <div 
        className="bg-white rounded-2xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Заголовок модалки */}
        <div className="flex justify-between items-center p-4 border-b">
          <h2 className="text-xl font-bold text-gray-800">Детали товара</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 text-3xl leading-none"
          >
            ×
          </button>
        </div>

        {/* Контент */}
        <div className="p-6">
          {loading && (
            <div className="text-center py-12 text-gray-500">Загрузка...</div>
          )}

          {error && (
            <div className="text-center py-12 text-red-500">{error}</div>
          )}

          {product && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {/* Фото */}
              <div className="aspect-square bg-gray-100 rounded-lg flex items-center justify-center text-6xl text-gray-300">
                🌸
              </div>

              {/* Информация */}
              <div>
                <h3 className="text-2xl font-bold text-gray-800">{product.name}</h3>

                <div className="flex items-baseline gap-1 mt-2">
                  <span className="text-3xl font-bold text-pink-600">{product.price}</span>
                  <span className="text-gray-500 text-sm">BYN</span>
                </div>
                {product.old_price && (
                  <p className="text-gray-400 text-sm line-through">{product.old_price} BYN</p>
                )}

                <p className="text-gray-600 mt-4">{product.description || 'Описание отсутствует'}</p>

                {product.tags && product.tags.length > 0 && (
                  <div className="flex flex-wrap gap-2 mt-4">
                    {product.tags.map((tag) => (
                      <span key={tag} className="bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <button className="mt-6 w-full bg-pink-500 text-white py-3 rounded-lg hover:bg-pink-600 transition text-lg font-medium">
                  🛒 Добавить в корзину
                </button>

                <div className="mt-6 border-t pt-4">
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
}

export default ProductModal