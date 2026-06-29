import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import {
  FaShoppingCart,
  FaTrash,
  FaPlus,
  FaMinus,
  FaLeaf,
  FaExclamationCircle,
} from 'react-icons/fa'
import { useCart } from '../context/CartContext'
import { useAuth } from '../context/AuthContext'
import ProductModal from '../components/ProductModal'

const CartPage = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { items, removeFromCart, updateQuantity, clearCart } = useCart()
  const [selectedIds, setSelectedIds] = useState<string[]>(items.map(item => item.product_id))
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null)
  const [authNotification, setAuthNotification] = useState(false)

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

  const handleCheckout = () => {
    if (!user) {
      setAuthNotification(true)
      setTimeout(() => setAuthNotification(false), 3000)
      return
    }
    navigate('/checkout')
  }

  if (items.length === 0) {
    return (
      <div className="text-center py-16">
        <FaShoppingCart className="text-5xl text-gray-300 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Корзина пуста</h2>
        <p className="text-gray-400 mb-4">Добавьте товары из каталога</p>
        <Link to="/catalog" className="text-[#8A9A86] hover:underline font-medium inline-block">
          Перейти в каталог
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto animate-fade-in-up">
      {/* Уведомление о необходимости авторизации */}
      {authNotification && (
        <div className="mb-4 bg-amber-50 text-amber-700 px-4 py-3 rounded-xl flex items-center gap-2 border border-amber-200">
          <FaExclamationCircle />
          <span>Войдите в аккаунт, чтобы оформить заказ</span>
        </div>
      )}

      <h1 className="text-3xl font-bold text-[#1C1C1C] mb-6 flex items-center gap-2">
        <FaShoppingCart className="text-[#8A9A86]" />
        Корзина
      </h1>

      <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] overflow-hidden border border-gray-100">
        {/* Шапка таблицы */}
        <div className="flex items-center gap-4 p-4 bg-gray-50/50 border-b border-gray-100">
          <input
            type="checkbox"
            checked={selectedIds.length === items.length}
            onChange={toggleSelectAll}
            className="w-5 h-5 accent-[#8A9A86] rounded"
          />
          <span className="text-sm font-medium text-[#1C1C1C]">Выбрать всё</span>
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
                isSelected ? 'bg-[#8A9A86]/5' : 'hover:bg-gray-50/50'
              }`}
              onClick={() => openModal(item.product_id)}
            >
              <div onClick={(e) => e.stopPropagation()}>
                <input
                  type="checkbox"
                  checked={isSelected}
                  onChange={() => toggleSelect(item.product_id)}
                  className="w-5 h-5 accent-[#8A9A86] rounded"
                />
              </div>

              <div className="w-16 h-16 bg-gray-50 rounded-xl flex items-center justify-center text-2xl flex-shrink-0">
                <FaLeaf className="text-gray-300 text-2xl" />
              </div>

              <div className="flex-1 min-w-0">
                <h3 className="font-semibold text-[#1C1C1C] truncate text-base">{item.name}</h3>
                <p className="text-sm text-gray-400">{item.price} BYN</p>
              </div>

              <div className="flex items-center gap-2 w-24 justify-center" onClick={(e) => e.stopPropagation()}>
                <button
                  onClick={() => updateQuantity(item.product_id, item.quantity - 1)}
                  className="w-8 h-8 border border-gray-200 rounded-full hover:border-[#8A9A86] hover:bg-[#8A9A86]/5 transition flex items-center justify-center"
                >
                  <FaMinus className="text-xs text-[#1C1C1C]" />
                </button>
                <span className="w-8 text-center font-medium text-[#1C1C1C]">{item.quantity}</span>
                <button
                  onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                  className="w-8 h-8 border border-gray-200 rounded-full hover:border-[#8A9A86] hover:bg-[#8A9A86]/5 transition flex items-center justify-center"
                >
                  <FaPlus className="text-xs text-[#1C1C1C]" />
                </button>
              </div>

              <div className="w-24 text-center font-medium text-[#8A9A86]">
                {item.price * item.quantity} BYN
              </div>

              <button
                onClick={(e) => {
                  e.stopPropagation()
                  removeFromCart(item.product_id)
                }}
                className="text-gray-300 hover:text-red-500 text-sm w-8 text-center transition"
              >
                <FaTrash />
              </button>
            </div>
          )
        })}
      </div>

      {/* Итог */}
      <div className="mt-6 bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 flex flex-wrap justify-between items-center gap-4 border border-gray-100">
        <div>
          <span className="text-[#1C1C1C]">
            Выбрано товаров: <strong>{selectedIds.length}</strong>
          </span>
          <span className="text-lg font-semibold ml-6 text-[#1C1C1C]">
            Итого: <span className="text-[#8A9A86]">{selectedTotal} BYN</span>
          </span>
        </div>
        <div className="flex gap-4">
          <button
            onClick={clearCart}
            className="text-gray-400 hover:text-red-500 text-sm transition flex items-center gap-1.5"
          >
            <FaTrash /> Очистить
          </button>
          <button
            onClick={handleCheckout}
            className={`px-6 py-2.5 rounded-xl transition flex items-center gap-2 text-sm font-medium ${
              selectedIds.length > 0
                ? 'bg-[#8A9A86] text-white hover:bg-[#7A8A76]'
                : 'bg-gray-100 text-gray-400 cursor-not-allowed'
            }`}
            disabled={selectedIds.length === 0}
          >
            Оформить выбранные
          </button>
        </div>
      </div>

      <ProductModal productId={selectedProductId} onClose={closeModal} />
    </div>
  )
}

export default CartPage