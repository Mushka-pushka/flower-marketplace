import { useState, useEffect } from 'react'
import { FaStar, FaRegStar, FaUser, FaEdit, FaTrash } from 'react-icons/fa'
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
  const [hasUserReviewed, setHasUserReviewed] = useState(false)

  // Состояния для редактирования
  const [editingReview, setEditingReview] = useState<string | null>(null)
  const [editRating, setEditRating] = useState(5)
  const [editComment, setEditComment] = useState('')

  // Проверяем, есть ли у пользователя отзыв на этот товар
  useEffect(() => {
    if (user && reviews.length > 0) {
      const userReview = reviews.find(r => r.user_id === user.id)
      setHasUserReviewed(!!userReview)
    }
  }, [reviews, user])

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [reviewsData, canReviewData] = await Promise.all([
          getProductReviews(productId),
          user ? canReviewProduct(productId) : Promise.resolve(false),
        ])
        console.log('📥 Reviews loaded:', reviewsData)
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
    
    //Проверяем, не оставлял ли пользователь уже отзыв
    if (hasUserReviewed) {
      toast.error('Вы уже оставили отзыв на этот товар')
      return
    }

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
      toast.success('Спасибо за ваш отзыв!')
      // Обновляем список отзывов
      const data = await getProductReviews(productId)
      setReviews(data || [])
      setRating(5)
      setComment('')
      setHasUserReviewed(true)
    } catch (error: any) {
      console.error('Ошибка отправки отзыва:', error)
      const errorMsg = error.response?.data?.error || 'Не удалось отправить отзыв'
      toast.error(errorMsg)
    } finally {
      setSubmitting(false)
    }
  }

  // Удаление отзыва
  const handleDeleteReview = async (reviewId: string) => {
    if (!confirm('Удалить отзыв?')) return
    try {
      await deleteReview(reviewId)
      const data = await getProductReviews(productId)
      setReviews(data || [])
      setHasUserReviewed(false)
      toast.success('Отзыв удален')
    } catch (error) {
      console.error('Ошибка удаления отзыва:', error)
      toast.error('Не удалось удалить отзыв')
    }
  }

  // Начало редактирования отзыва
  const handleStartEdit = (review: Review) => {
    setEditingReview(review.id)
    setEditRating(review.rating)
    setEditComment(review.comment)
  }

  // Отмена редактирования
  const handleCancelEdit = () => {
    setEditingReview(null)
    setEditRating(5)
    setEditComment('')
  }

  // Сохранение отредактированного отзыва
  const handleSaveEdit = async (reviewId: string) => {
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

      {/* Форма отправки отзыва — только если пользователь может оставить отзыв и ещё не оставлял */}
      {user && canReview && !hasUserReviewed && (
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

      {/* Сообщение, если пользователь уже оставил отзыв */}
      {user && hasUserReviewed && (
        <div className="mb-4 p-3 bg-green-50 rounded-xl border border-green-200 text-sm text-green-700">
          ✓ Вы уже оставили отзыв на этот товар. Вы можете редактировать или удалить его ниже.
        </div>
      )}

      {user && !canReview && (
        <p className="text-sm text-gray-400 mb-3">
          Вы можете оставить отзыв только после получения заказа с этим товаром
        </p>
      )}

      <div className="space-y-3 max-h-48 overflow-y-auto">
        {reviews.length === 0 ? (
          <p className="text-gray-400 text-sm">Пока нет отзывов</p>
        ) : (
          reviews.map((review) => {
            const isOwnReview = review.user_id === user?.id
            const isEditing = editingReview === review.id

            return (
              <div key={review.id} className="border-b border-gray-100 pb-2 last:border-0">
                {isEditing ? (
                  // Режим редактирования
                  <div className="p-3 bg-gray-50 rounded-xl border border-gray-200">
                    <div className="flex items-center gap-1 mb-2">
                      {[1, 2, 3, 4, 5].map((star) => (
                        <button
                          key={star}
                          type="button"
                          onClick={() => setEditRating(star)}
                          className="text-xl hover:scale-110 transition"
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
                        onClick={() => handleSaveEdit(review.id)}
                        className="px-3 py-1 text-sm bg-[#8A9A86] text-white rounded-lg hover:bg-[#7A8A76] transition"
                      >
                        Сохранить
                      </button>
                      <button
                        onClick={handleCancelEdit}
                        className="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                      >
                        Отмена
                      </button>
                    </div>
                  </div>
                ) : (
                  // Обычный просмотр отзыва с аватаром
                  <>
                    <div className="flex items-center gap-2">
                      {review.user_avatar ? (
                        <img 
                          src={`http://localhost:8081${review.user_avatar}`} 
                          alt={review.user_name || 'Пользователь'} 
                          className="w-6 h-6 rounded-full object-cover border border-gray-200"
                          onError={(e) => {
                            // Если картинка не загрузилась — показываем иконку
                            e.currentTarget.style.display = 'none'
                          }}
                        />
                      ) : null}
                      {!review.user_avatar && (
                        <FaUser className="text-gray-400 text-sm" />
                      )}
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

                    {/* Кнопки редактирования/удаления для своих отзывов */}
                    {isOwnReview && (
                      <div className="flex gap-3 mt-1">
                        <button
                          onClick={() => handleStartEdit(review)}
                          className="text-xs text-blue-500 hover:text-blue-700 transition flex items-center gap-1"
                        >
                          <FaEdit className="text-[10px]" /> Редактировать
                        </button>
                        <button
                          onClick={() => handleDeleteReview(review.id)}
                          className="text-xs text-red-500 hover:text-red-700 transition flex items-center gap-1"
                        >
                          <FaTrash className="text-[10px]" /> Удалить
                        </button>
                      </div>
                    )}
                  </>
                )}
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}

export default Reviews