import { Link } from 'react-router-dom'
import { useCart } from '../context/CartContext'

const CartPage = () => {
  const { items, removeFromCart, updateQuantity, totalPrice, clearCart } = useCart()

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
        {items.map((item) => (
          <div key={item.id} className="flex items-center gap-4 p-4 border-b last:border-b-0">
            <div className="w-16 h-16 bg-gray-100 rounded flex items-center justify-center text-2xl">
              🌸
            </div>
            
            <div className="flex-1">
              <h3 className="font-semibold text-gray-800">{item.name}</h3>
              <p className="text-pink-600 font-bold">{item.price} BYN</p>
            </div>

            <div className="flex items-center gap-2">
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

            <button
              onClick={() => removeFromCart(item.product_id)}
              className="text-red-400 hover:text-red-600 text-sm"
            >
              ✕
            </button>
          </div>
        ))}
      </div>

      <div className="mt-6 bg-white rounded-lg shadow p-4 flex justify-between items-center">
        <span className="text-lg font-semibold">
          Итого: <span className="text-pink-600">{totalPrice} BYN</span>
        </span>
        <div className="flex gap-4">
          <button
            onClick={clearCart}
            className="text-gray-400 hover:text-red-500 text-sm"
          >
            Очистить корзину
          </button>
          <Link
            to="/checkout"
            className="bg-pink-500 text-white px-6 py-2 rounded-lg hover:bg-pink-600 transition"
          >
            Оформить заказ
          </Link>
        </div>
      </div>
    </div>
  )
}

export default CartPage