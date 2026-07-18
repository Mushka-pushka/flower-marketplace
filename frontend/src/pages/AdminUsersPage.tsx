import { useEffect, useState } from 'react'
import { FaUsers, FaUserCheck, FaUserTimes, FaSearch } from 'react-icons/fa'
import { adminGetUsers, adminUpdateUserStatus } from '../api/admin.api'

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  role: string
  is_active: boolean
  created_at: string
}

const AdminUsersPage = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filterRole, setFilterRole] = useState('')

  useEffect(() => {
    fetchUsers()
  }, [search, filterRole])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const data = await adminGetUsers({ search, role: filterRole })
      setUsers(data)
    } catch (error) {
      console.error('Ошибка загрузки пользователей:', error)
    } finally {
      setLoading(false)
    }
  }

  const toggleUserStatus = async (userId: string, currentStatus: boolean) => {
    try {
      await adminUpdateUserStatus(userId, !currentStatus)
      setUsers(users.map(u => 
        u.id === userId ? { ...u, is_active: !currentStatus } : u
      ))
    } catch (error) {
      console.error('Ошибка обновления статуса:', error)
    }
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка...</div>
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-4 flex items-center gap-2">
        <FaUsers className="text-[#8A9A86]" />
        Управление пользователями
      </h2>

      {/* Фильтры */}
      <div className="flex flex-wrap gap-3 mb-4">
        <div className="flex-1 min-w-[200px]">
          <div className="relative">
            <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Поиск по email или имени..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86]"
            />
          </div>
        </div>
        <select
          value={filterRole}
          onChange={(e) => setFilterRole(e.target.value)}
          className="px-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] bg-white"
        >
          <option value="">Все роли</option>
          <option value="customer">Покупатель</option>
          <option value="seller">Продавец</option>
          <option value="admin">Админ</option>
        </select>
      </div>

      {/* Таблица */}
      <div className="overflow-x-auto bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] border border-gray-100">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-100">
            <tr>
              <th className="text-left p-3 font-medium text-[#1C1C1C]">Пользователь</th>
              <th className="text-left p-3 font-medium text-[#1C1C1C]">Email</th>
              <th className="text-left p-3 font-medium text-[#1C1C1C]">Роль</th>
              <th className="text-left p-3 font-medium text-[#1C1C1C]">Дата регистрации</th>
              <th className="text-left p-3 font-medium text-[#1C1C1C]">Статус</th>
              <th className="text-center p-3 font-medium text-[#1C1C1C]">Действия</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id} className="border-b border-gray-50 hover:bg-gray-50/50 transition">
                <td className="p-3">
                  {user.first_name} {user.last_name}
                </td>
                <td className="p-3 text-gray-600">{user.email}</td>
                <td className="p-3">
                  <span className={`px-2 py-0.5 rounded-full text-xs font-medium
                    ${user.role === 'admin' ? 'bg-purple-100 text-purple-700' : ''}
                    ${user.role === 'seller' ? 'bg-blue-100 text-blue-700' : ''}
                    ${user.role === 'customer' ? 'bg-gray-100 text-gray-700' : ''}
                  `}>
                    {user.role === 'customer' ? 'Покупатель' : 
                     user.role === 'seller' ? 'Продавец' : 'Админ'}
                  </span>
                </td>
                <td className="p-3 text-gray-600">
                  {new Date(user.created_at).toLocaleDateString('ru-RU')}
                </td>
                <td className="p-3">
                  <span className={`px-2 py-0.5 rounded-full text-xs font-medium
                    ${user.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}
                  `}>
                    {user.is_active ? 'Активен' : 'Заблокирован'}
                  </span>
                </td>
                <td className="p-3 text-center">
                  <button
                    onClick={() => toggleUserStatus(user.id, user.is_active)}
                    className={`px-3 py-1 rounded-lg text-xs font-medium transition flex items-center gap-1 mx-auto
                      ${user.is_active 
                        ? 'bg-red-50 text-red-600 hover:bg-red-100' 
                        : 'bg-green-50 text-green-600 hover:bg-green-100'
                      }`}
                  >
                    {user.is_active ? <FaUserTimes /> : <FaUserCheck />}
                    {user.is_active ? 'Заблокировать' : 'Разблокировать'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default AdminUsersPage