import { useState, useEffect, useRef } from 'react'
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
          className="w-full px-4 py-2 border border-gray-300 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-pink-400"
        />
        <button
          type="submit"
          className="bg-pink-500 text-white px-4 py-2 rounded-r-lg hover:bg-pink-600 transition"
        >
          🔍
        </button>
      </form>

      {/* Автодополнение */}
      {showSuggestions && suggestions.length > 0 && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-20 max-h-60 overflow-y-auto">
          {suggestions.map((suggestion, index) => (
            <div
              key={index}
              onClick={() => handleSuggestionClick(suggestion)}
              className="px-4 py-2 hover:bg-gray-50 cursor-pointer flex items-center gap-2"
            >
              <span className="text-gray-400 text-sm">
                {suggestion.type === 'product' && '🌸'}
                {suggestion.type === 'category' && '📁'}
                {suggestion.type === 'tag' && '🏷️'}
              </span>
              <span className="text-gray-800">{suggestion.text}</span>
              <span className="text-gray-400 text-xs ml-auto">{suggestion.type}</span>
            </div>
          ))}
        </div>
      )}

      {loading && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-20 p-3 text-center text-gray-400 text-sm">
          Загрузка...
        </div>
      )}
    </div>
  )
}

export default SearchBar