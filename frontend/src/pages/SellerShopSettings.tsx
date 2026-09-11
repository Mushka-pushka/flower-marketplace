import { useState, useEffect } from 'react'
import {
  FaStore,
  FaSave,
  FaCheckCircle,
  FaPaperPlane,
  FaExclamationCircle,
} from 'react-icons/fa'
import { useAuth } from '../context/AuthContext'
import { updateShopName, getShopInfo } from '../api/admin.api'
import { toast } from 'react-hot-toast'
import client from '../api/client'

const SellerShopSettings = () => {
    const { user, updateUser } = useAuth()
    const [shopName, setShopName] = useState('')
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)
    const [requesting, setRequesting] = useState(false)
    const [shopInfo, setShopInfo] = useState<{
        id: string
        name: string
        is_verified: boolean
        rating: number
    } | null>(null)

    useEffect(() => {
        const fetchShopInfo = async () => {
            try {
                const data = await getShopInfo()
                setShopInfo(data)
                setShopName(data.name)
            } catch (error) {
                console.error('Ошибка загрузки информации о магазине:', error)
                toast.error('Не удалось загрузить данные магазина')
            } finally {
                setLoading(false)
            }
        }
        fetchShopInfo()
    }, [])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!shopName.trim()) {
            toast.error('Название магазина не может быть пустым')
            return
        }

        setSaving(true)
        try {
            const result = await updateShopName(shopName.trim())
            setShopInfo(prev => prev ? { ...prev, name: result.shop_name } : null)
            if (user) {
                updateUser({ shop_name: result.shop_name })
            }
            toast.success('Название магазина обновлено!')
        } catch (error) {
            console.error('Ошибка обновления названия:', error)
            toast.error('Не удалось обновить название магазина')
        } finally {
            setSaving(false)
        }
    }

    // Отправка магазина на верификацию
    const handleRequestVerification = async () => {
        if (!confirm('Отправить заявку на верификацию магазина? Администратор рассмотрит её в ближайшее время.')) {
            return
        }

        setRequesting(true)
        try {
            await client.post('/admin/shop/request-verification')
            toast.success('Заявка отправлена! Ожидайте проверки администратором.')
            // Обновляем данные магазина
            const data = await getShopInfo()
            setShopInfo(data)
        } catch (error: any) {
            console.error('Ошибка отправки заявки:', error)
            toast.error(error.response?.data?.error || 'Не удалось отправить заявку')
        } finally {
            setRequesting(false)
        }
    }

    if (loading) {
        return <div className="text-center py-8 text-gray-400">Загрузка...</div>
    }

    return (
        <div className="max-w-2xl mx-auto">
            <h2 className="text-2xl font-bold text-[#1C1C1C] mb-6 flex items-center gap-2">
                <FaStore className="text-[#8A9A86]" />
                Мой магазин
            </h2>

            {/* Статус верификации */}
            <div className={`rounded-xl p-5 mb-6 border ${
                shopInfo?.is_verified
                    ? 'bg-green-50 border-green-200'
                    : 'bg-amber-50 border-amber-200'
            }`}>
                <div className="flex items-start gap-3">
                    {shopInfo?.is_verified ? (
                        <FaCheckCircle className="text-green-600 text-2xl flex-shrink-0 mt-0.5" />
                    ) : (
                        <FaExclamationCircle className="text-amber-600 text-2xl flex-shrink-0 mt-0.5" />
                    )}
                    <div className="flex-1">
                        <h3 className={`font-semibold ${
                            shopInfo?.is_verified ? 'text-green-800' : 'text-amber-800'
                        }`}>
                            {shopInfo?.is_verified ? 'Магазин верифицирован' : 'Магазин не верифицирован'}
                        </h3>
                        <p className={`text-sm mt-1 ${
                            shopInfo?.is_verified ? 'text-green-700' : 'text-amber-700'
                        }`}>
                            {shopInfo?.is_verified
                                ? 'Ваши товары отображаются в публичном каталоге.'
                                : 'Ваши товары не отображаются в каталоге. Отправьте заявку на верификацию, чтобы покупатели могли их видеть.'
                            }
                        </p>

                        {!shopInfo?.is_verified && (
                            <button
                                onClick={handleRequestVerification}
                                disabled={requesting}
                                className="mt-3 px-4 py-2 bg-amber-600 text-white rounded-xl hover:bg-amber-700 transition text-sm font-medium flex items-center gap-2 disabled:opacity-50"
                            >
                                <FaPaperPlane />
                                {requesting ? 'Отправка...' : 'Отправить на верификацию'}
                            </button>
                        )}
                    </div>
                </div>
            </div>

            <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                            Название магазина
                        </label>
                        <input
                            type="text"
                            value={shopName}
                            onChange={(e) => setShopName(e.target.value)}
                            className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                            placeholder="Введите название магазина"
                            required
                        />
                        <p className="text-xs text-gray-400 mt-1">
                            Название будет отображаться в каталоге и в заказах
                        </p>
                    </div>

                    <button
                        type="submit"
                        disabled={saving}
                        className="flex items-center gap-2 px-6 py-2.5 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50"
                    >
                        <FaSave /> {saving ? 'Сохранение...' : 'Сохранить'}
                    </button>
                </form>

                {shopInfo && (
                    <div className="mt-6 pt-6 border-t border-gray-100">
                        <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                                <span className="text-gray-400">ID магазина:</span>
                                <p className="text-[#1C1C1C] font-medium">{shopInfo.id.slice(0, 8)}...</p>
                            </div>
                            <div>
                                <span className="text-gray-400">Рейтинг:</span>
                                <p className="text-[#1C1C1C] font-medium">{shopInfo.rating || 'Нет оценок'}</p>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
    )
}

export default SellerShopSettings