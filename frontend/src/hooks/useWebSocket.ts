import { useEffect, useRef, useState } from 'react'

interface WebSocketMessage {
  event: string
  order_id: string
  status?: string
  timestamp: number
}

export const useWebSocket = (url: string) => {
  const [isConnected, setIsConnected] = useState(false)
  const [messages, setMessages] = useState<WebSocketMessage[]>([])
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        setMessages(prev => [...prev, data])
        
        // Если это обновление статуса заказа
        if (data.event === 'order.status_changed') {
          // Триггерим обновление в UI
          window.dispatchEvent(new CustomEvent('orderStatusChanged', { 
            detail: { orderId: data.order_id, status: data.status }
          }))
        }
      } catch (error) {
        console.error('WebSocket message error:', error)
      }
    }

    ws.onclose = () => {
      setIsConnected(false)
      console.log('WebSocket disconnected')
      // Переподключение через 3 секунды
      setTimeout(() => {
        if (wsRef.current?.readyState !== WebSocket.OPEN) {
          // Reconnect logic
        }
      }, 3000)
    }

    return () => {
      ws.close()
    }
  }, [url])

  const sendMessage = (message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    }
  }

  return { isConnected, messages, sendMessage }
}