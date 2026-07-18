import { useState, useEffect } from 'react'
import { FaStar, FaRegStar, FaUser } from 'react-icons/fa'
import { getProductReviews, createReview, updateReview, deleteReview } from '../api/catalog.api'
import { canReviewProduct } from '../api/order.api'
import type { Review } from '../api/catalog.api'
import { useAuth } from '../context/AuthContext'
import { toast } from 'react-hot-toast'

interface ReviewsProps {
  productId: string
}

const Reviews = ({ productId }: ReviewsProps) => {
  const { user } = useAuth()
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(true)
  const [rating, setRating] = useState(5)
  const [comment, setComment] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [canReview, setCanReview] = useState(false)
  const [editingReview, setEditingReview] = useState<string | null>(null)
  const [editRating, setEditRating] = useState(5)
  const [editComment, setEditComment] = useState('')

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [reviewsData, canReviewData] = await Promise.all([
          getProductReviews(productId),
          user ? canReviewProduct(productId) : Promise.resolve(false),
        ])
        setReviews(reviewsData || [])
        setCanReview(canReviewData)
      } catch (error) {
        console.error('Ошибка загрузки данных:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [productId, user])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!user) {
      toast.error('Войдите в аккаунт, чтобы оставить отзыв')
      return
    }

    if (!canReview) {
      toast.error('Вы можете оставить отзыв только на товары, которые заказывали и получили')
      return
    }

    setSubmitting(true)
    try {
      await createReview({
        product_id: productId,
        rating,
        comment,
      })
      const data = await getProductReviews(productId)
      setReviews(data || [])
      setRating(5)
      setComment('')
      toast.success('Спасибо за ваш отзыв!')
    } catch (error) {
      console.error('Ошибка отправки отзыва:', error)
      toast.error('Не удалось отправить отзыв')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteReview = async (reviewId: string) => {
    if (!confirm('Удалить отзыв?')) return
    try {
      await deleteReview(reviewId)
      const data = await getProductReviews(productId)
      setReviews(data || [])
      toast.success('Отзыв удален')
    } catch (error) {
      console.error('Ошибка удаления отзыва:', error)
      toast.error('Не удалось удалить отзыв')
    }
  }

  const handleEditReview = async (reviewId: string) => {
    try {
      await updateReview(reviewId, {
        rating: editRating,
        comment: editComment,
      })
      const data = await getProductReviews(productId)
      setReviews(data || [])
      setEditingReview(null)
      toast.success('Отзыв обновлен')
    } catch (error) {
      console.error('Ошибка обновления отзыва:', error)
      toast.error('Не удалось обновить отзыв')
    }
  }

  const cancelEdit = () => {
    setEditingReview(null)
    setEditRating(5)
    setEditComment('')
  }

  if (loading) {
    return <div className="text-gray-400 text-sm">Загрузка отзывов...</div>
  }

  const averageRating = reviews.length > 0
    ? reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length
    : 0

  return (
    <div className="mt-4">
      <div className="flex items-center gap-2 mb-3">
        <div className="flex items-center gap-0.5">
          {[1, 2, 3, 4, 5].map((star) => (
            <span key={star} className="text-lg">
              {star <= Math.round(averageRating) ? (
                <FaStar className="text-[#8A9A86]" />
              ) : (
                <FaRegStar className="text-gray-300" />
              )}
            </span>
          ))}
        </div>
        <span className="text-sm font-medium text-[#1C1C1C]">
          {averageRating.toFixed(1)}
        </span>
        <span className="text-sm text-gray-400">
          ({reviews.length} отзывов)
        </span>
      </div>

      {/* Форма отправки отзыва — только если пользователь может оставить отзыв */}
      {user && canReview && (
        <form onSubmit={handleSubmit} className="mb-4 p-4 bg-gray-50 rounded-xl border border-gray-100">
          <h4 className="font-medium text-[#1C1C1C] mb-2">Оставить отзыв</h4>
          <div className="flex items-center gap-1 mb-2">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                type="button"
                onClick={() => setRating(star)}
                className="text-2xl hover:scale-110 transition"
              >
                {star <= rating ? (
                  <FaStar className="text-[#8A9A86]" />
                ) : (
                  <FaRegStar className="text-gray-300" />
                )}
              </button>
            ))}
          </div>
          <textarea
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            placeholder="Поделитесь впечатлениями о товаре..."
            className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-sm resize-none"
            rows={2}
            required
          />
          <button
            type="submit"
            disabled={submitting}
            className="mt-2 bg-[#8A9A86] text-white px-4 py-1.5 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50"
          >
            {submitting ? 'Отправка...' : 'Отправить отзыв'}
          </button>
        </form>
      )}

      {user && !canReview && (
        <p className="text-sm text-gray-400 mb-3">
          💡 Вы можете оставить отзыв только после получения заказа с этим товаром
        </p>
      )}

      <div className="space-y-3 max-h-48 overflow-y-auto">
        {reviews.length === 0 ? (
          <p className="text-gray-400 text-sm">Пока нет отзывов</p>
        ) : (
          reviews.map((review) => (
            <div key={review.id} className="border-b border-gray-100 pb-2 last:border-0">
              {editingReview === review.id ? (
                // Режим редактирования
                <div className="p-3 bg-gray-50 rounded-xl border border-gray-100">
                  <div className="flex items-center gap-1 mb-2">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <button
                        key={star}
                        type="button"
                        onClick={() => setEditRating(star)}
                        className="text-2xl hover:scale-110 transition"
                      >
                        {star <= editRating ? (
                          <FaStar className="text-[#8A9A86]" />
                        ) : (
                          <FaRegStar className="text-gray-300" />
                        )}
                      </button>
                    ))}
                  </div>
                  <textarea
                    value={editComment}
                    onChange={(e) => setEditComment(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-sm resize-none"
                    rows={2}
                  />
                  <div className="flex gap-2 mt-2">
                    <button
                      onClick={() => handleEditReview(review.id)}
                      className="bg-[#8A9A86] text-white px-4 py-1.5 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium"
                    >
                      Сохранить
                    </button>
                    <button
                      onClick={cancelEdit}
                      className="bg-gray-200 text-gray-700 px-4 py-1.5 rounded-xl hover:bg-gray-300 transition text-sm font-medium"
                    >
                      Отмена
                    </button>
                  </div>
                </div>
              ) : (
                // Режим просмотра
                <>
                  <div className="flex items-center gap-2">
                    <FaUser className="text-gray-400 text-sm" />
                    <span className="text-sm font-medium text-[#1C1C1C]">
                      {review.user_name || 'Пользователь'}
                    </span>
                    <div className="flex items-center gap-0.5 ml-auto">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <span key={star} className="text-xs">
                          {star <= review.rating ? (
                            <FaStar className="text-[#8A9A86]" />
                          ) : (
                            <FaRegStar className="text-gray-300" />
                          )}
                        </span>
                      ))}
                    </div>
                  </div>
                  <p className="text-sm text-gray-600 mt-0.5">{review.comment}</p>
                  <p className="text-xs text-gray-400 mt-0.5">
                    {new Date(review.created_at).toLocaleDateString('ru-RU')}
                  </p>
                  {/* Кнопки редактирования и удаления */}
                  {review.user_id === user?.id && (
                    <div className="flex gap-3 mt-1">
                      <button
                        onClick={() => {
                          setEditingReview(review.id)
                          setEditRating(review.rating)
                          setEditComment(review.comment)
                        }}
                        className="text-xs text-blue-500 hover:text-blue-700 transition"
                      >
                        Редактировать
                      </button>
                      <button
                        onClick={() => handleDeleteReview(review.id)}
                        className="text-xs text-red-500 hover:text-red-700 transition"
                      >
                        Удалить
                      </button>
                    </div>
                  )}
                </>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default Reviews