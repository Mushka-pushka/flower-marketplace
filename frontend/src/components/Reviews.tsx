import { useState, useEffect } from 'react'
import { FaStar, FaRegStar, FaUser } from 'react-icons/fa'
import { getProductReviews, createReview } from '../api/catalog.api'
import { canReviewProduct } from '../api/order.api'
import type { Review } from '../api/catalog.api'
import { useAuth } from '../context/AuthContext'

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
      alert('Войдите в аккаунт, чтобы оставить отзыв')
      return
    }

    if (!canReview) {
      alert('Вы можете оставить отзыв только на товары, которые заказывали и получили')
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
      alert('Спасибо за ваш отзыв!')
    } catch (error) {
      console.error('Ошибка отправки отзыва:', error)
      alert('Не удалось отправить отзыв')
    } finally {
      setSubmitting(false)
    }
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
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default Reviews