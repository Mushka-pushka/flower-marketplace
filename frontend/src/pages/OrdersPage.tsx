import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import {
  FaBox,
  FaTimes,
  FaShoppingBag,
  FaChevronLeft,
  FaChevronRight,
  FaStar,
  FaRegStar,
} from 'react-icons/fa'
import { getMyOrderItems, getOrderDetails, cancelOrder } from '../api/order.api'
import type { OrderItemWithStatus, OrderDetails } from '../api/order.api'
import OrderTimeline from '../components/OrderTimeline'
import { useAuth } from '../context/AuthContext'
import { toast } from 'react-hot-toast'
import { createReview } from '../api/catalog.api'
import client from '../api/client'
import ProductModal from '../components/ProductModal'

const OrdersPage = () => {
  const { user } = useAuth()
  const [orderItems, setOrderItems] = useState<OrderItemWithStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [selectedOrder, setSelectedOrder] = useState<OrderDetails | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [limit] = useState(10)
  const [offset, setOffset] = useState(0)
  const [total, setTotal] = useState(0)

  // Состояния для модалки отзыва
  const [reviewModal, setReviewModal] = useState<{
    productId: string
    productName: string
  } | null>(null)
  const [reviewRating, setReviewRating] = useState(5)
  const [reviewComment, setReviewComment] = useState('')
  const [submittingReview, setSubmittingReview] = useState(false)
  const [hasReviewed, setHasReviewed] = useState<Record<string, boolean>>({})

  // Состояние для модалки товара из каталога
  const [productModalId, setProductModalId] = useState<string | null>(null)

  const fetchOrderItems = async () => {
    if (!user) {
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      const data = await getMyOrderItems()
      console.log('Order items data:', data)
      if (Array.isArray(data)) {
        setOrderItems(data)
        setTotal(data.length)
      } else {
        setOrderItems([])
        setTotal(0)
      }
    } catch (err: any) {
      console.error('Ошибка загрузки заказов:', err)
      setError(err.response?.data?.error || 'Ошибка загрузки заказов')
      setOrderItems([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  // Загружаем список отзывов пользователя
  const fetchUserReviews = async () => {
    try {
      const response = await client.get('/catalog/reviews/me')
      const reviewMap: Record<string, boolean> = {}
      response.data.forEach((review: any) => {
        reviewMap[review.product_id] = true
      })
      setHasReviewed(reviewMap)
    } catch (error) {
      console.error('Ошибка загрузки отзывов пользователя:', error)
    }
  }

  useEffect(() => {
    fetchOrderItems()
  }, [user, limit, offset])

  useEffect(() => {
    if (isModalOpen) {
      document.body.style.overflow = 'hidden'
      // Загружаем отзывы пользователя при открытии модалки
      if (user) {
        fetchUserReviews()
      }
    } else {
      document.body.style.overflow = 'auto'
    }
    return () => {
      document.body.style.overflow = 'auto'
    }
  }, [isModalOpen, user])

  // Блокировка скролла при открытии модалки товара
  useEffect(() => {
    if (productModalId) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'auto'
    }
    return () => {
      document.body.style.overflow = 'auto'
    }
  }, [productModalId])

  const handleOrderClick = async (orderId: string) => {
    try {
      const details = await getOrderDetails(orderId)
      console.log('Order details:', details)
      console.log('Items:', details.items)
      setSelectedOrder(details)
      setIsModalOpen(true)
    } catch (err) {
      console.error('Ошибка загрузки деталей заказа:', err)
      toast.error('Не удалось загрузить детали заказа')
    }
  }

  const handleCancelOrder = async (orderId: string) => {
    if (!confirm('Вы уверены, что хотите отменить заказ?')) return
    try {
      await cancelOrder(orderId)
      toast.success('Заказ отменен')
      setIsModalOpen(false)
      await fetchOrderItems()
    } catch (error) {
      console.error('Ошибка отмены заказа:', error)
      toast.error('Не удалось отменить заказ')
    }
  }

  // Обработчики для отзыва из заказа
  const handleOpenReviewModal = (productId: string, productName: string) => {
    setReviewModal({ productId, productName })
    setReviewRating(5)
    setReviewComment('')
  }

  const handleSubmitReview = async () => {
    if (!reviewModal) return
    if (!reviewComment.trim()) {
      toast.error('Пожалуйста, напишите комментарий')
      return
    }
    setSubmittingReview(true)
    try {
      await createReview({
        product_id: reviewModal.productId,
        rating: reviewRating,
        comment: reviewComment.trim(),
      })
      toast.success('Спасибо за ваш отзыв!')
      setReviewModal(null)
      // Обновляем список отзывов
      await fetchUserReviews()
      await fetchOrderItems()
    } catch (error) {
      console.error('Ошибка отправки отзыва:', error)
      toast.error('Не удалось отправить отзыв')
    } finally {
      setSubmittingReview(false)
    }
  }

  // Обработчики для модалки товара
  const handleViewProduct = (productId: string) => {
    setProductModalId(productId)
  }

  const closeProductModal = () => {
    setProductModalId(null)
  }

  const getStatusLabel = (status: string) => {
    const map: Record<string, string> = {
      pending: 'Ожидает подтверждения',
      paid: 'Оплачен',
      confirmed: 'Подтверждён',
      preparing: 'Собирается',
      packing: 'Упаковывается',
      delivery: 'В доставке',
      delivered: 'Доставлен',
      cancelled: 'Отменён',
    }
    return map[status] || status
  }

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'text-yellow-600 bg-yellow-50 border-yellow-200',
      paid: 'text-blue-600 bg-blue-50 border-blue-200',
      confirmed: 'text-blue-600 bg-blue-50 border-blue-200',
      preparing: 'text-purple-600 bg-purple-50 border-purple-200',
      packing: 'text-indigo-600 bg-indigo-50 border-indigo-200',
      delivery: 'text-orange-600 bg-orange-50 border-orange-200',
      delivered: 'text-green-600 bg-green-50 border-green-200',
      cancelled: 'text-red-600 bg-red-50 border-red-200',
    }
    return map[status] || 'text-gray-600 bg-gray-50 border-gray-200'
  }

  const totalPages = Math.ceil(total / limit)
  const currentPage = Math.floor(offset / limit) + 1

  const goToPreviousPage = () => {
    setOffset(Math.max(0, offset - limit))
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const goToNextPage = () => {
    setOffset(offset + limit)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка заказов...</div>
  }

  if (error) {
    return <div className="text-center py-8 text-red-500">{error}</div>
  }

  if (!orderItems || orderItems.length === 0) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
          <FaBox className="text-[#8A9A86]" />
          Мои заказы
        </h2>
        <p className="text-gray-400 text-base">У вас пока нет заказов</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
        <FaBox className="text-[#8A9A86]" />
        Мои заказы
      </h2>
      <div className="space-y-3">
        {orderItems.map((item) => (
          <div
            key={item.id}
            onClick={() => handleOrderClick(item.order_id)}
            className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] hover:shadow-[0_8px_30px_rgba(0,0,0,0.08)] transition-all duration-300 p-4 cursor-pointer border border-gray-100"
          >
            <div className="flex justify-between items-start">
              <div>
                <p className="text-sm font-medium text-[#1C1C1C]">{item.product_name}</p>
                <p className="text-sm text-gray-400">
                  {new Date(item.created_at).toLocaleDateString('ru-RU', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                  })}
                </p>
                <p className="text-xs text-gray-400 mt-0.5">
                  Количество: {item.quantity} шт.
                </p>
              </div>
              <div className="text-right">
                <span className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(item.order_status)}`}>
                  {getStatusLabel(item.order_status)}
                </span>
                <p className="text-lg font-bold text-[#8A9A86] mt-1">{item.total} BYN</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {total > limit && (
        <div className="flex justify-center items-center gap-2 mt-6">
          <button
            onClick={goToPreviousPage}
            disabled={offset === 0}
            className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition flex items-center gap-1"
          >
            <FaChevronLeft className="text-sm" />
            Назад
          </button>
          <span className="px-4 py-2 text-[#1C1C1C] font-medium">
            {currentPage} / {totalPages}
          </span>
          <button
            onClick={goToNextPage}
            disabled={offset + limit >= total}
            className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition flex items-center gap-1"
          >
            Вперед
            <FaChevronRight className="text-sm" />
          </button>
        </div>
      )}

      {isModalOpen &&
        selectedOrder &&
        createPortal(
          <div
            className="fixed inset-0 z-[9999] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
            onClick={() => setIsModalOpen(false)}
          >
            <div
              className="bg-white rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 shadow-[0_8px_40px_rgba(0,0,0,0.08)] border border-gray-100 animate-fade-in-up"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="flex justify-between items-center mb-4 pb-3 border-b border-gray-100">
                <h2 className="text-xl font-bold text-[#1C1C1C] flex items-center gap-2">
                  <FaShoppingBag className="text-[#8A9A86]" />
                  {selectedOrder.items.length > 0
                    ? selectedOrder.items.map((item) => item.product_name || 'Товар').join(', ')
                    : `Заказ #${selectedOrder.order.id.slice(0, 8)}`}
                </h2>
                <button
                  onClick={() => setIsModalOpen(false)}
                  className="text-gray-400 hover:text-gray-600 text-2xl leading-none transition flex-shrink-0"
                >
                  <FaTimes />
                </button>
              </div>

              <div className="mb-4">
                <p className="text-sm text-gray-400">
                  Сумма: <span className="font-bold text-[#8A9A86]">{selectedOrder.order.total_amount} BYN</span>
                </p>
                <p className="text-sm text-gray-400">
                  Статус:{' '}
                  <span
                    className={`px-3 py-1 rounded-full text-xs font-medium border ${getStatusColor(
                      selectedOrder.order.current_status
                    )}`}
                  >
                    {getStatusLabel(selectedOrder.order.current_status)}
                  </span>
                </p>
              </div>

              <div className="border-t border-gray-100 pt-4">
                <h3 className="font-semibold text-[#1C1C1C] mb-3">История статусов</h3>
                <OrderTimeline
                  statuses={selectedOrder.statuses}
                  orderId={selectedOrder.order.id}
                  currentStatus={selectedOrder.order.current_status}
                  onStatusUpdate={() => {
                    fetchOrderItems()
                  }}
                />
              </div>

              <div className="border-t border-gray-100 pt-4 mt-4">
                <h3 className="font-semibold text-[#1C1C1C] mb-2">Товары</h3>
                <div className="space-y-2">
                  {selectedOrder.items.map((item) => {
                    const isDelivered = selectedOrder.order.current_status === 'delivered'
                    const alreadyReviewed = hasReviewed[item.product_id]
                    const canReview = isDelivered && !alreadyReviewed

                    return (
                      <div
                        key={item.id}
                        className="flex justify-between items-center text-sm border-b border-gray-100 pb-2 last:border-0"
                      >
                        <div className="flex-1">
                          <span className="text-gray-400">
                            {item.quantity} × {item.price} BYN
                          </span>
                          <span className="text-[#1C1C1C] font-medium ml-2">
                            {item.total} BYN
                          </span>
                          {item.product_name && (
                            <span className="text-[#1C1C1C] ml-2 font-medium">
                              — {item.product_name}
                            </span>
                          )}
                        </div>

                        <div className="flex items-center gap-2 ml-4">
                          {canReview && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation()
                                handleOpenReviewModal(item.product_id, item.product_name || 'Товар')
                              }}
                              className="px-3 py-1 text-xs bg-[#8A9A86] text-white rounded-lg hover:bg-[#7A8A76] transition whitespace-nowrap"
                            >
                              Оставить отзыв
                            </button>
                          )}
                          
                          {/* Кнопка "Посмотреть товар" */}
                          {isDelivered && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation()
                                handleViewProduct(item.product_id)
                              }}
                              className="px-3 py-1 text-xs bg-gray-100 text-[#1C1C1C] rounded-lg hover:bg-gray-200 transition whitespace-nowrap"
                            >
                              Посмотреть товар
                            </button>
                          )}
                          
                          {isDelivered && alreadyReviewed && (
                            <span className="text-xs text-green-600 whitespace-nowrap">
                              ✓ Отзыв оставлен
                            </span>
                          )}
                        </div>
                      </div>
                    )
                  })}
                </div>
              </div>

              {selectedOrder.order.current_status === 'pending' && (
                <div className="border-t border-gray-100 pt-4 mt-4">
                  <button
                    onClick={() => handleCancelOrder(selectedOrder.order.id)}
                    className="w-full px-6 py-2.5 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition border border-red-200 text-sm font-medium"
                  >
                    Отменить заказ
                  </button>
                </div>
              )}
            </div>
          </div>,
          document.body
        )}

      {/* Модалка для оставления отзыва из заказа */}
      {reviewModal &&
        createPortal(
          <div
            className="fixed inset-0 z-[10000] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
            onClick={() => setReviewModal(null)}
          >
            <div
              className="bg-white rounded-2xl w-full max-w-md p-6 shadow-lg border border-gray-100 animate-fade-in-up"
              onClick={(e) => e.stopPropagation()}
            >
              <h3 className="text-xl font-bold text-[#1C1C1C] mb-2">
                Оставить отзыв
              </h3>
              <p className="text-sm text-gray-400 mb-4">
                на товар «{reviewModal.productName}»
              </p>

              <div className="flex items-center gap-1 mb-3">
                <span className="text-sm text-gray-400 mr-2">Оценка:</span>
                {[1, 2, 3, 4, 5].map((star) => (
                  <button
                    key={star}
                    type="button"
                    onClick={() => setReviewRating(star)}
                    className="text-2xl hover:scale-110 transition focus:outline-none"
                  >
                    {star <= reviewRating ? (
                      <FaStar className="text-[#8A9A86]" />
                    ) : (
                      <FaRegStar className="text-gray-300" />
                    )}
                  </button>
                ))}
              </div>

              <textarea
                value={reviewComment}
                onChange={(e) => setReviewComment(e.target.value)}
                placeholder="Поделитесь впечатлениями о товаре..."
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition resize-none text-sm"
                rows={4}
              />

              <div className="flex gap-3 mt-4">
                <button
                  onClick={() => setReviewModal(null)}
                  className="flex-1 px-4 py-2.5 border border-gray-200 rounded-xl hover:bg-gray-50 transition text-sm font-medium"
                >
                  Отмена
                </button>
                <button
                  onClick={handleSubmitReview}
                  disabled={submittingReview || !reviewComment.trim()}
                  className="flex-1 bg-[#8A9A86] text-white px-4 py-2.5 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {submittingReview ? 'Отправка...' : 'Отправить отзыв'}
                </button>
              </div>
            </div>
          </div>,
          document.body
        )}

      {/* Модалка товара из каталога */}
      {productModalId && (
        <ProductModal productId={productModalId} onClose={closeProductModal} />
      )}
    </div>
  )
}

export default OrdersPage