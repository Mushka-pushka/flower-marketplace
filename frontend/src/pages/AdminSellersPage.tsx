import { useEffect, useState } from 'react'
import {
  FaStore,
  FaCheckCircle,
  FaTimesCircle,
  FaUserCheck,
  FaUserTimes,
  FaSearch,
  FaSpinner,
  FaStar,
  FaRegStar,
} from 'react-icons/fa'
import { adminGetSellers, adminVerifySeller, type SellerWithShop } from '../api/admin.api'
import { toast } from 'react-hot-toast'

const AdminSellersPage = () => {
  const [sellers, setSellers] = useState<SellerWithShop[]>([])
  const [loading, setLoading] = useState(true)
  const [filterVerified, setFilterVerified] = useState<'all' | 'verified' | 'unverified'>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [actionLoading, setActionLoading] = useState<string | null>(null)

  useEffect(() => {
    fetchSellers()
  }, [filterVerified])

  const fetchSellers = async () => {
    try {
      setLoading(true)
      const params: { verified?: boolean } = {}
      if (filterVerified === 'verified') params.verified = true
      if (filterVerified === 'unverified') params.verified = false
      
      const data = await adminGetSellers(params)
      setSellers(Array.isArray(data) ? data : [])
    } catch (error) {
      console.error('Ошибка загрузки продавцов:', error)
      toast.error('Не удалось загрузить список продавцов')
    } finally {
      setLoading(false)
    }
  }

  const handleVerify = async (shopId: string, verify: boolean) => {
    const action = verify ? 'верифицировать' : 'отклонить'
    if (!confirm(`Вы уверены, что хотите ${action} этого продавца?`)) return

    setActionLoading(shopId)
    try {
      await adminVerifySeller(shopId, verify)
      toast.success(`Продавец ${verify ? 'верифицирован' : 'отклонён'}`)
      await fetchSellers()
    } catch (error) {
      console.error('Ошибка верификации:', error)
      toast.error('Не удалось обновить статус продавца')
    } finally {
      setActionLoading(null)
    }
  }

  const getStatusBadge = (isVerified: boolean) => {
    if (isVerified) {
      return (
        <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-green-50 text-green-700 border border-green-200">
          <FaCheckCircle className="text-xs" />
          Верифицирован
        </span>
      )
    }
    return (
      <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-yellow-50 text-yellow-700 border border-yellow-200">
        <FaTimesCircle className="text-xs" />
        Не верифицирован
      </span>
    )
  }

  const getActiveBadge = (isActive: boolean) => {
    if (isActive) {
      return (
        <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-green-50 text-green-700">
          <FaUserCheck className="text-xs" />
          Активен
        </span>
      )
    }
    return (
      <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-red-50 text-red-700">
        <FaUserTimes className="text-xs" />
        Заблокирован
      </span>
    )
  }

  const renderStars = (rating: number) => {
    const fullStars = Math.floor(rating)
    const hasHalfStar = rating % 1 >= 0.5
    const emptyStars = 5 - fullStars - (hasHalfStar ? 1 : 0)

    return (
      <div className="flex items-center gap-0.5">
        {[...Array(fullStars)].map((_, i) => (
          <FaStar key={`full-${i}`} className="text-[#8A9A86] text-sm" />
        ))}
        {hasHalfStar && <FaStar className="text-[#8A9A86] text-sm opacity-50" />}
        {[...Array(emptyStars)].map((_, i) => (
          <FaRegStar key={`empty-${i}`} className="text-gray-300 text-sm" />
        ))}
        <span className="text-xs text-gray-400 ml-1">({rating?.toFixed(1) || '0'})</span>
      </div>
    )
  }

  const filteredSellers = sellers.filter((seller) => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return (
      seller.shop_name?.toLowerCase().includes(query) ||
      seller.email?.toLowerCase().includes(query) ||
      `${seller.first_name} ${seller.last_name}`.toLowerCase().includes(query)
    )
  })

  if (loading) {
    return (
      <div className="text-center py-12">
        <FaSpinner className="text-3xl text-[#8A9A86] animate-spin mx-auto" />
        <p className="text-gray-400 mt-3">Загрузка продавцов...</p>
      </div>
    )
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaStore className="text-[#8A9A86]" />
        Модерация продавцов
      </h2>
      <p className="text-gray-400 text-base mb-4">
        Подтверждение и отклонение заявок на открытие магазина
      </p>

      {/* Фильтры */}
      <div className="flex flex-wrap items-center gap-3 mb-4">
        <div className="flex-1 min-w-[200px] relative">
          <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm" />
          <input
            type="text"
            placeholder="Поиск по названию магазина или email..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          />
        </div>

        <select
          value={filterVerified}
          onChange={(e) => setFilterVerified(e.target.value as 'all' | 'verified' | 'unverified')}
          className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
        >
          <option value="all">Все продавцы</option>
          <option value="verified">Верифицированные</option>
          <option value="unverified">Не верифицированные</option>
        </select>

        <span className="text-sm text-gray-400 ml-auto whitespace-nowrap">
          Найдено: {filteredSellers.length} из {sellers.length}
        </span>
      </div>

      {/* Список продавцов */}
      {filteredSellers.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border border-gray-100">
          <p className="text-gray-400 text-lg">Продавцы не найдены</p>
          <p className="text-gray-400 text-sm mt-1">
            {filterVerified !== 'all' 
              ? 'Попробуйте изменить фильтр' 
              : 'На платформе пока нет продавцов'}
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredSellers.map((seller) => (
            <div
              key={seller.shop_id}
              className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-5 border border-gray-100 hover:shadow-[0_8px_30px_rgba(0,0,0,0.06)] transition-all"
            >
              <div className="flex flex-col md:flex-row md:items-start gap-4">
                {/* Левая часть - информация */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-3 flex-wrap">
                    <h3 className="text-lg font-semibold text-[#1C1C1C]">
                      {seller.shop_name || 'Магазин без названия'}
                    </h3>
                    {getStatusBadge(seller.is_verified)}
                    {getActiveBadge(seller.is_active)}
                  </div>

                  {seller.shop_description && (
                    <p className="text-sm text-gray-500 mt-1">{seller.shop_description}</p>
                  )}

                  <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm">
                    <div className="flex items-center gap-1 text-gray-400">
                      <FaStore className="text-[#8A9A86]" />
                      <span className="text-[#1C1C1C] font-medium">
                        {seller.first_name} {seller.last_name}
                      </span>
                    </div>
                    <div className="text-gray-400">
                      Email: <span className="text-[#1C1C1C]">{seller.email}</span>
                    </div>
                    {seller.phone && (
                      <div className="text-gray-400">
                        Телефон: <span className="text-[#1C1C1C]">{seller.phone}</span>
                      </div>
                    )}
                    <div className="text-gray-400">
                      Дата регистрации:{' '}
                      <span className="text-[#1C1C1C]">
                        {new Date(seller.created_at).toLocaleDateString('ru-RU', {
                          day: '2-digit',
                          month: '2-digit',
                          year: 'numeric',
                        })}
                      </span>
                    </div>
                    <div className="text-gray-400">
                      ID магазина:{' '}
                      <span className="text-[#1C1C1C] font-mono text-xs">
                        {seller.shop_id.slice(0, 8)}
                      </span>
                    </div>
                  </div>

                  <div className="mt-2">{renderStars(seller.rating || 0)}</div>
                </div>

                {/* Правая часть - действия */}
                <div className="flex flex-col gap-2 min-w-[140px]">
                  {!seller.is_verified ? (
                    <>
                      <button
                        onClick={() => handleVerify(seller.shop_id, true)}
                        disabled={actionLoading === seller.shop_id}
                        className="px-4 py-2 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {actionLoading === seller.shop_id ? (
                          <FaSpinner className="animate-spin" />
                        ) : (
                          <FaCheckCircle />
                        )}
                        Верифицировать
                      </button>
                      <button
                        onClick={() => handleVerify(seller.shop_id, false)}
                        disabled={actionLoading === seller.shop_id}
                        className="px-4 py-2 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition text-sm font-medium flex items-center justify-center gap-2 border border-red-200 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <FaTimesCircle />
                        Отклонить
                      </button>
                    </>
                  ) : (
                    <button
                      onClick={() => handleVerify(seller.shop_id, false)}
                      disabled={actionLoading === seller.shop_id}
                      className="px-4 py-2 bg-red-50 text-red-600 rounded-xl hover:bg-red-100 transition text-sm font-medium flex items-center justify-center gap-2 border border-red-200 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {actionLoading === seller.shop_id ? (
                        <FaSpinner className="animate-spin" />
                      ) : (
                        <FaTimesCircle />
                      )}
                      Снять верификацию
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default AdminSellersPage