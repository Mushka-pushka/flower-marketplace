import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  FaShoppingCart,
  FaExclamationCircle,
  FaCreditCard,
  FaInfoCircle,
} from 'react-icons/fa'
import { useCart } from '../context/CartContext'
import { useAuth } from '../context/AuthContext'
import { createOrder } from '../api/order.api'
import { createPayment, getPaymentStatus } from '../api/payment.api'
import { createAddress } from '../api/catalog.api'

const CheckoutPage = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { items, totalPrice, clearCart } = useCart()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [paymentStatus, setPaymentStatus] = useState<'idle' | 'processing' | 'success' | 'failed'>('idle')
  const [shopId, setShopId] = useState<string | null>(null)

  const [form, setForm] = useState({
    address: '',
    entrance: '',
    floor: '',
    intercom: '',
    comment: '',
  })

  useEffect(() => {
    if (items.length === 0) return
    const firstItemShopId = items[0]?.shop_id
    if (!firstItemShopId) return
    const allSameShop = items.every(item => item.shop_id === firstItemShopId)
    if (!allSameShop) {
      setError('Все товары в корзине должны быть из одного магазина')
      setShopId(null)
      return
    }
    setShopId(firstItemShopId)
    setError('')
  }, [items])

  if (items.length === 0) {
    return (
      <div className="text-center py-16">
        <FaShoppingCart className="text-5xl text-gray-300 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Корзина пуста</h2>
        <button onClick={() => navigate('/catalog')} className="text-[#8A9A86] hover:underline font-medium inline-block">
          Перейти в каталог
        </button>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="text-center py-16">
        <FaExclamationCircle className="text-5xl text-amber-500 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Войдите в аккаунт</h2>
        <p className="text-gray-400 mb-4">Чтобы оформить заказ, войдите в свой аккаунт</p>
        <button onClick={() => navigate('/login')} className="text-[#8A9A86] hover:underline font-medium inline-block">
          Войти
        </button>
      </div>
    )
  }

  if (!shopId) {
    return (
      <div className="text-center py-16">
        <FaExclamationCircle className="text-5xl text-amber-500 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Не удалось определить магазин</h2>
        <p className="text-gray-400 mb-4">Попробуйте добавить товары заново</p>
        <button onClick={() => navigate('/catalog')} className="text-[#8A9A86] hover:underline font-medium inline-block">
          Перейти в каталог
        </button>
      </div>
    )
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    setPaymentStatus('processing')

    try {
      // 1. Создаём адрес доставки
      const addressData = {
        name: 'Доставка',
        address: form.address,
        entrance: form.entrance,
        floor: form.floor,
        intercom: form.intercom,
        comment: form.comment,
        is_default: false,
      }

      const address = await createAddress(addressData)
      console.log('Адрес создан:', address)

      // 2. Создаём заказ (оплата только онлайн — payment_type_id = 3)
      const orderData = {
        shop_id: shopId,
        delivery_address_id: address.id,
        payment_type_id: 3, // 3 = онлайн-оплата
        comment: form.comment,
        items: items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
        })),
      }

      const order = await createOrder(orderData)
      console.log('Заказ создан:', order)

      // 3. Создаём платёж (всегда онлайн)
      const paymentData = {
        order_id: order.id,
        amount: totalPrice,
        payment_method: 'online',
      }

      const payment = await createPayment(paymentData)
      console.log('Платёж создан:', payment)

      // 4. Ожидаем оплату
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
          setError('Оплата не прошла. Попробуйте оформить заказ снова.')
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

            {/* Уведомление о согласовании доставки */}
            <div className="bg-amber-50 border border-amber-200 rounded-xl p-4">
              <div className="flex items-start gap-3">
                <FaInfoCircle className="text-amber-600 text-lg mt-0.5 flex-shrink-0" />
                <div>
                  <p className="text-sm font-medium text-amber-800">Согласование доставки</p>
                  <p className="text-xs text-amber-700 mt-0.5">
                    После подтверждения заказа продавец свяжется с вами по указанному номеру телефона для согласования даты и времени доставки.
                  </p>
                </div>
              </div>
            </div>

            {/* Онлайн-оплата (информационный блок) */}
            <div className="bg-[#8A9A86]/5 border border-[#8A9A86]/20 rounded-xl p-4">
              <div className="flex items-start gap-3">
                <FaCreditCard className="text-[#8A9A86] text-lg mt-0.5 flex-shrink-0" />
                <div>
                  <p className="text-sm font-medium text-[#1C1C1C]">Онлайн-оплата</p>
                  <p className="text-xs text-gray-500 mt-0.5">
                    Оплата производится на сайте. После оформления заказа вы будете перенаправлены на страницу оплаты.
                  </p>
                </div>
              </div>
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
              {loading ? 'Обработка...' : 'Перейти к оплате'}
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