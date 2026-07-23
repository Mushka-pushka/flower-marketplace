import { useState, useEffect, useRef, type KeyboardEvent } from 'react'
import { FaSearch, FaLeaf, FaFolder, FaTag } from 'react-icons/fa'
import { getAutocomplete } from '../../api/catalog.api'
import type { AutocompleteSuggestion } from '../../api/catalog.api'

interface SearchBarProps {
  onSearch: (query: string) => void
  initialValue?: string
}

const SearchBar = ({ onSearch, initialValue = '' }: SearchBarProps) => {
  const [query, setQuery] = useState(initialValue)
  const [suggestions, setSuggestions] = useState<AutocompleteSuggestion[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [loading, setLoading] = useState(false)
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const wrapperRef = useRef<HTMLDivElement>(null)
  const listRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setShowSuggestions(false)
        setSelectedIndex(-1)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    const fetchSuggestions = async () => {
      if (query.length < 2) {
        setSuggestions([])
        setSelectedIndex(-1)
        return
      }
      setLoading(true)
      try {
        const data = await getAutocomplete(query, 6)
        setSuggestions(Array.isArray(data) ? data : [])
        setShowSuggestions(true)
        setSelectedIndex(-1)
      } catch (error) {
        console.error('Ошибка автодополнения:', error)
        setSuggestions([])
      } finally {
        setLoading(false)
      }
    }

    const timer = setTimeout(fetchSuggestions, 300)
    return () => clearTimeout(timer)
  }, [query])

  useEffect(() => {
    if (selectedIndex >= 0 && listRef.current) {
      const items = listRef.current.children
      if (items[selectedIndex]) {
        items[selectedIndex].scrollIntoView({ block: 'nearest' })
      }
    }
  }, [selectedIndex])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setShowSuggestions(false)
    setSelectedIndex(-1)
    if (query.trim()) {
      console.log('SearchBar submitting:', query.trim())
      onSearch(query.trim())
    }
  }

  const handleSuggestionClick = (suggestion: AutocompleteSuggestion) => {
    setQuery(suggestion.text)
    setShowSuggestions(false)
    setSelectedIndex(-1)
    console.log('SearchBar suggestion clicked:', suggestion.text)
    onSearch(suggestion.text)
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (!showSuggestions || !Array.isArray(suggestions) || suggestions.length === 0) {
      if (e.key === 'Enter' && query.trim()) {
        handleSubmit(e)
      }
      return
    }

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setSelectedIndex(prev => 
          prev < suggestions.length - 1 ? prev + 1 : prev
        )
        break
      case 'ArrowUp':
        e.preventDefault()
        setSelectedIndex(prev => prev > 0 ? prev - 1 : -1)
        break
      case 'Enter':
        e.preventDefault()
        if (selectedIndex >= 0 && selectedIndex < suggestions.length) {
          handleSuggestionClick(suggestions[selectedIndex])
        } else if (query.trim()) {
          handleSubmit(e)
        }
        break
      case 'Escape':
        setShowSuggestions(false)
        setSelectedIndex(-1)
        break
    }
  }

  return (
    <div ref={wrapperRef} className="relative w-full max-w-md">
      <form onSubmit={handleSubmit} className="flex items-center">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => query.length >= 2 && setShowSuggestions(true)}
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

      {showSuggestions && Array.isArray(suggestions) && suggestions.length > 0 && (
        <div 
          ref={listRef}
          className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] z-20 max-h-60 overflow-y-auto"
        >
          {suggestions.map((suggestion, index) => (
            <div
              key={index}
              onClick={() => handleSuggestionClick(suggestion)}
              onMouseEnter={() => setSelectedIndex(index)}
              className={`px-4 py-2.5 cursor-pointer flex items-center gap-2 transition ${
                index === selectedIndex ? 'bg-[#8A9A86]/10' : 'hover:bg-gray-50'
              }`}
            >
              <span className="text-gray-400 text-sm flex items-center justify-center w-5">
                {suggestion.type === 'product' && <FaLeaf className="text-[#8A9A86]" />}
                {suggestion.type === 'category' && <FaFolder className="text-[#8A9A86]" />}
                {suggestion.type === 'tag' && <FaTag className="text-[#8A9A86]" />}
              </span>
              <span className="text-[#1C1C1C]">{suggestion.text}</span>
              <span className="text-gray-400 text-xs ml-auto capitalize">{suggestion.type}</span>
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