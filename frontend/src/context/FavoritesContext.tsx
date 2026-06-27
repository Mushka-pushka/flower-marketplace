import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'

interface FavoriteItem {
  id: string
  product_id: string
  name: string
  price: number
  image?: string
}

interface FavoritesContextType {
  items: FavoriteItem[]
  addFavorite: (product: { id: string; name: string; price: number }) => void
  removeFavorite: (productId: string) => void
  isFavorite: (productId: string) => boolean
  toggleFavorite: (product: { id: string; name: string; price: number }) => void
}

const FAVORITES_STORAGE_KEY = 'flower_favorites'

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined)

export const FavoritesProvider = ({ children }: { children: ReactNode }) => {
  const [items, setItems] = useState<FavoriteItem[]>(() => {
    const saved = localStorage.getItem(FAVORITES_STORAGE_KEY)
    if (saved) {
      try {
        return JSON.parse(saved)
      } catch {
        return []
      }
    }
    return []
  })

  useEffect(() => {
    localStorage.setItem(FAVORITES_STORAGE_KEY, JSON.stringify(items))
  }, [items])

  const addFavorite = (product: { id: string; name: string; price: number }) => {
    setItems((prev) => {
      if (prev.some((item) => item.product_id === product.id)) return prev
      return [
        ...prev,
        {
          id: crypto.randomUUID(),
          product_id: product.id,
          name: product.name,
          price: product.price,
        },
      ]
    })
  }

  const removeFavorite = (productId: string) => {
    setItems((prev) => prev.filter((item) => item.product_id !== productId))
  }

  const isFavorite = (productId: string) => {
    return items.some((item) => item.product_id === productId)
  }

  const toggleFavorite = (product: { id: string; name: string; price: number }) => {
    if (isFavorite(product.id)) {
      removeFavorite(product.id)
    } else {
      addFavorite(product)
    }
  }

  return (
    <FavoritesContext.Provider value={{ items, addFavorite, removeFavorite, isFavorite, toggleFavorite }}>
      {children}
    </FavoritesContext.Provider>
  )
}

export const useFavorites = () => {
  const context = useContext(FavoritesContext)
  if (!context) throw new Error('useFavorites must be used within FavoritesProvider')
  return context
}