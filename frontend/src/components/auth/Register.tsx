import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { FaUserPlus } from 'react-icons/fa'
import { register } from '../../api/auth.api'

const Register = () => {
  const navigate = useNavigate()
  const [form, setForm] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    phone: '',
    role: 'customer' as 'customer' | 'seller',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await register(form)
      navigate('/login')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-16 bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-8 animate-fade-in-up border border-gray-100">
      <h2 className="text-3xl font-bold text-[#1C1C1C] text-center mb-2">Регистрация</h2>
      <p className="text-center text-gray-400 text-sm mb-6">Создайте новый аккаунт</p>

      {error && (
        <div className="bg-red-50 text-red-500 p-3 rounded-lg mb-4 text-sm">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Имя</label>
          <input
            type="text"
            name="first_name"
            value={form.first_name}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Фамилия</label>
          <input
            type="text"
            name="last_name"
            value={form.last_name}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Email</label>
          <input
            type="email"
            name="email"
            value={form.email}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Пароль</label>
          <input
            type="password"
            name="password"
            value={form.password}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Телефон</label>
          <input
            type="text"
            name="phone"
            value={form.phone}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
          />
        </div>
        <div className="mb-6">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Роль</label>
          <select
            name="role"
            value={form.role}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white"
          >
            <option value="customer">Покупатель</option>
            <option value="seller">Продавец</option>
          </select>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-[#8A9A86] text-white py-3 rounded-xl hover:bg-[#7A8A76] transition flex items-center justify-center gap-2 text-base font-medium"
        >
          <FaUserPlus />
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
      </form>

      <p className="text-center text-sm text-gray-400 mt-4">
        Уже есть аккаунт?{' '}
        <Link to="/login" className="text-[#8A9A86] hover:underline font-medium">
          Войти
        </Link>
      </p>
    </div>
  )
}

export default Register