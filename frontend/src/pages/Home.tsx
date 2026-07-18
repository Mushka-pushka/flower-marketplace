import { Link } from 'react-router-dom'
import { FaLeaf, FaTruck, FaShieldAlt, FaSmile, FaArrowRight } from 'react-icons/fa'

const Home = () => {
  return (
    <div className="animate-fade-in-up">
      {/* Hero секция */}
      <div className="text-center py-16 md:py-24">
        <FaLeaf className="text-6xl text-[#8A9A86] mx-auto mb-6" />
        <h1 className="text-4xl md:text-6xl font-bold text-[#1C1C1C] mb-4">
          Цветы с доставкой по Беларуси
        </h1>
        <p className="text-lg text-gray-400 max-w-2xl mx-auto mb-8">
          Свежие букеты от лучших флористов. Создайте настроение себе и своим близким!
        </p>
        <Link
          to="/catalog"
          className="inline-flex items-center gap-2 bg-[#8A9A86] text-white px-8 py-4 rounded-xl hover:bg-[#7A8A76] transition text-base font-medium"
        >
          Смотреть каталог <FaArrowRight />
        </Link>
      </div>

      {/* Преимущества */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 py-12 border-t border-gray-100">
        <div className="text-center p-6">
          <FaTruck className="text-4xl text-[#8A9A86] mx-auto mb-3" />
          <h3 className="font-semibold text-[#1C1C1C] text-lg">Быстрая доставка</h3>
          <p className="text-gray-400 text-sm">Доставим цветы в день заказа</p>
        </div>
        <div className="text-center p-6">
          <FaShieldAlt className="text-4xl text-[#8A9A86] mx-auto mb-3" />
          <h3 className="font-semibold text-[#1C1C1C] text-lg">Гарантия свежести</h3>
          <p className="text-gray-400 text-sm">Только свежие цветы от проверенных продавцов</p>
        </div>
        <div className="text-center p-6">
          <FaSmile className="text-4xl text-[#8A9A86] mx-auto mb-3" />
          <h3 className="font-semibold text-[#1C1C1C] text-lg">Более 1000 букетов</h3>
          <p className="text-gray-400 text-sm">На любой вкус и повод</p>
        </div>
      </div>

      {/* Популярные категории (заглушка для будущего контента) */}
      <div className="py-12 border-t border-gray-100">
        <h2 className="text-2xl font-bold text-[#1C1C1C] text-center mb-6">
          Популярные категории
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {['Розы', 'Тюльпаны', 'Пионы', 'Хризантемы'].map((category) => (
            <Link
              key={category}
              to={`/catalog?category=${category.toLowerCase()}`}
              className="bg-gray-50 hover:bg-[#8A9A86]/10 rounded-xl p-6 text-center transition border border-gray-100"
            >
              <span className="text-[#1C1C1C] font-medium">{category}</span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  )
}

export default Home