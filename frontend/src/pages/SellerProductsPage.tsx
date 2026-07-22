import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import {
  FaPlus,
  FaEdit,
  FaTrash,
  FaLeaf,
  FaTimes,
  FaSave,
} from 'react-icons/fa'
import { getSellerProducts, createSellerProduct, updateSellerProduct, deleteSellerProduct } from '../api/catalog.api'
import { getCategories } from '../api/catalog.api'
import type { Product, Category } from '../api/catalog.api'
import { toast } from 'react-hot-toast'

const SellerProductsPage = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingProduct, setEditingProduct] = useState<Product | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [selectedImages, setSelectedImages] = useState<File[]>([])
  const [imagePreviews, setImagePreviews] = useState<string[]>([])

  // Форма
  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    old_price: '',
    stock: '',
    category_id: '',
    unit: 'букет',
    packaging: '',
    tags: '',
    is_active: true,
    is_featured: false,
  })

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      setLoading(true)
      const [productsData, categoriesData] = await Promise.all([
        getSellerProducts(),
        getCategories(),
      ])
      setProducts(Array.isArray(productsData) ? productsData : [])
      setCategories(Array.isArray(categoriesData) ? categoriesData : [])
    } catch (error) {
      console.error('Ошибка загрузки данных:', error)
      toast.error('Не удалось загрузить данные')
    } finally {
      setLoading(false)
    }
  }

  const resetForm = () => {
    setForm({
      name: '',
      description: '',
      price: '',
      old_price: '',
      stock: '',
      category_id: '',
      unit: 'букет',
      packaging: '',
      tags: '',
      is_active: true,
      is_featured: false,
    })
    setSelectedImages([])
    setImagePreviews([])
    setEditingProduct(null)
  }

  const handleOpenCreate = () => {
    resetForm()
    setIsModalOpen(true)
  }

  const handleOpenEdit = (product: Product) => {
    setEditingProduct(product)
    setForm({
      name: product.name || '',
      description: product.description || '',
      price: String(product.price || ''),
      old_price: product.old_price ? String(product.old_price) : '',
      stock: String(product.stock || ''),
      category_id: product.category_id || '',
      unit: product.unit || 'букет',
      packaging: product.packaging || '',
      tags: (product.tags || []).join(', '),
      is_active: product.is_active ?? true,
      is_featured: product.is_featured ?? false,
    })
    setSelectedImages([])
    setImagePreviews([])
    setIsModalOpen(true)
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || [])
    if (files.length + selectedImages.length > 5) {
      toast.error('Можно загрузить не более 5 фото')
      return
    }
    setSelectedImages([...selectedImages, ...files])
    const previews = files.map((file) => URL.createObjectURL(file))
    setImagePreviews([...imagePreviews, ...previews])
  }

  const removeImage = (index: number) => {
    setSelectedImages(selectedImages.filter((_, i) => i !== index))
    setImagePreviews(imagePreviews.filter((_, i) => i !== index))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)

    try {
      const formData = new FormData()
      formData.append('name', form.name)
      formData.append('description', form.description)
      formData.append('price', form.price)
      if (form.old_price) formData.append('old_price', form.old_price)
      formData.append('stock', form.stock)
      formData.append('category_id', form.category_id)
      formData.append('unit', form.unit)
      formData.append('packaging', form.packaging)
      formData.append('tags', form.tags)
      formData.append('is_active', String(form.is_active))
      formData.append('is_featured', String(form.is_featured))

      selectedImages.forEach((file) => {
        formData.append('images', file)
      })

      if (editingProduct) {
        await updateSellerProduct(editingProduct.id, formData)
        toast.success('Товар обновлён')
      } else {
        await createSellerProduct(formData)
        toast.success('Товар создан')
      }

      setIsModalOpen(false)
      resetForm()
      await fetchData()
    } catch (error: any) {
      console.error('Ошибка сохранения:', error)
      toast.error(error.response?.data?.error || 'Ошибка сохранения товара')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Вы уверены, что хотите удалить товар?')) return
    try {
      await deleteSellerProduct(id)
      toast.success('Товар удалён')
      await fetchData()
    } catch (error) {
      console.error('Ошибка удаления:', error)
      toast.error('Не удалось удалить товар')
    }
  }

  const getCategoryName = (categoryId: string) => {
    const cat = categories.find((c) => c.id === categoryId)
    return cat?.name || '—'
  }

  if (loading) {
    return <div className="text-center py-8 text-gray-400">Загрузка товаров...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
          <FaLeaf className="text-[#8A9A86]" />
          Мои товары
        </h2>
        <button
          onClick={handleOpenCreate}
          className="flex items-center gap-2 px-4 py-2.5 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium"
        >
          <FaPlus /> Добавить товар
        </button>
      </div>

      {products.length === 0 ? (
        <p className="text-gray-400 text-base">У вас пока нет товаров</p>
      ) : (
        <div className="overflow-x-auto bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] border border-gray-100">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Фото</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Название</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Категория</th>
                <th className="text-right p-3 font-medium text-[#1C1C1C]">Цена</th>
                <th className="text-center p-3 font-medium text-[#1C1C1C]">Остаток</th>
                <th className="text-center p-3 font-medium text-[#1C1C1C]">Статус</th>
                <th className="text-center p-3 font-medium text-[#1C1C1C]">Действия</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product) => (
                <tr key={product.id} className="border-b border-gray-50 hover:bg-gray-50/50 transition">
                  <td className="p-3">
                    <div className="w-12 h-12 bg-gray-50 rounded-lg flex items-center justify-center overflow-hidden">
                      <FaLeaf className="text-gray-300 text-2xl" />
                    </div>
                  </td>
                  <td className="p-3">
                    <span className="font-medium text-[#1C1C1C]">{product.name}</span>
                  </td>
                  <td className="p-3 text-gray-600">{getCategoryName(product.category_id)}</td>
                  <td className="p-3 text-right font-medium text-[#8A9A86]">{product.price} BYN</td>
                  <td className="p-3 text-center">{product.stock}</td>
                  <td className="p-3 text-center">
                    {product.is_active ? (
                      <span className="text-green-600 bg-green-50 px-2 py-1 rounded-full text-xs">Активен</span>
                    ) : (
                      <span className="text-red-600 bg-red-50 px-2 py-1 rounded-full text-xs">Неактивен</span>
                    )}
                  </td>
                  <td className="p-3 text-center">
                    <div className="flex justify-center gap-2">
                      <button
                        onClick={() => handleOpenEdit(product)}
                        className="text-blue-500 hover:text-blue-700 transition"
                      >
                        <FaEdit />
                      </button>
                      <button
                        onClick={() => handleDelete(product.id)}
                        className="text-red-500 hover:text-red-700 transition"
                      >
                        <FaTrash />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Модалка добавления/редактирования */}
      {isModalOpen &&
        createPortal(
          <div
            className="fixed inset-0 z-[9999] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4"
            onClick={() => setIsModalOpen(false)}
          >
            <div
              className="bg-white rounded-2xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6 shadow-[0_8px_40px_rgba(0,0,0,0.08)] border border-gray-100 animate-fade-in-up"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="flex justify-between items-center mb-4 pb-3 border-b border-gray-100">
                <h3 className="text-xl font-bold text-[#1C1C1C]">
                  {editingProduct ? 'Редактировать товар' : 'Добавить товар'}
                </h3>
                <button
                  onClick={() => setIsModalOpen(false)}
                  className="text-gray-400 hover:text-gray-600 text-2xl leading-none transition"
                >
                  <FaTimes />
                </button>
              </div>

              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Название *</label>
                  <input
                    type="text"
                    value={form.name}
                    onChange={(e) => setForm({ ...form, name: e.target.value })}
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Описание</label>
                  <textarea
                    value={form.description}
                    onChange={(e) => setForm({ ...form, description: e.target.value })}
                    rows={3}
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition resize-none"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Цена (BYN) *</label>
                    <input
                      type="number"
                      step="0.01"
                      value={form.price}
                      onChange={(e) => setForm({ ...form, price: e.target.value })}
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Старая цена</label>
                    <input
                      type="number"
                      step="0.01"
                      value={form.old_price}
                      onChange={(e) => setForm({ ...form, old_price: e.target.value })}
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Количество *</label>
                    <input
                      type="number"
                      value={form.stock}
                      onChange={(e) => setForm({ ...form, stock: e.target.value })}
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Категория *</label>
                    <select
                      value={form.category_id}
                      onChange={(e) => setForm({ ...form, category_id: e.target.value })}
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white"
                      required
                    >
                      <option value="">Выберите категорию</option>
                      {categories.map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.name}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Единица</label>
                    <select
                      value={form.unit}
                      onChange={(e) => setForm({ ...form, unit: e.target.value })}
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white"
                    >
                      <option value="букет">Букет</option>
                      <option value="шт">Штука</option>
                      <option value="композиция">Композиция</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Упаковка</label>
                    <input
                      type="text"
                      value={form.packaging}
                      onChange={(e) => setForm({ ...form, packaging: e.target.value })}
                      placeholder="Крафт, плёнка..."
                      className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Теги (через запятую)</label>
                  <input
                    type="text"
                    value={form.tags}
                    onChange={(e) => setForm({ ...form, tags: e.target.value })}
                    placeholder="розы, романтика, свадьба"
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-[#1C1C1C] mb-1">Фото (до 5 шт.)</label>
                  <input
                    type="file"
                    accept="image/*"
                    multiple
                    onChange={handleImageChange}
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition"
                  />
                  <div className="flex flex-wrap gap-2 mt-2">
                    {imagePreviews.map((preview, index) => (
                      <div key={index} className="relative w-16 h-16 rounded-lg overflow-hidden border border-gray-200">
                        <img src={preview} alt={`Preview ${index}`} className="w-full h-full object-cover" />
                        <button
                          type="button"
                          onClick={() => removeImage(index)}
                          className="absolute top-0 right-0 bg-red-500 text-white rounded-full w-5 h-5 flex items-center justify-center text-xs"
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="flex gap-4">
                  <label className="flex items-center gap-2 text-sm text-[#1C1C1C]">
                    <input
                      type="checkbox"
                      checked={form.is_active}
                      onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                      className="w-4 h-4 accent-[#8A9A86]"
                    />
                    Активен (виден в каталоге)
                  </label>
                  <label className="flex items-center gap-2 text-sm text-[#1C1C1C]">
                    <input
                      type="checkbox"
                      checked={form.is_featured}
                      onChange={(e) => setForm({ ...form, is_featured: e.target.checked })}
                      className="w-4 h-4 accent-[#8A9A86]"
                    />
                    Рекомендуемый
                  </label>
                </div>

                <div className="flex gap-3 pt-4 border-t border-gray-100">
                  <button
                    type="button"
                    onClick={() => setIsModalOpen(false)}
                    className="flex-1 px-4 py-2.5 border border-gray-200 rounded-xl hover:bg-gray-50 transition text-sm font-medium"
                  >
                    Отмена
                  </button>
                  <button
                    type="submit"
                    disabled={submitting}
                    className="flex-1 bg-[#8A9A86] text-white px-4 py-2.5 rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium disabled:opacity-50 flex items-center justify-center gap-2"
                  >
                    <FaSave /> {submitting ? 'Сохранение...' : 'Сохранить'}
                  </button>
                </div>
              </form>
            </div>
          </div>,
          document.body
        )}
    </div>
  )
}

export default SellerProductsPage