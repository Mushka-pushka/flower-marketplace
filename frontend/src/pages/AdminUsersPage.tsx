import { useEffect, useState } from 'react'
import { FaUsers, FaUserCheck, FaUserTimes, FaSearch, FaEnvelope, FaCalendar } from 'react-icons/fa'
import { adminGetUsers, adminUpdateUserStatus, type UsersListResponse } from '../api/admin.api'

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  phone: string
  role: string
  is_active: boolean
  created_at: string
  updated_at: string
}

const AdminUsersPage = () => {
  const [allUsers, setAllUsers] = useState<User[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [filterRole, setFilterRole] = useState('')
  const [limit] = useState(10)
  const [offset, setOffset] = useState(0)

  // Загружаем пользователей при монтировании и при смене роли
  useEffect(() => {
    fetchUsers()
  }, [filterRole])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      // Загружаем ВСЕХ пользователей с большим лимитом
      const data: UsersListResponse = await adminGetUsers({
        role: filterRole || undefined,
        limit: 1000, // Загружаем всех
        offset: 0
      })
      
      console.log('Загружено пользователей:', data.users?.length)
      
      // Скрываем админа
      const filteredUsers = (data.users || []).filter(
        (user: User) => user.role !== 'admin'
      )
      
      setAllUsers(filteredUsers)
      setUsers(filteredUsers)
    } catch (error) {
      console.error('Ошибка загрузки пользователей:', error)
    } finally {
      setLoading(false)
    }
  }

  // Фильтрация на фронтенде (поиск)
  useEffect(() => {
    console.log('Поиск:', searchQuery)
    const filtered = allUsers.filter((user) => {
      if (searchQuery) {
        const query = searchQuery.toLowerCase().trim()
        const fullName = `${user.first_name || ''} ${user.last_name || ''}`.toLowerCase()
        const email = user.email.toLowerCase()
        
        // Проверяем совпадение
        const match = fullName.includes(query) || email.includes(query)
        if (match) {
          console.log('Найден:', user.email, user.first_name, user.last_name)
        }
        return match
      }
      return true
    })
    
    console.log('Найдено пользователей:', filtered.length)
    setUsers(filtered)
    setOffset(0)
  }, [searchQuery, allUsers])

  const toggleUserStatus = async (userId: string, currentStatus: boolean) => {
    if (!confirm(`Вы уверены, что хотите ${currentStatus ? 'заблокировать' : 'разблокировать'} этого пользователя?`)) return

    try {
      await adminUpdateUserStatus(userId, !currentStatus)
      setAllUsers(allUsers.map(u => 
        u.id === userId ? { ...u, is_active: !currentStatus } : u
      ))
    } catch (error) {
      console.error('Ошибка обновления статуса:', error)
    }
  }

  const getRoleLabel = (role: string) => {
    const map: Record<string, string> = {
      customer: 'Покупатель',
      seller: 'Продавец',
      admin: 'Администратор',
    }
    return map[role] || role
  }

  const getRoleColor = (role: string) => {
    const map: Record<string, string> = {
      customer: 'bg-blue-50 text-blue-600 border-blue-200',
      seller: 'bg-purple-50 text-purple-600 border-purple-200',
      admin: 'bg-red-50 text-red-600 border-red-200',
    }
    return map[role] || 'bg-gray-50 text-gray-600 border-gray-200'
  }

  const paginatedUsers = users.slice(offset, offset + limit)
  const totalPages = Math.ceil(users.length / limit)
  const currentPage = Math.floor(offset / limit) + 1

  const goToPreviousPage = () => {
    setOffset(Math.max(0, offset - limit))
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const goToNextPage = () => {
    setOffset(Math.min(offset + limit, users.length - limit))
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка пользователей...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaUsers className="text-[#8A9A86]" />
          Управление пользователями
          <span className="text-sm text-gray-400 font-normal ml-2">
            ({users.length} пользователей)
          </span>
        </h2>
      </div>

      {/* Фильтры */}
      <div className="flex flex-wrap items-center gap-3 mb-4">
        <div className="flex-1 min-w-[200px] relative">
          <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm" />
          <input
            type="text"
            placeholder="Поиск по email или имени..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          />
        </div>

        <select
          value={filterRole}
          onChange={(e) => {
            setFilterRole(e.target.value)
            setOffset(0)
          }}
          className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
        >
          <option value="">Все роли</option>
          <option value="customer">Покупатели</option>
          <option value="seller">Продавцы</option>
        </select>

        <span className="text-sm text-gray-400 ml-auto whitespace-nowrap">
          Найдено: {users.length}
        </span>
      </div>

      {/* Таблица */}
      <div className="overflow-x-auto bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] border border-gray-100">
        {paginatedUsers.length === 0 ? (
          <div className="text-center py-12 text-gray-400">
            <FaUsers className="text-4xl mx-auto mb-3 text-gray-300" />
            <p>Пользователи не найдены</p>
          </div>
        ) : (
          <>
            <table className="w-full text-sm">
              <thead className="bg-gray-50 border-b border-gray-100">
                <tr>
                  <th className="text-left p-3 font-medium text-[#1C1C1C]">Пользователь</th>
                  <th className="text-left p-3 font-medium text-[#1C1C1C]">Email</th>
                  <th className="text-left p-3 font-medium text-[#1C1C1C]">Роль</th>
                  <th className="text-left p-3 font-medium text-[#1C1C1C]">Дата регистрации</th>
                  <th className="text-center p-3 font-medium text-[#1C1C1C]">Статус</th>
                  <th className="text-center p-3 font-medium text-[#1C1C1C]">Действия</th>
                </tr>
              </thead>
              <tbody>
                {paginatedUsers.map((user) => (
                  <tr key={user.id} className="border-b border-gray-50 hover:bg-gray-50/50 transition">
                    <td className="p-3">
                      <span className="font-medium text-[#1C1C1C]">
                        {user.first_name || user.last_name 
                          ? `${user.first_name || ''} ${user.last_name || ''}`.trim()
                          : 'Без имени'}
                      </span>
                    </td>
                    <td className="p-3">
                      <div className="flex items-center gap-2">
                        <FaEnvelope className="text-gray-300 text-xs" />
                        <span className="text-gray-600">{user.email}</span>
                      </div>
                    </td>
                    <td className="p-3">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-medium border ${getRoleColor(user.role)}`}>
                        {getRoleLabel(user.role)}
                      </span>
                    </td>
                    <td className="p-3">
                      <div className="flex items-center gap-2">
                        <FaCalendar className="text-gray-300 text-xs" />
                        <span className="text-gray-500 text-xs">
                          {new Date(user.created_at).toLocaleDateString('ru-RU', {
                            day: '2-digit',
                            month: '2-digit',
                            year: 'numeric',
                          })}
                        </span>
                      </div>
                    </td>
                    <td className="p-3 text-center">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium border ${
                        user.is_active 
                          ? 'bg-green-50 text-green-600 border-green-200' 
                          : 'bg-red-50 text-red-600 border-red-200'
                      }`}>
                        {user.is_active ? 'Активен' : 'Заблокирован'}
                      </span>
                    </td>
                    <td className="p-3 text-center">
                      <button
                        onClick={() => toggleUserStatus(user.id, user.is_active)}
                        className={`px-3 py-1.5 rounded-lg text-xs font-medium transition flex items-center gap-1.5 mx-auto ${
                          user.is_active 
                            ? 'bg-red-50 text-red-600 hover:bg-red-100 border border-red-200' 
                            : 'bg-green-50 text-green-600 hover:bg-green-100 border border-green-200'
                        }`}
                      >
                        {user.is_active ? (
                          <>
                            <FaUserTimes className="text-xs" />
                            Заблокировать
                          </>
                        ) : (
                          <>
                            <FaUserCheck className="text-xs" />
                            Разблокировать
                          </>
                        )}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Пагинация */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center gap-2 p-4 border-t border-gray-100">
                <button
                  onClick={goToPreviousPage}
                  disabled={offset === 0}
                  className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition text-sm"
                >
                  Назад
                </button>
                <span className="px-4 py-2 text-[#1C1C1C] font-medium text-sm">
                  {currentPage} / {totalPages}
                </span>
                <button
                  onClick={goToNextPage}
                  disabled={offset + limit >= users.length}
                  className="px-4 py-2 border border-gray-200 rounded-xl hover:border-[#8A9A86] disabled:opacity-50 disabled:cursor-not-allowed transition text-sm"
                >
                  Вперед
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

export default AdminUsersPage