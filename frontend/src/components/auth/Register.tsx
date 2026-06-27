import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
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
    <div className="max-w-md mx-auto mt-12 bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-pink-50/50 animate-fade-in-up">
      <h2 className="text-3xl font-bold gradient-text text-center mb-6">Регистрация</h2>
      {error && (
        <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">
          {error}
        </div>
      )}
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Имя</label>
          <input
            type="text"
            name="first_name"
            value={form.first_name}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Фамилия</label>
          <input
            type="text"
            name="last_name"
            value={form.last_name}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Email</label>
          <input
            type="email"
            name="email"
            value={form.email}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Пароль</label>
          <input
            type="password"
            name="password"
            value={form.password}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Телефон</label>
          <input
            type="text"
            name="phone"
            value={form.phone}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
          />
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-medium mb-1">Роль</label>
          <select
            name="role"
            value={form.role}
            onChange={handleChange}
            className="input-primary w-full px-4 py-2 rounded-lg"
          >
            <option value="customer">Покупатель</option>
            <option value="seller">Продавец</option>
          </select>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="btn-primary w-full py-3 rounded-full text-lg font-medium"
        >
          {loading ? 'Регистрация...' : 'Зарегистрироваться'}
        </button>
      </form>
      <p className="text-center text-sm text-gray-500 mt-4">
        Уже есть аккаунт? <Link to="/login" className="text-pink-500 hover:underline">Войти</Link>
      </p>
    </div>
  )
}

export default Register