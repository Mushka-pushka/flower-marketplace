import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useCart } from '../context/CartContext'
import { createOrder } from '../api/order.api'
import { createPayment, getPaymentStatus } from '../api/payment.api'

const CheckoutPage = () => {
  const navigate = useNavigate()
  const { items, totalPrice, clearCart } = useCart()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [paymentStatus, setPaymentStatus] = useState<'idle' | 'processing' | 'success' | 'failed'>('idle')

  const [form, setForm] = useState({
    address: '',
    entrance: '',
    floor: '',
    intercom: '',
    comment: '',
    deliveryDate: '',
    deliveryTime: '',
    paymentMethod: 'card',
  })

  if (items.length === 0) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-gray-600">🛒 Корзина пуста</h2>
        <button
          onClick={() => navigate('/catalog')}
          className="text-pink-500 hover:underline mt-4 inline-block"
        >
          Перейти в каталог
        </button>
      </div>
    )
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    setPaymentStatus('processing')

    try {
      const orderData = {
        customer_id: '6b75b13b-2b7b-4df1-b700-b39ac0bc1d45',
        shop_id: '11111111-1111-1111-1111-111111111111',
        delivery_address_id: '11111111-1111-1111-1111-111111111111',
        payment_type_id: form.paymentMethod === 'card' ? 1 : form.paymentMethod === 'cash' ? 2 : 3,
        delivery_date: form.deliveryDate,
        delivery_time: form.deliveryTime,
        comment: form.comment,
        items: items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      }

      const order = await createOrder(orderData)
      console.log('✅ Заказ создан:', order)

      const paymentData = {
        order_id: order.id,
        amount: totalPrice,
        payment_method: form.paymentMethod,
      }

      const payment = await createPayment(paymentData)
      console.log('💳 Платёж создан:', payment)

      let attempts = 0
      const maxAttempts = 10
      let paymentCompleted = false

      while (attempts < maxAttempts && !paymentCompleted) {
        await new Promise((resolve) => setTimeout(resolve, 1000))
        attempts++

        const statusResponse = await getPaymentStatus(payment.id)
        console.log(`🔍 Попытка ${attempts}: статус платежа - ${statusResponse.status}`)

        if (statusResponse.status === 'completed') {
          paymentCompleted = true
          setPaymentStatus('success')
          clearCart()
          navigate('/checkout/success', { state: { orderId: order.id } })
          break
        }

        if (statusResponse.status === 'failed') {
          setPaymentStatus('failed')
          setError('Оплата не прошла. Попробуйте другой способ оплаты.')
          break
        }
      }

      if (!paymentCompleted && paymentStatus !== 'failed') {
        setError('Превышено время ожидания оплаты. Попробуйте снова.')
        setPaymentStatus('failed')
      }
    } catch (err: any) {
      console.error('Ошибка:', err)
      setError(err.response?.data?.error || 'Ошибка оформления заказа')
      setPaymentStatus('failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-3xl mx-auto animate-fade-in-up">
      <h1 className="text-4xl font-bold gradient-text mb-6">📦 Оформление заказа</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2 bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-6 border border-pink-50/50">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Адрес доставки *
              </label>
              <input
                type="text"
                name="address"
                value={form.address}
                onChange={handleChange}
                placeholder="г. Минск, ул. Независимости, д. 10, кв. 25"
                className="input-primary w-full px-4 py-2 rounded-lg"
                required
              />
            </div>

            <div className="grid grid-cols-3 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Подъезд
                </label>
                <input
                  type="text"
                  name="entrance"
                  value={form.entrance}
                  onChange={handleChange}
                  placeholder="1"
                  className="input-primary w-full px-4 py-2 rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Этаж
                </label>
                <input
                  type="text"
                  name="floor"
                  value={form.floor}
                  onChange={handleChange}
                  placeholder="5"
                  className="input-primary w-full px-4 py-2 rounded-lg"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Домофон
                </label>
                <input
                  type="text"
                  name="intercom"
                  value={form.intercom}
                  onChange={handleChange}
                  placeholder="25"
                  className="input-primary w-full px-4 py-2 rounded-lg"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Дата доставки *
                </label>
                <input
                  type="date"
                  name="deliveryDate"
                  value={form.deliveryDate}
                  onChange={handleChange}
                  className="input-primary w-full px-4 py-2 rounded-lg"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Время доставки *
                </label>
                <select
                  name="deliveryTime"
                  value={form.deliveryTime}
                  onChange={handleChange}
                  className="input-primary w-full px-4 py-2 rounded-lg"
                  required
                >
                  <option value="">Выберите время</option>
                  <option value="10:00-12:00">10:00 – 12:00</option>
                  <option value="12:00-14:00">12:00 – 14:00</option>
                  <option value="14:00-16:00">14:00 – 16:00</option>
                  <option value="16:00-18:00">16:00 – 18:00</option>
                  <option value="18:00-20:00">18:00 – 20:00</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Способ оплаты *
              </label>
              <select
                name="paymentMethod"
                value={form.paymentMethod}
                onChange={handleChange}
                className="input-primary w-full px-4 py-2 rounded-lg"
                required
              >
                <option value="card">Картой курьеру</option>
                <option value="cash">Наличными курьеру</option>
                <option value="online">Онлайн на сайте</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Комментарий к заказу
              </label>
              <textarea
                name="comment"
                value={form.comment}
                onChange={handleChange}
                placeholder="Позвоните за 15 минут до доставки"
                rows={3}
                className="input-primary w-full px-4 py-2 rounded-lg"
              />
            </div>

            {error && (
              <div className="bg-red-50 text-red-600 p-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="btn-primary w-full py-3 rounded-full text-lg font-medium"
            >
              {loading ? 'Обработка...' : '✅ Оформить заказ'}
            </button>
          </form>
        </div>

        <div className="md:col-span-1">
          <div className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-6 sticky top-4 border border-pink-50/50">
            <h2 className="font-semibold text-gray-700 mb-4">🛒 Ваш заказ</h2>
            <div className="space-y-2 max-h-48 overflow-y-auto">
              {items.map((item) => (
                <div key={item.id} className="flex justify-between text-sm">
                  <span className="text-gray-600">{item.name} × {item.quantity}</span>
                  <span className="text-gray-800 font-medium">{item.price * item.quantity} BYN</span>
                </div>
              ))}
            </div>
            <div className="border-t pt-3 mt-3 border-pink-100">
              <div className="flex justify-between font-semibold text-lg">
                <span>Итого</span>
                <span className="text-pink-600">{totalPrice} BYN</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CheckoutPage