import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  FaShoppingCart,
  FaCreditCard,
  FaMoneyBillWave,
  FaWallet,
} from 'react-icons/fa'
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
      <div className="text-center py-16">
        <FaShoppingCart className="text-5xl text-gray-300 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Корзина пуста</h2>
        <button
          onClick={() => navigate('/catalog')}
          className="text-[#8A9A86] hover:underline font-medium inline-block"
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
      console.log('Заказ создан:', order)

      const paymentData = {
        order_id: order.id,
        amount: totalPrice,
        payment_method: form.paymentMethod,
      }

      const payment = await createPayment(paymentData)
      console.log('Платёж создан:', payment)

      let attempts = 0
      const maxAttempts = 10
      let paymentCompleted = false

      while (attempts < maxAttempts && !paymentCompleted) {
        await new Promise((resolve) => setTimeout(resolve, 1000))
        attempts++

        const statusResponse = await getPaymentStatus(payment.id)
        console.log(`Попытка ${attempts}: статус платежа - ${statusResponse.status}`)

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
      <h1 className="text-3xl font-bold text-[#1C1C1C] mb-6 flex items-center gap-2">
        <FaShoppingCart className="text-[#8A9A86]" />
        Оформление заказа
      </h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2 bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                Адрес доставки *
              </label>
              <input
                type="text"
                name="address"
                value={form.address}
                onChange={handleChange}
                placeholder="г. Минск, ул. Независимости, д. 10, кв. 25"
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                required
              />
            </div>

            <div className="grid grid-cols-3 gap-3">
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                  Подъезд
                </label>
                <input
                  type="text"
                  name="entrance"
                  value={form.entrance}
                  onChange={handleChange}
                  placeholder="1"
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                  Этаж
                </label>
                <input
                  type="text"
                  name="floor"
                  value={form.floor}
                  onChange={handleChange}
                  placeholder="5"
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                  Домофон
                </label>
                <input
                  type="text"
                  name="intercom"
                  value={form.intercom}
                  onChange={handleChange}
                  placeholder="25"
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                  Дата доставки *
                </label>
                <input
                  type="date"
                  name="deliveryDate"
                  value={form.deliveryDate}
                  onChange={handleChange}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                  Время доставки *
                </label>
                <select
                  name="deliveryTime"
                  value={form.deliveryTime}
                  onChange={handleChange}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C]"
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
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                Способ оплаты *
              </label>
              <select
                name="paymentMethod"
                value={form.paymentMethod}
                onChange={handleChange}
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C]"
                required
              >
                <option value="card">
                  <FaCreditCard /> Картой курьеру
                </option>
                <option value="cash">
                  <FaMoneyBillWave /> Наличными курьеру
                </option>
                <option value="online">
                  <FaWallet /> Онлайн на сайте
                </option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1.5">
                Комментарий к заказу
              </label>
              <textarea
                name="comment"
                value={form.comment}
                onChange={handleChange}
                placeholder="Позвоните за 15 минут до доставки"
                rows={3}
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] resize-none"
              />
            </div>

            {error && (
              <div className="bg-red-50 text-red-500 p-3 rounded-xl text-sm border border-red-100">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-[#8A9A86] text-white py-3 rounded-xl hover:bg-[#7A8A76] transition flex items-center justify-center gap-2 text-base font-medium disabled:opacity-50"
            >
              {loading ? 'Обработка...' : 'Оформить заказ'}
            </button>
          </form>
        </div>

        <div className="md:col-span-1">
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 sticky top-4 border border-gray-100">
            <h2 className="font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
              <FaShoppingCart className="text-[#8A9A86]" />
              Ваш заказ
            </h2>
            <div className="space-y-2 max-h-48 overflow-y-auto">
              {items.map((item) => (
                <div key={item.id} className="flex justify-between text-sm border-b border-gray-50 pb-1.5">
                  <span className="text-[#1C1C1C]">{item.name} × {item.quantity}</span>
                  <span className="text-[#1C1C1C] font-medium">{item.price * item.quantity} BYN</span>
                </div>
              ))}
            </div>
            <div className="border-t border-gray-100 pt-3 mt-3">
              <div className="flex justify-between font-semibold text-base">
                <span className="text-[#1C1C1C]">Итого</span>
                <span className="text-[#8A9A86]">{totalPrice} BYN</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CheckoutPage