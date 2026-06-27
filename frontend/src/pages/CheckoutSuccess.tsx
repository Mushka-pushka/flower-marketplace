import { Link } from 'react-router-dom'

const CheckoutSuccess = () => {
  return (
    <div className="max-w-2xl mx-auto text-center py-20 animate-fade-in-up">
      <div className="text-7xl mb-6">🎉</div>
      <h1 className="text-4xl font-bold gradient-text mb-4">
        Заказ успешно оформлен!
      </h1>
      <p className="text-lg text-gray-600 mb-8">
        Спасибо за ваш заказ. В ближайшее время с вами свяжется продавец для подтверждения.
      </p>
      <Link
        to="/profile"
        className="btn-primary px-8 py-3 rounded-full text-lg font-medium inline-block"
      >
        Перейти к моим заказам
      </Link>
    </div>
  )
}

export default CheckoutSuccess