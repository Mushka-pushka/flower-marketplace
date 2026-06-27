import { useEffect, useState } from 'react'
import { searchProducts } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'
import ProductModal from '../components/ProductModal'
import SearchBar from '../components/catalog/SearchBar'
import Filters from '../components/catalog/Filters'

const Catalog = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)
  const [filters, setFilters] = useState({
    q: '',
    category: '',
    minPrice: undefined as number | undefined,
    maxPrice: undefined as number | undefined,
    sortBy: 'relevance' as string,
  })

  const fetchProducts = async () => {
    try {
      setLoading(true)
      setError('')
      const params = {
        q: filters.q || undefined,
        category: filters.category || undefined,
        min_price: filters.minPrice,
        max_price: filters.maxPrice,
        sort_by: filters.sortBy || undefined,
        limit: 24,
      }
      const response = await searchProducts(params)
      setProducts(response?.items || []) // <-- ЗАЩИТА ОТ null
    } catch (err: any) {
      console.error('Ошибка загрузки товаров:', err)
      setError(err.response?.data?.error || 'Ошибка загрузки товаров')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
  }, [filters])

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  const handleSearch = (query: string) => {
    setFilters((prev) => ({ ...prev, q: query }))
  }

  const handleFilter = (newFilters: {
    category?: string
    minPrice?: number
    maxPrice?: number
    sortBy?: string
  }) => {
    setFilters((prev) => ({
      ...prev,
      category: newFilters.category !== undefined ? newFilters.category : '',
      minPrice: newFilters.minPrice !== undefined ? newFilters.minPrice : undefined,
      maxPrice: newFilters.maxPrice !== undefined ? newFilters.maxPrice : undefined,
      sortBy: newFilters.sortBy !== undefined ? newFilters.sortBy : 'relevance',
    }))
  }

  if (loading) {
    return <div className="text-center py-12 text-gray-500">Загрузка...</div>
  }

  if (error) {
    return <div className="text-center py-12 text-red-500">{error}</div>
  }

  return (
    <div>
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-3xl font-bold">🌺 Каталог цветов</h1>
        <SearchBar onSearch={handleSearch} initialValue={filters.q} />
      </div>

      <Filters 
        onFilter={handleFilter} 
        initialFilters={{
          category: filters.category || '',
          minPrice: filters.minPrice,
          maxPrice: filters.maxPrice,
          sortBy: filters.sortBy || 'relevance',
        }}
      />

      {products.length === 0 ? (
        <div className="text-center py-12 mt-6">
          <p className="text-gray-500 text-lg">🌸 Товаров не найдено</p>
          <p className="text-gray-400 text-sm mt-2">Попробуйте изменить параметры фильтра</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 mt-6">
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
      )}

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </div>
  )
}

export default Catalog