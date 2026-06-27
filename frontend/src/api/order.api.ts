import client from './client'

export interface Order {
  id: string
  customer_id: string
  shop_id: string
  delivery_address_id: string
  payment_type_id: number
  total_amount: number
  delivery_date: string
  delivery_time: string
  comment: string
  current_status: string
  created_at: string
  updated_at: string
}

export const createOrder = async (data: {
  customer_id: string
  shop_id: string
  delivery_address_id: string
  payment_type_id: number
  delivery_date: string
  delivery_time: string
  comment: string
  items: { product_id: string; quantity: number }[]
}): Promise<Order> => {
  const response = await client.post('/orders', data)
  return response.data
}