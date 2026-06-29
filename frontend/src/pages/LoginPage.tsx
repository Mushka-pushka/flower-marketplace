import { FaSignInAlt } from 'react-icons/fa'

const LoginPage = () => {
  return (
    <div className="max-w-md mx-auto mt-16 bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] p-8 border border-gray-100">
      <div className="text-center">
        <FaSignInAlt className="text-4xl text-[#8A9A86] mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-[#1C1C1C] mb-2">Вход</h2>
        <p className="text-gray-400 text-sm">Страница входа</p>
      </div>
    </div>
  )
}

export default LoginPage