import { FaUsers } from 'react-icons/fa'

const AdminUsersPage = () => {
  return (
    <div>
      <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2 flex items-center gap-2">
        <FaUsers className="text-[#8A9A86]" />
        Управление пользователями
      </h2>
      <p className="text-gray-400 text-base">Список всех пользователей, блокировка</p>
    </div>
  )
}

export default AdminUsersPage