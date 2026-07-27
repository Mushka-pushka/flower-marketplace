import { useState, useEffect } from 'react'
import {
  FaStar,
  FaRegStar,
  FaReply,
  FaEdit,
  FaTrash,
  FaSpinner,
} from 'react-icons/fa'
import { getSellerReviews, addReplyToReview, updateReplyOnReview, deleteReplyFromReview, deleteReviewBySeller } from '../api/catalog.api'
import type { ReviewWithReply } from '../api/catalog.api'
import { toast } from 'react-hot-toast'
import { useAuth } from '../context/AuthContext'

const SellerReviewsPage = () => {
  const { user } = useAuth()
  const [reviews, setReviews] = useState<ReviewWithReply[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [limit] = useState(20)
  const [offset, setOffset] = useState(0)
  const [replyingTo, setReplyingTo] = useState<string | null>(null)
  const [replyText, setReplyText] = useState('')
  const [editingReply, setEditingReply] = useState<string | null>(null)
  const [editReplyText, setEditReplyText] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [expandedReviews, setExpandedReviews] = useState<Set<string>>(new Set())

  const fetchReviews = async () => {
    if (!user?.shop_id) {
      setLoading(false)
      return
    }

    try {
      setLoading(true)
      const data = await getSellerReviews({ limit, offset })
      setReviews(data.reviews || [])
      setTotal(data.total || 0)
    } catch (error) {
      console.error('Ошибка загрузки отзывов:', error)
      toast.error('Не удалось загрузить отзывы')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchReviews()
  }, [user, limit, offset])

  const handleAddReply = async (reviewId: string) => {
    if (!replyText.trim()) {
      toast.error('Введите текст ответа')
      return
    }

    setSubmitting(true)
    try {
      await addReplyToReview(reviewId, replyText.trim())
      toast.success('Ответ добавлен')
      setReplyingTo(null)
      setReplyText('')
      await fetchReviews()
    } catch (error: any) {
      console.error('Ошибка добавления ответа:', error)
      toast.error(error.response?.data?.error || 'Не удалось добавить ответ')
    } finally {
      setSubmitting(false)
    }
  }

  const handleUpdateReply = async (reviewId: string) => {
    if (!editReplyText.trim()) {
      toast.error('Введите текст ответа')
      return
    }

    setSubmitting(true)
    try {
      await updateReplyOnReview(reviewId, editReplyText.trim())
      toast.success('Ответ обновлён')
      setEditingReply(null)
      setEditReplyText('')
      await fetchReviews()
    } catch (error: any) {
      console.error('Ошибка обновления ответа:', error)
      toast.error(error.response?.data?.error || 'Не удалось обновить ответ')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteReply = async (reviewId: string) => {
    if (!confirm('Удалить ответ на отзыв?')) return

    try {
      await deleteReplyFromReview(reviewId)
      toast.success('Ответ удалён')
      await fetchReviews()
    } catch (error: any) {
      console.error('Ошибка удаления ответа:', error)
      toast.error(error.response?.data?.error || 'Не удалось удалить ответ')
    }
  }

  // Удаление отзыва
  const handleDeleteReview = async (reviewId: string) => {
    if (!confirm('Вы уверены, что хотите удалить этот отзыв? Это действие нельзя отменить.')) return

    try {
      await deleteReviewBySeller(reviewId)
      toast.success('Отзыв удалён')
      await fetchReviews()
    } catch (error: any) {
      console.error('Ошибка удаления отзыва:', error)
      toast.error(error.response?.data?.error || 'Не удалось удалить отзыв')
    }
  }

  const toggleExpand = (reviewId: string) => {
    setExpandedReviews(prev => {
      const newSet = new Set(prev)
      if (newSet.has(reviewId)) {
        newSet.delete(reviewId)
      } else {
        newSet.add(reviewId)
      }
      return newSet
    })
  }

  const renderStars = (rating: number) => {
    return (
      <div className="flex items-center gap-0.5">
        {[1, 2, 3, 4, 5].map((star) => (
          <span key={star}>
            {star <= rating ? (
              <FaStar className="text-[#8A9A86] text-sm" />
            ) : (
              <FaRegStar className="text-gray-300 text-sm" />
            )}
          </span>
        ))}
      </div>
    )
  }

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  if (loading) {
    return (
      <div className="text-center py-8">
        <FaSpinner className="text-2xl text-[#8A9A86] animate-spin mx-auto" />
        <p className="text-gray-400 mt-2">Загрузка отзывов...</p>
      </div>
    )
  }

  if (reviews.length === 0) {
    return (
      <div>
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
          <FaStar className="text-[#8A9A86]" />
          Отзывы покупателей
        </h2>
        <p className="text-gray-400 text-center py-8">
          На ваши товары пока нет отзывов
        </p>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaStar className="text-[#8A9A86]" />
          Отзывы покупателей
          <span className="text-sm text-gray-400 font-normal ml-2">
            ({total} отзывов)
          </span>
        </h2>
      </div>

      {/* Статистика */}
      {reviews.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
            <div className="flex items-center gap-2">
              <FaStar className="text-[#8A9A86] text-lg" />
              <p className="text-2xl font-bold text-[#1C1C1C]">
                {(reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length).toFixed(1)}
              </p>
            </div>
            <p className="text-sm text-gray-400">Средний рейтинг</p>
          </div>
          
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
            <div className="flex items-center gap-2">
              <FaRegStar className="text-[#8A9A86] text-lg" />
              <p className="text-2xl font-bold text-[#1C1C1C]">{reviews.length}</p>
            </div>
            <p className="text-sm text-gray-400">Всего отзывов</p>
          </div>
          
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
            <div className="flex items-center gap-2">
              <FaStar className="text-green-500 text-lg" />
              <p className="text-2xl font-bold text-[#1C1C1C]">
                {reviews.filter(r => r.rating >= 4).length}
              </p>
            </div>
            <p className="text-sm text-gray-400">Положительных (4-5★)</p>
          </div>
          
          <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-4 border border-gray-100">
            <div className="flex items-center gap-2">
              <FaReply className="text-[#8A9A86] text-lg" />
              <p className="text-2xl font-bold text-[#1C1C1C]">
                {reviews.filter(r => r.reply).length}
              </p>
            </div>
            <p className="text-sm text-gray-400">С ответами</p>
          </div>
        </div>
      )}

      {/* Список отзывов */}
      <div className="space-y-4">
        {reviews.map((review) => {
          const isExpanded = expandedReviews.has(review.id)
          const hasReply = !!review.reply
          const isReplying = replyingTo === review.id
          const isEditing = editingReply === review.id

          return (
            <div
              key={review.id}
              className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition-all"
            >
              <div className="flex items-start gap-4">
                {/* Аватар */}
                <div className="w-10 h-10 rounded-full bg-[#8A9A86]/10 flex items-center justify-center flex-shrink-0">
                  {review.user_avatar ? (
                    <img
                      src={`http://localhost:8081${review.user_avatar}`}
                      alt={review.user_name}
                      className="w-full h-full rounded-full object-cover"
                    />
                  ) : (
                    <span className="text-[#8A9A86] font-medium text-sm">
                      {review.user_name ? getInitials(review.user_name) : 'U'}
                    </span>
                  )}
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 flex-wrap">
                    <span className="font-medium text-[#1C1C1C]">
                      {review.user_name || 'Пользователь'}
                    </span>
                    {renderStars(review.rating)}
                    <span className="text-xs text-gray-400">
                      {new Date(review.created_at).toLocaleDateString('ru-RU', {
                        day: '2-digit',
                        month: '2-digit',
                        year: 'numeric',
                      })}
                    </span>
                  </div>

                  <p className="text-gray-600 mt-1">{review.comment}</p>

                  <button
                    onClick={() => toggleExpand(review.id)}
                    className="text-xs text-[#8A9A86] hover:underline mt-1 font-medium"
                  >
                    {isExpanded ? 'Скрыть' : 'Подробнее'}
                  </button>

                  {isExpanded && (
                    <div className="mt-2">
                      <p className="text-sm text-gray-400">
                        Товар: <span className="text-[#1C1C1C] font-medium">{review.product_name}</span>
                      </p>
                      <p className="text-sm text-gray-400">
                        ID отзыва: <span className="font-mono">{review.id.slice(0, 8)}</span>
                      </p>
                    </div>
                  )}

                  {/* Блок с ответом продавца */}
                  {hasReply && !isEditing && (
                    <div className="mt-3 pl-4 border-l-2 border-[#8A9A86] bg-gray-50/50 rounded-r-lg p-3">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className="text-xs font-medium text-[#8A9A86]">Ваш ответ:</span>
                          <span className="text-xs text-gray-400">
                            {review.reply_created_at && new Date(review.reply_created_at).toLocaleDateString('ru-RU')}
                          </span>
                        </div>
                        <div className="flex gap-2">
                          <button
                            onClick={() => {
                              setEditingReply(review.id)
                              setEditReplyText(review.reply || '')
                            }}
                            className="text-blue-500 hover:text-blue-700 text-xs transition"
                          >
                            <FaEdit />
                          </button>
                          <button
                            onClick={() => handleDeleteReply(review.id)}
                            className="text-red-500 hover:text-red-700 text-xs transition"
                          >
                            <FaTrash />
                          </button>
                        </div>
                      </div>
                      <p className="text-sm text-gray-700 mt-1">{review.reply}</p>
                      {review.reply_updated_at && review.reply_created_at !== review.reply_updated_at && (
                        <span className="text-[10px] text-gray-400">(отредактировано)</span>
                      )}
                    </div>
                  )}

                  {/* Форма редактирования ответа */}
                  {isEditing && (
                    <div className="mt-3 pl-4 border-l-2 border-[#8A9A86] bg-gray-50 rounded-r-lg p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs font-medium text-[#8A9A86]">Редактировать ответ:</span>
                      </div>
                      <textarea
                        value={editReplyText}
                        onChange={(e) => setEditReplyText(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-sm resize-none"
                        rows={2}
                        placeholder="Введите ответ..."
                      />
                      <div className="flex gap-2 mt-2">
                        <button
                          onClick={() => handleUpdateReply(review.id)}
                          disabled={submitting || !editReplyText.trim()}
                          className="px-3 py-1 text-sm bg-[#8A9A86] text-white rounded-lg hover:bg-[#7A8A76] transition disabled:opacity-50 flex items-center gap-1"
                        >
                          {submitting ? <FaSpinner className="animate-spin" /> : <FaEdit />}
                          Сохранить
                        </button>
                        <button
                          onClick={() => {
                            setEditingReply(null)
                            setEditReplyText('')
                          }}
                          className="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                        >
                          Отмена
                        </button>
                      </div>
                    </div>
                  )}

                  {/* Кнопки действий продавца */}
                  <div className="mt-3 flex flex-wrap items-center gap-3">
                    {/* Кнопка "Ответить" */}
                    {!hasReply && !isReplying && (
                      <button
                        onClick={() => {
                          setReplyingTo(review.id)
                          setReplyText('')
                        }}
                        className="text-sm text-[#8A9A86] hover:text-[#7A8A76] transition flex items-center gap-1.5 font-medium"
                      >
                        <FaReply /> Ответить на отзыв
                      </button>
                    )}

                    {/* Кнопка "Удалить отзыв" */}
                    <button
                      onClick={() => handleDeleteReview(review.id)}
                      className="text-sm text-red-500 hover:text-red-700 transition flex items-center gap-1.5 font-medium"
                    >
                      <FaTrash /> Удалить отзыв
                    </button>
                  </div>

                  {/* Форма добавления ответа */}
                  {isReplying && (
                    <div className="mt-3 pl-4 border-l-2 border-[#8A9A86] bg-gray-50 rounded-r-lg p-3">
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs font-medium text-[#8A9A86]">Ваш ответ:</span>
                      </div>
                      <textarea
                        value={replyText}
                        onChange={(e) => setReplyText(e.target.value)}
                        className="w-full px-3 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-sm resize-none"
                        rows={2}
                        placeholder="Введите ответ на отзыв..."
                        autoFocus
                      />
                      <div className="flex gap-2 mt-2">
                        <button
                          onClick={() => handleAddReply(review.id)}
                          disabled={submitting || !replyText.trim()}
                          className="px-3 py-1 text-sm bg-[#8A9A86] text-white rounded-lg hover:bg-[#7A8A76] transition disabled:opacity-50 flex items-center gap-1"
                        >
                          {submitting ? <FaSpinner className="animate-spin" /> : <FaReply />}
                          Отправить
                        </button>
                        <button
                          onClick={() => {
                            setReplyingTo(null)
                            setReplyText('')
                          }}
                          className="px-3 py-1 text-sm bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                        >
                          Отмена
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {/* Пагинация */}
      {total > limit && (
        <div className="flex justify-center items-center gap-2 mt-6">
          <button
            onClick={() => setOffset(Math.max(0, offset - limit))}
            disabled={offset === 0}
            className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition text-sm"
          >
            Назад
          </button>
          <span className="px-4 py-2 text-[#1C1C1C] font-medium text-sm">
            {Math.floor(offset / limit) + 1} / {Math.ceil(total / limit)}
          </span>
          <button
            onClick={() => setOffset(offset + limit)}
            disabled={offset + limit >= total}
            className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition text-sm"
          >
            Вперед
          </button>
        </div>
      )}
    </div>
  )
}

export default SellerReviewsPage