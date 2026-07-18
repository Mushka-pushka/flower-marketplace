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

export interface OrderItem {
  id: string
  order_id: string
  product_id: string
  quantity: number
  price: number
  total: number
  packaging: string
  created_at: string
}

export interface OrderStatus {
  id: string
  order_id: string
  status: string
  changed_by: string
  comment: string
  created_at: string
}

export interface OrderDetails {
  order: Order
  items: OrderItem[]
  statuses: OrderStatus[]
}

// Создание заказа
export const createOrder = async (data: {
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

// Получение заказов пользователя
export const getMyOrders = async (customerId: string): Promise<Order[]> => {
  const response = await client.get('/orders/customer', { params: { customer_id: customerId } })
  return response.data
}

// Получение деталей заказа
export const getOrderDetails = async (orderId: string): Promise<OrderDetails> => {
  const response = await client.get('/orders', { params: { id: orderId } })
  return response.data
}

// Проверка, может ли пользователь оставить отзыв на товар
export const canReviewProduct = async (productId: string): Promise<boolean> => {
  try {
    const response = await client.get('/orders/can-review', { params: { product_id: productId } })
    return response.data.can_review
  } catch {
    return false
  }
}