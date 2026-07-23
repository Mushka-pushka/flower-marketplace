import { createContext, useContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { useAuth } from './AuthContext'

interface CartItem {
  id: string
  product_id: string
  name: string
  price: number
  quantity: number
  shop_id: string  
  image?: string
  stock?: number
}

interface CartContextType {
  items: CartItem[]
  addToCart: (product: { id: string; name: string; price: number; shop_id: string; image?: string }) => void
  removeFromCart: (productId: string) => void
  updateQuantity: (productId: string, quantity: number) => void
  clearCart: () => void
  totalItems: number
  totalPrice: number
}

const CartContext = createContext<CartContextType | undefined>(undefined)

export const CartProvider = ({ children }: { children: ReactNode }) => {
  const { user } = useAuth()
  const [items, setItems] = useState<CartItem[]>([])

  const getStorageKey = () => {
    return user ? `flower_cart_${user.id}` : 'flower_cart_guest'
  }

  useEffect(() => {
    const key = getStorageKey()
    
    // Очищаем старые ключи (без ID) ПРИ КАЖДОЙ ЗАГРУЗКЕ
    if (user) {
      const oldKeys = Object.keys(localStorage).filter(
        k => k.startsWith('flower_cart_') && k !== key
      )
      oldKeys.forEach(k => {
        console.log('Удаляем старый ключ корзины:', k)
        localStorage.removeItem(k)
      })
    }
    
    const saved = localStorage.getItem(key)
    console.log('Loading cart from:', key, saved)
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        if (Array.isArray(parsed) && parsed.length > 0) {
          setItems(parsed)
        }
        // Если parsed пустой — НЕ ОЧИЩАЕМ items
      } catch (error) {
        console.warn('Failed to parse cart data:', error)
        // НЕ ОЧИЩАЕМ items при ошибке
      }
    }
    // НЕ ОЧИЩАЕМ items, если данных нет
  }, [user])

  useEffect(() => {
    const key = getStorageKey()
    console.log('Saving cart to:', key, items)
    localStorage.setItem(key, JSON.stringify(items))
  }, [items, user])

  const addToCart = (product: { id: string; name: string; price: number; shop_id: string; image?: string }) => {
    setItems((prev) => {
      const existing = prev.find((item) => item.product_id === product.id)
      if (existing) {
        return prev.map((item) =>
          item.product_id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        )
      }
      return [
        ...prev,
        {
          id: crypto.randomUUID(),
          product_id: product.id,
          name: product.name,
          price: product.price,
          shop_id: product.shop_id,
          image: product.image, 
          quantity: 1,
        },
      ]
    })
  }

  const removeFromCart = (productId: string) => {
    setItems((prev) => {
      const newItems = prev.filter((item) => item.product_id !== productId)
      // Если корзина стала пустой — удаляем ключ из localStorage
      if (newItems.length === 0) {
        const key = getStorageKey()
        localStorage.removeItem(key)
        console.log('Cart empty, removed key:', key)
      }
      return newItems
    })
  }

  const updateQuantity = (productId: string, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(productId)
      return
    }
    setItems((prev) =>
      prev.map((item) =>
        item.product_id === productId ? { ...item, quantity } : item
      )
    )
  }

  const clearCart = () => {
    const key = getStorageKey()
    localStorage.removeItem(key)
    console.log('Cart cleared, removed key:', key)
    setItems([])
  }

  const totalItems = items.reduce((sum, item) => sum + item.quantity, 0)
  const totalPrice = items.reduce((sum, item) => sum + item.price * item.quantity, 0)

  return (
    <CartContext.Provider
      value={{
        items,
        addToCart,
        removeFromCart,
        updateQuantity,
        clearCart,
        totalItems,
        totalPrice,
      }}
    >
      {children}
    </CartContext.Provider>
  )
}

export const useCart = () => {
  const context = useContext(CartContext)
  if (!context) throw new Error('useCart must be used within CartProvider')
  return context
}