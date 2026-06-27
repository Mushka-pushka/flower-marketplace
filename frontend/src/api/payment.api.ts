import client from './client'

export interface Payment {
  id: string
  order_id: string
  amount: number
  status: 'pending' | 'completed' | 'failed' | 'refunded'
  payment_method: string
  transaction_id: string
  payment_url: string
  created_at: string
  updated_at: string
}

export const createPayment = async (data: {
  order_id: string
  amount: number
  payment_method: string
}): Promise<Payment> => {
  const response = await client.post('/payments', data)
  return response.data
}

export const getPaymentStatus = async (paymentId: string): Promise<Payment> => {
  const response = await client.get('/payments', { params: { id: paymentId } })
  return response.data
}