import { useState, useEffect, useRef } from 'react'
import { FaSearch } from 'react-icons/fa'
import { getAutocomplete } from '../../api/catalog.api'

interface Suggestion {
  text: string
  type: string
  slug: string
  score: number
}

interface SearchBarProps {
  onSearch: (query: string) => void
  initialValue?: string
}

const SearchBar = ({ onSearch, initialValue = '' }: SearchBarProps) => {
  const [query, setQuery] = useState(initialValue)
  const [suggestions, setSuggestions] = useState<Suggestion[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [loading, setLoading] = useState(false)
  const wrapperRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setShowSuggestions(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    const fetchSuggestions = async () => {
      if (query.length < 2) {
        setSuggestions([])
        return
      }
      setLoading(true)
      try {
        const data = await getAutocomplete(query, 6)
        setSuggestions(data)
        setShowSuggestions(true)
      } catch (error) {
        console.error('Ошибка автодополнения:', error)
      } finally {
        setLoading(false)
      }
    }

    const timer = setTimeout(fetchSuggestions, 300)
    return () => clearTimeout(timer)
  }, [query])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setShowSuggestions(false)
    onSearch(query)
  }

  const handleSuggestionClick = (suggestion: Suggestion) => {
    setQuery(suggestion.text)
    setShowSuggestions(false)
    onSearch(suggestion.text)
  }

  return (
    <div ref={wrapperRef} className="relative w-full max-w-md">
      <form onSubmit={handleSubmit} className="flex items-center">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Поиск цветов..."
          className="w-full px-4 py-2.5 border border-gray-200 rounded-l-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] text-sm"
        />
        <button
          type="submit"
          className="bg-[#8A9A86] text-white px-5 py-2.5 rounded-r-xl hover:bg-[#7A8A76] transition flex items-center gap-2 text-sm font-medium"
        >
          <FaSearch />
        </button>
      </form>

      {showSuggestions && suggestions.length > 0 && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] z-20 max-h-60 overflow-y-auto">
          {suggestions.map((suggestion, index) => (
            <div
              key={index}
              onClick={() => handleSuggestionClick(suggestion)}
              className="px-4 py-2.5 hover:bg-gray-50 cursor-pointer flex items-center gap-2 transition"
            >
              <span className="text-gray-400 text-sm">
                {suggestion.type === 'product' && '🌸'}
                {suggestion.type === 'category' && '📁'}
                {suggestion.type === 'tag' && '🏷️'}
              </span>
              <span className="text-[#1C1C1C]">{suggestion.text}</span>
              <span className="text-gray-400 text-xs ml-auto">{suggestion.type}</span>
            </div>
          ))}
        </div>
      )}

      {loading && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] z-20 p-3 text-center text-gray-400 text-sm">
          Загрузка...
        </div>
      )}
    </div>
  )
}

export default SearchBar