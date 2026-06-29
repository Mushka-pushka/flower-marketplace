import { Link } from 'react-router-dom'
import { FaCheckCircle } from 'react-icons/fa'

const CheckoutSuccess = () => {
  return (
    <div className="max-w-2xl mx-auto text-center py-20 animate-fade-in-up">
      <FaCheckCircle className="text-6xl text-[#8A9A86] mx-auto mb-6" />
      <h1 className="text-3xl font-bold text-[#1C1C1C] mb-4">
        Заказ успешно оформлен!
      </h1>
      <p className="text-base text-gray-400 mb-8 max-w-md mx-auto">
        Спасибо за ваш заказ. В ближайшее время с вами свяжется продавец для подтверждения.
      </p>
      <Link
        to="/profile"
        className="inline-block bg-[#8A9A86] text-white px-8 py-3 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium"
      >
        Перейти к моим заказам
      </Link>
    </div>
  )
}

export default CheckoutSuccess