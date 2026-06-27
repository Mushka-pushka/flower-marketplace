import { Link } from 'react-router-dom'

const CheckoutSuccess = () => {
  return (
    <div className="max-w-2xl mx-auto text-center py-16">
      <div className="text-6xl mb-6">🎉</div>
      <h1 className="text-3xl font-bold text-gray-800 mb-4">Заказ успешно оформлен!</h1>
      <p className="text-gray-600 mb-8">
        Спасибо за ваш заказ. В ближайшее время с вами свяжется продавец для подтверждения.
      </p>
      <Link
        to="/profile"
        className="bg-pink-500 text-white px-6 py-3 rounded-lg hover:bg-pink-600 transition"
      >
        Перейти к моим заказам
      </Link>
    </div>
  )
}

export default CheckoutSuccess