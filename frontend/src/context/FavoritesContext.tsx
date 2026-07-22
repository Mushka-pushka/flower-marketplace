import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from './AuthContext'

interface FavoriteItem {
  id: string
  product_id: string
  name: string
  price: number
  image?: string
}

interface FavoritesContextType {
  items: FavoriteItem[]
  loading: boolean
  addFavorite: (product: { id: string; name: string; price: number; image?: string }) => void 
  removeFavorite: (productId: string) => void
  isFavorite: (productId: string) => boolean
  toggleFavorite: (product: { id: string; name: string; price: number; image?: string }) => void  
}

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined)

export const FavoritesProvider = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth()
  const [items, setItems] = useState<FavoriteItem[]>([])
  const [loading, setLoading] = useState(true)

  const getStorageKey = () => {
    return user ? `flower_favorites_${user.id}` : 'flower_favorites_guest'
  }

  useEffect(() => {
    const loadFavorites = async () => {
      setLoading(true)
      const key = getStorageKey()
      const oldKey = 'flower_favorites' // старый ключ без ID
      
      let saved = localStorage.getItem(key)
      
      if (!saved) {
        const oldData = localStorage.getItem(oldKey)
        if (oldData) {
          // Переносим данные из старого ключа в новый
          localStorage.setItem(key, oldData)
          saved = oldData
          // Удаляем старый ключ (опционально)
          // localStorage.removeItem(oldKey)
        }
      }
      
      if (saved) {
        try {
          setItems(JSON.parse(saved))
        } catch {
          setItems([])
        }
      } else {
        setItems([])
      }
      setLoading(false)
    }

    loadFavorites()
  }, [user])

  useEffect(() => {
    const key = getStorageKey()
    localStorage.setItem(key, JSON.stringify(items))
  }, [items, user])

  const addFavorite = (product: { id: string; name: string; price: number; image?: string }) => {
    setItems((prev) => {
      if (prev.some((item) => item.product_id === product.id)) return prev
      return [
        ...prev,
        {
          id: crypto.randomUUID(),
          product_id: product.id,
          name: product.name,
          price: product.price,
          image: product.image || '',
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

  const toggleFavorite = (product: { id: string; name: string; price: number; image?: string }) => {
    if (isFavorite(product.id)) {
      removeFavorite(product.id)
    } else {
      addFavorite(product)
    }
  }

  return (
    <FavoritesContext.Provider 
      value={{ 
        items, 
        loading, 
        addFavorite, 
        removeFavorite, 
        isFavorite, 
        toggleFavorite 
      }}
    >
      {children}
    </FavoritesContext.Provider>
  )
}

export const useFavorites = () => {
  const context = useContext(FavoritesContext)
  if (!context) throw new Error('useFavorites must be used within FavoritesProvider')
  return context
}