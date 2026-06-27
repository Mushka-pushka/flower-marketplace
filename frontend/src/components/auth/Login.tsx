import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login } from '../../api/auth.api'

const Login = () => {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await login(email, password)
      navigate('/')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto mt-12 bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg p-8 border border-pink-50/50 animate-fade-in-up">
      <h2 className="text-3xl font-bold gradient-text text-center mb-6">Вход</h2>
      {error && (
        <div className="bg-red-50 text-red-600 p-3 rounded-lg mb-4 text-sm">
          {error}
        </div>
      )}
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-gray-700 text-sm font-medium mb-1">Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-medium mb-1">Пароль</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="input-primary w-full px-4 py-2 rounded-lg"
            required
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="btn-primary w-full py-3 rounded-full text-lg font-medium"
        >
          {loading ? 'Вход...' : 'Войти'}
        </button>
      </form>
      <p className="text-center text-sm text-gray-500 mt-4">
        Нет аккаунта? <Link to="/register" className="text-pink-500 hover:underline">Зарегистрироваться</Link>
      </p>
    </div>
  )
}

export default Login