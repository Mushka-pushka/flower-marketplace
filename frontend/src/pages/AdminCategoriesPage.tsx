import { useEffect, useState } from 'react'
import {
  FaFolder,
  FaPlus,
  FaEdit,
  FaTrash,
  FaTimes,
  FaSave,
  FaSpinner,
  FaSearch,
  FaFolderOpen,
} from 'react-icons/fa'
import { toast } from 'react-hot-toast'
import client from '../api/client'

interface Category {
  id: string
  name: string
  slug: string
  description?: string
  parent_id?: string | null
  sort_order: number
  created_at: string
  product_count?: number
}

const AdminCategoriesPage = () => {
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingCategory, setEditingCategory] = useState<Category | null>(null)
  const [submitting, setSubmitting] = useState(false)

  // Форма
  const [form, setForm] = useState({
    name: '',
    slug: '',
    description: '',
    parent_id: '',
    sort_order: 0,
  })

  // Загрузка категорий
  const fetchCategories = async () => {
    try {
      setLoading(true)
      const response = await client.get('/admin/categories')
      setCategories(Array.isArray(response.data) ? response.data : [])
    } catch (error) {
      console.error('Ошибка загрузки категорий:', error)
      toast.error('Не удалось загрузить категории')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCategories()
  }, [])

  // Открытие модалки для создания
  const handleOpenCreate = () => {
    setEditingCategory(null)
    setForm({
      name: '',
      slug: '',
      description: '',
      parent_id: '',
      sort_order: 0,
    })
    setIsModalOpen(true)
  }

  // Открытие модалки для редактирования
  const handleOpenEdit = (category: Category) => {
    setEditingCategory(category)
    setForm({
      name: category.name || '',
      slug: category.slug || '',
      description: category.description || '',
      parent_id: category.parent_id || '',
      sort_order: category.sort_order || 0,
    })
    setIsModalOpen(true)
  }

  // Генерация slug из названия
  const generateSlug = (name: string) => {
    return name
      .toLowerCase()
      .replace(/[^a-zа-яё0-9]/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-|-$/g, '')
  }

  // Обработчик изменения названия (автогенерация slug)
  const handleNameChange = (name: string) => {
    setForm(prev => ({
      ...prev,
      name,
      slug: generateSlug(name),
    }))
  }

  // Создание/обновление категории
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!form.name.trim() || !form.slug.trim()) {
      toast.error('Название и Slug обязательны')
      return
    }

    setSubmitting(true)
    try {
      const payload = {
        name: form.name.trim(),
        slug: form.slug.trim(),
        description: form.description.trim(),
        parent_id: form.parent_id || null,
        sort_order: Number(form.sort_order) || 0,
      }

      if (editingCategory) {
        await client.put(`/admin/categories?id=${editingCategory.id}`, payload)
        toast.success('Категория обновлена')
      } else {
        await client.post('/admin/categories', payload)
        toast.success('Категория создана')
      }

      setIsModalOpen(false)
      await fetchCategories()
    } catch (error: any) {
      console.error('Ошибка сохранения:', error)
      toast.error(error.response?.data?.error || 'Ошибка сохранения категории')
    } finally {
      setSubmitting(false)
    }
  }

  // Удаление категории
  const handleDelete = async (category: Category) => {
    if (!confirm(`Удалить категорию "${category.name}"?`)) return

    try {
      await client.delete(`/admin/categories?id=${category.id}`)
      toast.success('Категория удалена')
      await fetchCategories()
    } catch (error: any) {
      console.error('Ошибка удаления:', error)
      toast.error(error.response?.data?.error || 'Не удалось удалить категорию')
    }
  }

  // Фильтрация категорий
  const filteredCategories = categories.filter((cat) => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return (
      cat.name.toLowerCase().includes(query) ||
      cat.slug.toLowerCase().includes(query) ||
      (cat.description && cat.description.toLowerCase().includes(query))
    )
  })

  // Подкатегории (для отображения иерархии)
  const getSubcategories = (parentId: string | null) => {
    return filteredCategories.filter((cat) => cat.parent_id === parentId)
  }

  if (loading) {
    return (
      <div className="text-center py-12">
        <FaSpinner className="text-3xl text-[#8A9A86] animate-spin mx-auto" />
        <p className="text-gray-400 mt-3">Загрузка категорий...</p>
      </div>
    )
  }

  return (
    <div>
      {/* Заголовок */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h2 className="text-2xl font-bold text-[#1C1C1C] flex items-center gap-2">
            <FaFolderOpen className="text-[#8A9A86]" />
            Управление категориями
          </h2>
          <p className="text-gray-400 text-sm">
            Всего категорий: {categories.length}
          </p>
        </div>
        <button
          onClick={handleOpenCreate}
          className="flex items-center gap-2 px-5 py-2.5 bg-[#8A9A86] text-white rounded-xl hover:bg-[#7A8A76] transition text-sm font-medium"
        >
          <FaPlus /> Добавить категорию
        </button>
      </div>

      {/* Поиск */}
      <div className="mb-4">
        <div className="relative max-w-sm">
          <FaSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-sm" />
          <input
            type="text"
            placeholder="Поиск категорий..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C] text-sm"
          />
        </div>
      </div>

      {/* Список категорий */}
      {filteredCategories.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border border-gray-100">
          <FaFolder className="text-4xl text-gray-300 mx-auto mb-3" />
          <p className="text-gray-400 text-lg">
            {searchQuery ? 'Категории не найдены' : 'Категорий пока нет'}
          </p>
          <p className="text-gray-400 text-sm">
            {searchQuery ? 'Попробуйте изменить поисковый запрос' : 'Создайте первую категорию'}
          </p>
        </div>
      ) : (
        <div className="bg-white rounded-xl shadow-[0_4px_20px_rgba(0,0,0,0.04)] overflow-hidden border border-gray-100">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-100">
              <tr>
                <th className="text-left p-3 font-medium text-[#1C1C1C] w-12">#</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Название</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Slug</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Описание</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Порядок</th>
                <th className="text-left p-3 font-medium text-[#1C1C1C]">Дата</th>
                <th className="text-center p-3 font-medium text-[#1C1C1C]">Действия</th>
              </tr>
            </thead>
            <tbody>
              {filteredCategories.map((category, index) => {
                const subcategories = getSubcategories(category.id)
                const isParent = subcategories.length > 0

                return (
                  <tr
                    key={category.id}
                    className="border-b border-gray-50 hover:bg-gray-50/50 transition"
                  >
                    <td className="p-3 text-gray-400 text-center">{index + 1}</td>
                    <td className="p-3">
                      <div className="flex items-center gap-2">
                        <FaFolder className={`${isParent ? 'text-[#8A9A86]' : 'text-gray-300'} text-sm`} />
                        <span className="font-medium text-[#1C1C1C]">{category.name}</span>
                        {isParent && (
                          <span className="text-xs text-gray-400 ml-1">
                            ({subcategories.length} подкатегорий)
                          </span>
                        )}
                      </div>
                    </td>
                    <td className="p-3 text-gray-500 font-mono text-xs">{category.slug}</td>
                    <td className="p-3 text-gray-500 max-w-[200px] truncate">
                      {category.description || '—'}
                    </td>
                    <td className="p-3 text-gray-500">{category.sort_order}</td>
                    <td className="p-3 text-gray-400 text-xs">
                      {new Date(category.created_at).toLocaleDateString('ru-RU', {
                        day: '2-digit',
                        month: '2-digit',
                        year: 'numeric',
                      })}
                    </td>
                    <td className="p-3 text-center">
                      <div className="flex justify-center gap-2">
                        <button
                          onClick={() => handleOpenEdit(category)}
                          className="text-blue-500 hover:text-blue-700 transition p-1"
                          title="Редактировать"
                        >
                          <FaEdit />
                        </button>
                        <button
                          onClick={() => handleDelete(category)}
                          className="text-red-500 hover:text-red-700 transition p-1"
                          title="Удалить"
                        >
                          <FaTrash />
                        </button>
                      </div>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      {/* Модалка создания/редактирования */}
      {isModalOpen && (
        <div
          className="fixed inset-0 z-[9999] flex items-center justify-center p-4"
          onClick={() => setIsModalOpen(false)}
        >
          <div
            className="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto p-6 shadow-[0_20px_60px_rgba(0,0,0,0.15)] border border-gray-100 animate-fade-in-up"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-center mb-4 pb-3 border-b border-gray-100">
              <h3 className="text-xl font-bold text-[#1C1C1C] flex items-center gap-2">
                <FaFolder className="text-[#8A9A86]" />
                {editingCategory ? 'Редактировать категорию' : 'Новая категория'}
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
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                  Название *
                </label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => handleNameChange(e.target.value)}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                  placeholder="Например: Розы"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                  Slug (URL) *
                </label>
                <input
                  type="text"
                  value={form.slug}
                  onChange={(e) => setForm({ ...form, slug: e.target.value })}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] font-mono text-sm"
                  placeholder="Например: roses"
                  required
                />
                <p className="text-xs text-gray-400 mt-1">
                  Используется в URL. Только латиница, цифры и дефис.
                </p>
              </div>

              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                  Описание
                </label>
                <textarea
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  rows={3}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C] resize-none"
                  placeholder="Краткое описание категории"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                  Родительская категория
                </label>
                <select
                  value={form.parent_id}
                  onChange={(e) => setForm({ ...form, parent_id: e.target.value })}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition bg-white text-[#1C1C1C]"
                >
                  <option value="">Нет (корневая категория)</option>
                  {categories
                    .filter((cat) => cat.id !== editingCategory?.id)
                    .map((cat) => (
                      <option key={cat.id} value={cat.id}>
                        {cat.name}
                      </option>
                    ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-[#1C1C1C] mb-1">
                  Порядок сортировки
                </label>
                <input
                  type="number"
                  value={form.sort_order}
                  onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8A9A86] transition text-[#1C1C1C]"
                  placeholder="0"
                />
                <p className="text-xs text-gray-400 mt-1">
                  Меньше число — выше в списке
                </p>
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
                  <FaSave />
                  {submitting ? 'Сохранение...' : 'Сохранить'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default AdminCategoriesPage