import { Link } from 'react-router-dom'
import { FaLeaf } from 'react-icons/fa'

const Home = () => {
  return (
    <div className="text-center py-20 animate-fade-in-up">
      <FaLeaf className="text-6xl text-[#8A9A86] mx-auto mb-6" />
      <h1 className="text-4xl md:text-5xl font-bold text-[#1C1C1C] mb-4">
        Добро пожаловать!
      </h1>
      <p className="text-base text-gray-400 max-w-2xl mx-auto">
        Самые свежие цветы с доставкой по Беларуси.
        Создайте настроение себе и своим близким!
      </p>
      <div className="mt-8 flex justify-center gap-4">
        <Link
          to="/catalog"
          className="bg-[#8A9A86] text-white px-8 py-3 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium inline-block"
        >
          Смотреть каталог
        </Link>
      </div>
    </div>
  )
}

export default Home