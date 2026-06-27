import { useEffect, useState } from 'react'
import { searchProducts } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'
import ProductModal from '../components/ProductModal'

const Catalog = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true)
        const response = await searchProducts({ limit: 24 })
        setProducts(response.items)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Ошибка загрузки товаров')
      } finally {
        setLoading(false)
      }
    }

    fetchProducts()
  }, [])

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  if (loading) {
    return <div className="text-center py-12 text-gray-500">Загрузка...</div>
  }

  if (error) {
    return <div className="text-center py-12 text-red-500">{error}</div>
  }

  if (products.length === 0) {
    return <div className="text-center py-12 text-gray-500">Товаров пока нет</div>
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">🌺 Каталог цветов</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {products.map((product) => (
          <div
            key={product.id}
            onClick={() => openModal(product.id)}
            className="bg-white rounded-lg shadow hover:shadow-lg transition p-4 cursor-pointer h-full flex flex-col"
          >
            <div className="aspect-square bg-gray-100 rounded-lg mb-3 flex items-center justify-center text-gray-400 text-4xl">
              🌸
            </div>
            <h3 className="font-semibold text-gray-800 truncate">{product.name}</h3>
            
            <div className="flex items-baseline gap-0.5 mt-1">
              <span className="text-pink-600 font-bold text-xl">{product.price}</span>
              <span className="text-gray-500 text-sm font-medium">BYN</span>
            </div>

            {product.old_price && (
              <p className="text-gray-400 text-sm line-through">{product.old_price} BYN</p>
            )}

            {product.tags && product.tags.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {product.tags.slice(0, 3).map((tag) => (
                  <span key={tag} className="bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">
                    {tag}
                  </span>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Модальное окно с товаром */}
      <ProductModal
        productId={selectedProductId}
        onClose={closeModal}
      />
    </div>
  )
}

export default Catalog