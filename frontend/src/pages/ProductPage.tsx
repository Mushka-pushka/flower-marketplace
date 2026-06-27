import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getProductById } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'

const ProductPage = () => {
  const { id } = useParams<{ id: string }>()
  const [product, setProduct] = useState<Product | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const fetchProduct = async () => {
      if (!id) return
      try {
        setLoading(true)
        const data = await getProductById(id)
        setProduct(data)
      } catch (err: any) {
        setError('Товар не найден')
      } finally {
        setLoading(false)
      }
    }
    fetchProduct()
  }, [id])

  if (loading) return <div className="text-center py-12">Загрузка...</div>
  if (error || !product) return <div className="text-center py-12 text-red-500">{error || 'Товар не найден'}</div>

  return (
    <div className="max-w-4xl mx-auto">
      <Link to="/catalog" className="text-pink-500 hover:underline mb-4 inline-block">
        ← Назад к каталогу
      </Link>

      <div className="bg-white rounded-lg shadow p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Фото */}
          <div className="aspect-square bg-gray-100 rounded-lg flex items-center justify-center text-6xl text-gray-300">
            🌸
          </div>

          {/* Информация о товаре */}
          <div>
            <h1 className="text-3xl font-bold text-gray-800">{product.name}</h1>

            {/* Цена */}
            <div className="flex items-baseline gap-1 mt-2">
              <span className="text-3xl font-bold text-pink-600">{product.price}</span>
              <span className="text-gray-500 text-sm">BYN</span>
            </div>
            {product.old_price && (
              <p className="text-gray-400 text-sm line-through">{product.old_price} BYN</p>
            )}

            {/* Описание */}
            <p className="text-gray-600 mt-4">{product.description || 'Описание отсутствует'}</p>

            {/* Теги */}
            {product.tags && product.tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-4">
                {product.tags.map((tag) => (
                  <span key={tag} className="bg-gray-100 text-gray-600 text-sm px-3 py-1 rounded-full">
                    #{tag}
                  </span>
                ))}
              </div>
            )}

            {/* Кнопка "В корзину" */}
            <button className="mt-6 w-full bg-pink-500 text-white py-3 rounded-lg hover:bg-pink-600 transition text-lg font-medium">
              🛒 Добавить в корзину
            </button>

            {/* Отзывы (заглушка) */}
            <div className="mt-8 border-t pt-4">
              <h3 className="font-semibold text-gray-700">⭐ Отзывы</h3>
              <p className="text-gray-400 text-sm">Отзывы появятся позже</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ProductPage