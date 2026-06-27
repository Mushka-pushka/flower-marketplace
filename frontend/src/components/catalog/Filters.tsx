import { useEffect, useState } from 'react'
import { getCategories } from '../../api/catalog.api'
import type { Category } from '../../api/catalog.api'

interface FiltersProps {
  onFilter: (filters: {
    category?: string
    minPrice?: number
    maxPrice?: number
    sortBy?: string
  }) => void
  initialFilters?: {
    category?: string
    minPrice?: number
    maxPrice?: number
    sortBy?: string
  }
}

const Filters = ({ onFilter, initialFilters = {} }: FiltersProps) => {
  const [categories, setCategories] = useState<Category[]>([])
  const [selectedCategory, setSelectedCategory] = useState(initialFilters.category || '')
  const [minPrice, setMinPrice] = useState(initialFilters.minPrice?.toString() || '')
  const [maxPrice, setMaxPrice] = useState(initialFilters.maxPrice?.toString() || '')
  const [sortBy, setSortBy] = useState(initialFilters.sortBy || 'relevance')

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const data = await getCategories()
        setCategories(data)
      } catch (error) {
        console.error('Ошибка загрузки категорий:', error)
      }
    }
    fetchCategories()
  }, [])

  const applyFilters = () => {
    onFilter({
      category: selectedCategory || undefined,
      minPrice: minPrice ? Number(minPrice) : undefined,
      maxPrice: maxPrice ? Number(maxPrice) : undefined,
      sortBy: sortBy || undefined,
    })
  }

  const resetFilters = () => {
    setSelectedCategory('')
    setMinPrice('')
    setMaxPrice('')
    setSortBy('relevance')
    onFilter({
      category: '',
      minPrice: undefined,
      maxPrice: undefined,
      sortBy: 'relevance',
    })
  }

  const handlePriceChange = (value: string, setter: (val: string) => void) => {
    const cleaned = value.replace(/[^0-9.]/g, '')
    setter(cleaned)
  }

  return (
    <div className="bg-white rounded-lg shadow p-4 space-y-4">
      <div className="flex flex-wrap gap-4 items-end">
        <div className="flex-1 min-w-[150px]">
          <label className="block text-sm font-medium text-gray-700 mb-1">Категория</label>
          <select
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-400"
          >
            <option value="">Все категории</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.slug}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>

        <div className="flex-1 min-w-[100px]">
          <label className="block text-sm font-medium text-gray-700 mb-1">Цена от</label>
          <input
            type="text"
            value={minPrice}
            onChange={(e) => handlePriceChange(e.target.value, setMinPrice)}
            placeholder="0"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-400"
          />
        </div>

        <div className="flex-1 min-w-[100px]">
          <label className="block text-sm font-medium text-gray-700 mb-1">Цена до</label>
          <input
            type="text"
            value={maxPrice}
            onChange={(e) => handlePriceChange(e.target.value, setMaxPrice)}
            placeholder="1000"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-400"
          />
        </div>

        <div className="flex-1 min-w-[150px]">
          <label className="block text-sm font-medium text-gray-700 mb-1">Сортировка</label>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-pink-400"
          >
            <option value="relevance">По релевантности</option>
            <option value="price_asc">Сначала дешёвые</option>
            <option value="price_desc">Сначала дорогие</option>
            <option value="rating">По рейтингу</option>
            <option value="newest">Сначала новые</option>
          </select>
        </div>

        <div className="flex gap-2">
          <button
            onClick={applyFilters}
            className="bg-pink-500 text-white px-4 py-2 rounded-lg hover:bg-pink-600 transition"
          >
            Применить
          </button>
          <button
            onClick={resetFilters}
            className="bg-gray-200 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-300 transition"
          >
            Сбросить
          </button>
        </div>
      </div>
    </div>
  )
}

export default Filters