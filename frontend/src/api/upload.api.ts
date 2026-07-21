import client from './client'

// Загрузка аватара
export const uploadAvatar = async (file: File): Promise<{ avatar_url: string }> => {
  const formData = new FormData()
  formData.append('avatar', file)
  
  const response = await client.post('/auth/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return response.data
}