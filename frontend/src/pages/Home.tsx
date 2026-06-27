const Home = () => {
  return (
    <div className="text-center py-20 animate-fade-in-up">
      <h1 className="text-5xl md:text-6xl font-bold gradient-text mb-6">
        🌸 Добро пожаловать!
      </h1>
      <p className="text-xl text-gray-600 max-w-2xl mx-auto">
        Самые свежие цветы с доставкой по Беларуси. 
        Создайте настроение себе и своим близким!
      </p>
      <div className="mt-8 flex justify-center gap-4">
        <a href="/catalog" className="btn-primary px-8 py-3 rounded-full text-lg font-medium">
          Смотреть каталог
        </a>
      </div>
    </div>
  )
}

export default Home