import { useEffect, useState } from 'react'
import { FaFilter, FaUndo } from 'react-icons/fa'
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
    <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 space-y-4 border border-gray-100">
      <div className="flex flex-wrap gap-4 items-end">
        <div className="flex-1 min-w-[150px]">
          <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">Категория</label>
          <select
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
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
          <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">Цена от</label>
          <input
            type="text"
            value={minPrice}
            onChange={(e) => handlePriceChange(e.target.value, setMinPrice)}
            placeholder="0"
            className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] text-sm"
          />
        </div>

        <div className="flex-1 min-w-[100px]">
          <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">Цена до</label>
          <input
            type="text"
            value={maxPrice}
            onChange={(e) => handlePriceChange(e.target.value, setMaxPrice)}
            placeholder="1000"
            className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] text-sm"
          />
        </div>

        <div className="flex-1 min-w-[150px]">
          <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">Сортировка</label>
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
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
            className="bg-[#8A9A86] text-white px-5 py-2.5 rounded-xl hover:bg-[#7A8A76] transition flex items-center gap-2 text-sm font-medium"
          >
            <FaFilter className="text-sm" />
            Применить
          </button>
          <button
            onClick={resetFilters}
            className="bg-gray-100 text-[#1C1C1C] px-5 py-2.5 rounded-xl hover:bg-gray-200 transition flex items-center gap-2 text-sm font-medium"
          >
            <FaUndo className="text-sm" />
            Сбросить
          </button>
        </div>
      </div>
    </div>
  )
}

export default Filters