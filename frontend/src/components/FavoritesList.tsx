import { useFavorites } from '../context/FavoritesContext'

const FavoritesList = () => {
  const { items, removeFavorite } = useFavorites()

  if (items.length === 0) {
    return <p className="text-gray-500">У вас пока нет избранных товаров</p>
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {items.map((item) => (
        <div key={item.id} className="border rounded-lg p-4 flex justify-between items-center">
          <div>
            <p className="font-medium">{item.name}</p>
            <p className="text-pink-600 font-bold">{item.price} BYN</p>
          </div>
          <button
            onClick={() => removeFavorite(item.product_id)}
            className="text-red-400 hover:text-red-600"
          >
            ✕
          </button>
        </div>
      ))}
    </div>
  )
}

export default FavoritesList