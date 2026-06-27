import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useCart } from '../context/CartContext'
import ProductModal from '../components/ProductModal'

const CartPage = () => {
  const { items, removeFromCart, updateQuantity, totalPrice, clearCart } = useCart()
  const [selectedIds, setSelectedIds] = useState<string[]>(items.map(item => item.product_id))
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)

  const toggleSelect = (productId: string) => {
    setSelectedIds(prev =>
      prev.includes(productId)
        ? prev.filter(id => id !== productId)
        : [...prev, productId]
    )
  }

  const toggleSelectAll = () => {
    if (selectedIds.length === items.length) {
      setSelectedIds([])
    } else {
      setSelectedIds(items.map(item => item.product_id))
    }
  }

  const selectedItems = items.filter(item => selectedIds.includes(item.product_id))
  const selectedTotal = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0)

  const openModal = (productId: string) => {
    setSelectedProductId(productId)
    document.body.style.overflow = 'hidden'
  }

  const closeModal = () => {
    setSelectedProductId(null)
    document.body.style.overflow = 'auto'
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-gray-600">🛒 Корзина пуста</h2>
        <Link to="/catalog" className="text-pink-500 hover:underline mt-4 inline-block">
          Перейти в каталог
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">🛒 Корзина</h1>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        {/* Шапка таблицы */}
        <div className="flex items-center gap-4 p-4 bg-gray-50 border-b">
          <input
            type="checkbox"
            checked={selectedIds.length === items.length}
            onChange={toggleSelectAll}
            className="w-5 h-5 accent-pink-500"
          />
          <span className="text-sm font-medium text-gray-600">Выбрать всё</span>
          <span className="text-sm text-gray-400 ml-auto">Товар</span>
          <span className="text-sm text-gray-400 w-24 text-center">Кол-во</span>
          <span className="text-sm text-gray-400 w-24 text-center">Цена</span>
          <span className="w-8" />
        </div>

        {items.map((item) => {
          const isSelected = selectedIds.includes(item.product_id)
          return (
            <div
              key={item.id}
              className={`flex items-center gap-4 p-4 border-b last:border-b-0 transition cursor-pointer ${
                isSelected ? 'bg-pink-50' : ''
              }`}
              onClick={() => openModal(item.product_id)}
            >
                <div onClick={(e) => e.stopPropagation()}>
              <input
                type="checkbox"
                checked={isSelected}
                onChange={() => toggleSelect(item.product_id)}
                className="w-5 h-5 accent-pink-500"
              />
            </div>

              <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-2xl flex-shrink-0">
                🌸
              </div>

              <div className="flex-1 min-w-0">
                <h3 className="font-semibold text-gray-800 truncate">{item.name}</h3>
                <p className="text-sm text-gray-500">BYN {item.price}</p>
              </div>

              <div className="flex items-center gap-2 w-24 justify-center" onClick={(e) => e.stopPropagation()}>
                <button
                  onClick={() => updateQuantity(item.product_id, item.quantity - 1)}
                  className="w-8 h-8 border border-gray-300 rounded hover:bg-gray-50"
                >
                  −
                </button>
                <span className="w-8 text-center font-medium">{item.quantity}</span>
                <button
                  onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                  className="w-8 h-8 border border-gray-300 rounded hover:bg-gray-50"
                >
                  +
                </button>
              </div>

              <div className="w-24 text-center font-medium text-pink-600">
                {item.price * item.quantity} BYN
              </div>

              <button
                onClick={(e) => {
                  e.stopPropagation()
                  removeFromCart(item.product_id)
                }}
                className="text-red-400 hover:text-red-600 text-sm w-8 text-center"
              >
                ✕
              </button>
            </div>
          )
        })}
      </div>

      {/* Итог */}
      <div className="mt-6 bg-white rounded-lg shadow p-4 flex flex-wrap justify-between items-center gap-4">
        <div>
          <span className="text-gray-600">Выбрано товаров: <strong>{selectedIds.length}</strong></span>
          <span className="text-lg font-semibold ml-6">
            Итого: <span className="text-pink-600">{selectedTotal} BYN</span>
          </span>
        </div>
        <div className="flex gap-4">
          <button
            onClick={clearCart}
            className="text-gray-400 hover:text-red-500 text-sm"
          >
            Очистить корзину
          </button>
          <Link
            to={selectedIds.length > 0 ? '/checkout' : '#'}
            className={`px-6 py-2 rounded-lg transition ${
              selectedIds.length > 0
                ? 'bg-pink-500 text-white hover:bg-pink-600'
                : 'bg-gray-300 text-gray-500 cursor-not-allowed'
            }`}
          >
            Оформить выбранные
          </Link>
        </div>
      </div>

      {/* Модальное окно */}
      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </div>
  )
}

export default CartPage