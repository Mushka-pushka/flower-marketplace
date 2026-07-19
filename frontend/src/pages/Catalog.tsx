import { useEffect, useState } from 'react'
import { FaLeaf, FaChevronLeft, FaChevronRight } from 'react-icons/fa'
import { searchProducts } from '../api/catalog.api'
import type { Product } from '../api/catalog.api'
import ProductModal from '../components/ProductModal'
import SearchBar from '../components/catalog/SearchBar'
import Filters from '../components/catalog/Filters'
import { useSearchParams } from 'react-router-dom'

const Catalog = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  const [limit] = useState(24)
  const [offset, setOffset] = useState(() => {
    const offsetParam = searchParams.get('offset')
    return offsetParam ? Number(offsetParam) : 0
  })

  const [filters, setFilters] = useState({
    q: searchParams.get('q') || '',
    category: searchParams.get('category') || '',
    minPrice: searchParams.get('minPrice') ? Number(searchParams.get('minPrice')) : undefined,
    maxPrice: searchParams.get('maxPrice') ? Number(searchParams.get('maxPrice')) : undefined,
    sortBy: searchParams.get('sortBy') || 'relevance',
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
        limit,
        offset,
      }
      console.log('Sending to backend:', params)
      const response = await searchProducts(params)
      setProducts(response?.items || [])
      setTotal(response?.total || 0)
      
      // Обновляем URL
      const newParams = new URLSearchParams()
      if (filters.q) newParams.set('q', filters.q)
      if (filters.category) newParams.set('category', filters.category)
      if (filters.minPrice) newParams.set('minPrice', String(filters.minPrice))
      if (filters.maxPrice) newParams.set('maxPrice', String(filters.maxPrice))
      if (filters.sortBy && filters.sortBy !== 'relevance') newParams.set('sortBy', filters.sortBy)
      if (offset > 0) newParams.set('offset', String(offset))
      setSearchParams(newParams)
    } catch (err: any) {
      console.error('Ошибка загрузки товаров:', err)
      setError(err.response?.data?.error || 'Ошибка загрузки товаров')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchProducts()
  }, [filters, offset])

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  const handleSearch = (query: string) => {
    console.log('Catalog received search query:', query)
    setFilters((prev) => ({ ...prev, q: query }))
    setOffset(0) // Сбрасываем страницу при поиске
  }

  const handleFilter = (newFilters: {
    category?: string
    minPrice?: number
    maxPrice?: number
    sortBy?: string
  }) => {
    setOffset(0) // Сбрасываем страницу при изменении фильтров
    setFilters((prev) => ({
      ...prev,
      category: newFilters.category !== undefined ? newFilters.category : '',
      minPrice: newFilters.minPrice !== undefined ? newFilters.minPrice : undefined,
      maxPrice: newFilters.maxPrice !== undefined ? newFilters.maxPrice : undefined,
      sortBy: newFilters.sortBy !== undefined ? newFilters.sortBy : 'relevance',
    }))
  }

  const totalPages = Math.ceil(total / limit)
  const currentPage = Math.floor(offset / limit) + 1

  const goToPage = (page: number) => {
    setOffset((page - 1) * limit)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  if (loading) {
    return <div className="text-center py-12 text-gray-400">Загрузка...</div>
  }

  if (error) {
    return <div className="text-center py-12 text-red-500">{error}</div>
  }

  return (
    <div className="animate-fade-in-up">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-3xl font-bold text-[#1C1C1C]">Каталог цветов</h1>
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
          <p className="text-gray-400 text-lg">Товаров не найдено</p>
          <p className="text-gray-400 text-sm mt-2">Попробуйте изменить параметры фильтра</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5 mt-6">
            {products.map((product) => (
              <div
                key={product.id}
                onClick={() => openModal(product.id)}
                className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-5 cursor-pointer h-full flex flex-col border border-gray-100"
              >
                <div className="aspect-square bg-gray-50 rounded-xl mb-3 flex items-center justify-center text-4xl overflow-hidden">
                  <FaLeaf className="text-gray-300 text-5xl" />
                </div>
                <h3 className="font-semibold text-[#1C1C1C] truncate text-base">{product.name}</h3>
                <div className="flex items-baseline gap-0.5 mt-1">
                  <span className="text-[#8A9A86] font-bold text-xl">{product.price}</span>
                  <span className="text-gray-400 text-sm font-medium">BYN</span>
                </div>
                {product.old_price && (
                  <p className="text-gray-400 text-sm line-through">{product.old_price} BYN</p>
                )}
                {product.tags && product.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-2">
                    {product.tags.slice(0, 3).map((tag) => (
                      <span key={tag} className="bg-gray-50 text-[#1C1C1C] text-xs px-2 py-0.5 rounded-full border border-gray-200">
                        {tag}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Пагинация */}
          {totalPages > 1 && (
            <div className="flex justify-center items-center gap-2 mt-8">
              <button
                onClick={() => goToPage(currentPage - 1)}
                disabled={currentPage === 1}
                className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition flex items-center gap-1"
              >
                <FaChevronLeft className="text-sm" />
                Назад
              </button>
              <span className="px-4 py-2 text-[#1C1C1C] font-medium">
                {currentPage} / {totalPages}
              </span>
              <button
                onClick={() => goToPage(currentPage + 1)}
                disabled={currentPage === totalPages}
                className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition flex items-center gap-1"
              >
                Вперед
                <FaChevronRight className="text-sm" />
              </button>
            </div>
          )}
        </>
      )}

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </div>
  )
}

export default Catalog