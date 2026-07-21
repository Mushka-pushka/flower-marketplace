import { useState, useRef } from 'react'
import { FaUser, FaCamera, FaSave, FaKey, FaEye, FaEyeSlash } from 'react-icons/fa'
import { useAuth } from '../context/AuthContext'
import { updateProfile, changePassword } from '../api/auth.api'
import { uploadAvatar } from '../api/upload.api'
import { toast } from 'react-hot-toast'

const ProfileSettings = () => {
  const { user, updateUser } = useAuth() 
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Личные данные
  const [firstName, setFirstName] = useState(user?.first_name || '')
  const [lastName, setLastName] = useState(user?.last_name || '')
  const [phone, setPhone] = useState(user?.phone || '')
  const [avatarUrl, setAvatarUrl] = useState(user?.avatar_url || '')
  const [loading, setLoading] = useState(false)

  console.log('Avatar URL from user:', user?.avatar_url)
  console.log('Avatar URL state:', avatarUrl)

  // Полный URL для аватара
  const avatarFullUrl = avatarUrl ? `http://localhost:8081${avatarUrl}` : ''

  // Смена пароля
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showOldPassword, setShowOldPassword] = useState(false)
  const [showNewPassword, setShowNewPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [passwordLoading, setPasswordLoading] = useState(false)

  // Сохранение личных данных
  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const updatedUser = await updateProfile({
        first_name: firstName,
        last_name: lastName,
        phone,
      })
      console.log('Profile updated:', updatedUser)
      // Обновляем данные в контексте через updateUser
      if (user) {
        updateUser({
          first_name: updatedUser.first_name,
          last_name: updatedUser.last_name,
          phone: updatedUser.phone,
        })
        console.log('User updated in context via updateUser')
      }
      toast.success('Профиль обновлён')
    } catch (error) {
      console.error('Ошибка обновления профиля:', error)
      toast.error('Не удалось обновить профиль')
    } finally {
      setLoading(false)
    }
  }

  // Смена пароля
  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (newPassword.length < 6) {
      toast.error('Новый пароль должен содержать минимум 6 символов')
      return
    }
    
    if (newPassword !== confirmPassword) {
      toast.error('Пароли не совпадают')
      return
    }

    setPasswordLoading(true)
    try {
      await changePassword(oldPassword, newPassword)
      toast.success('Пароль успешно изменён')
      setOldPassword('')
      setNewPassword('')
      setConfirmPassword('')
    } catch (error: any) {
      console.error('Ошибка смены пароля:', error)
      toast.error(error.response?.data?.error || 'Не удалось сменить пароль')
    } finally {
      setPasswordLoading(false)
    }
  }

  // Загрузка аватара
  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    console.log('File selected:', file.name, file.type, file.size)

    // Проверка типа и размера
    if (!file.type.startsWith('image/')) {
      toast.error('Файл должен быть изображением')
      return
    }
    if (file.size > 5 * 1024 * 1024) {
      toast.error('Файл слишком большой (макс. 5MB)')
      return
    }

    try {
      const result = await uploadAvatar(file)
      console.log('Upload result:', result)
      console.log('Avatar URL from response:', result.avatar_url)
      
      setAvatarUrl(result.avatar_url)
      console.log('Avatar URL state updated to:', result.avatar_url)
      
      // Обновляем пользователя в контексте через updateUser
      if (user) {
        updateUser({ avatar_url: result.avatar_url })
        console.log('User updated in context with avatar via updateUser')
      }
      toast.success('Аватар обновлён')
    } catch (error) {
      console.error('Ошибка загрузки аватара:', error)
      toast.error('Не удалось загрузить аватар')
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-6 flex items-center gap-2">
        <FaUser className="text-[#8A9A86]" />
        Настройки профиля
      </h2>

      <div className="space-y-6">
        {/* Личные данные */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4">Личные данные</h3>
          
          <form onSubmit={handleSaveProfile} className="space-y-4">
            {/* Аватар */}
            <div className="flex items-center gap-4">
              <div className="relative w-20 h-20 rounded-full bg-gray-100 flex items-center justify-center overflow-hidden border-2 border-gray-200">
                {avatarFullUrl ? (
                  <img 
                    src={avatarFullUrl} 
                    alt="Avatar" 
                    className="w-full h-full object-cover"
                    key={avatarFullUrl}
                  />
                ) : (
                  <FaUser className="text-3xl text-gray-400" />
                )}
              </div>
              <div>
                <button
                  type="button"
                  onClick={() => fileInputRef.current?.click()}
                  className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-[#1C1C1C] rounded-xl hover:bg-gray-200 transition text-sm font-medium"
                >
                  <FaCamera /> Загрузить фото
                </button>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  onChange={handleAvatarUpload}
                  className="hidden"
                />
                <p className="text-xs text-gray-400 mt-1">PNG, JPG до 5MB</p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Имя</label>
                <input
                  type="text"
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Фамилия</label>
                <input
                  type="text"
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Email</label>
              <input
                type="email"
                value={user?.email || ''}
                disabled
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl bg-gray-50 text-gray-500 cursor-not-allowed"
              />
              <p className="text-xs text-gray-400 mt-1">Email нельзя изменить</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Телефон</label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="+375 (29) 123-45-67"
                className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="flex items-center gap-2 px-6 py-2.5 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50"
            >
              <FaSave /> {loading ? 'Сохранение...' : 'Сохранить изменения'}
            </button>
          </form>
        </div>

        {/* Смена пароля */}
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-6 border border-gray-100">
          <h3 className="text-lg font-semibold text-[#1C1C1C] mb-4 flex items-center gap-2">
            <FaKey className="text-[#8A9A86]" />
            Безопасность
          </h3>

          <form onSubmit={handleChangePassword} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Старый пароль</label>
              <div className="relative">
                <input
                  type={showOldPassword ? 'text' : 'password'}
                  value={oldPassword}
                  onChange={(e) => setOldPassword(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition pr-12"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowOldPassword(!showOldPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  {showOldPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                Новый пароль <span className="text-gray-400 text-xs">(минимум 6 символов)</span>
              </label>
              <div className="relative">
                <input
                  type={showNewPassword ? 'text' : 'password'}
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition pr-12"
                  required
                  minLength={6}
                />
                <button
                  type="button"
                  onClick={() => setShowNewPassword(!showNewPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  {showNewPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Подтверждение пароля</label>
              <div className="relative">
                <input
                  type={showConfirmPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition pr-12"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  {showConfirmPassword ? <FaEyeSlash /> : <FaEye />}
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={passwordLoading}
              className="flex items-center gap-2 px-6 py-2.5 bg-gray-800 text-white rounded-xl hover:bg-gray-700 transition text-sm font-medium disabled:opacity-50"
            >
              <FaKey /> {passwordLoading ? 'Смена...' : 'Сменить пароль'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}

export default ProfileSettings