import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { FaSignInAlt } from 'react-icons/fa'
import { login } from '../../api/auth.api'
import { useAuth } from '../../context/AuthContext'

const Login = () => {
  const navigate = useNavigate()
  const { login: authLogin } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await login(email, password)
      authLogin(response.user, response.access_token)
      navigate('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-16 bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-8 animate-fade-in-up border border-gray-100">
      <h2 className="text-3xl font-bold text-[#1C1C1C] text-center mb-2">Вход</h2>
      <p className="text-center text-gray-400 text-sm mb-6">Войдите в свой аккаунт</p>

      {error && (
        <div className="bg-red-50 text-red-500 p-3 rounded-lg mb-4 text-sm">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <div className="mb-6">
          <label className="block text-[#1C1C1C] text-sm font-medium mb-1">Пароль</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
            required
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full bg-[#8A9A86] text-white py-3 rounded-xl hover:bg-[#7A8A76] transition flex items-center justify-center gap-2 text-base font-medium"
        >
          <FaSignInAlt />
          {loading ? 'Вход...' : 'Войти'}
        </button>
      </form>

      <p className="text-center text-sm text-gray-400 mt-4">
        Нет аккаунта?{' '}
        <Link to="/register" className="text-[#8A9A86] hover:underline font-medium">
          Зарегистрироваться
        </Link>
      </p>
    </div>
  )
}

export default Login