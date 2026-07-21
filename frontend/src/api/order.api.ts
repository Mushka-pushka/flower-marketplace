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
  product_name?: string
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

// Интерфейс для ответа с пагинацией
export interface OrdersResponse {
  orders: Order[]
  total: number
  limit: number
  offset: number
}

// позиция заказа с данными о товаре и статусе
export interface OrderItemWithStatus {
  id: string
  order_id: string
  product_id: string
  product_name: string
  product_price: number
  quantity: number
  total: number
  order_status: string
  shop_id: string
  delivery_date: string
  delivery_time: string
  comment: string
  created_at: string
  updated_at: string
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

// Получение заказов пользователя (с пагинацией)
export const getMyOrders = async (customerId: string): Promise<Order[]> => {
  try {
    const response = await client.get('/orders/customer', {
      params: { customer_id: customerId }
    })
    
    console.log('getMyOrders raw response:', response.data) 
    
    if (response.data && Array.isArray(response.data.orders)) {
      return response.data.orders
    }
    
    if (Array.isArray(response.data)) {
      return response.data
    }
    
    if (response.data && typeof response.data === 'object') {
      // Проверяем все возможные поля
      for (const key of ['orders', 'items', 'data', 'results']) {
        if (Array.isArray(response.data[key])) {
          return response.data[key]
        }
      }
    }
    
    console.warn('Unexpected orders response format:', response.data)
    return []
  } catch (error) {
    console.error('Error fetching orders:', error)
    return []
  }
}

// получение всех товаров пользователя (как отдельные позиции)
export const getMyOrderItems = async (): Promise<OrderItemWithStatus[]> => {
  const response = await client.get('/orders/items')
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

// Обновление статуса заказа (для продавца)
export const updateOrderStatus = async (data: {
  order_id: string
  status: string
  comment?: string
}): Promise<void> => {
  await client.put('/orders/status', data)
}

// Отмена заказа
export const cancelOrder = async (orderId: string): Promise<void> => {
  await client.post('/orders/cancel', null, { params: { id: orderId } })
}