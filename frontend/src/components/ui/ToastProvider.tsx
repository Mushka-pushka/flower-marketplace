import { Toaster } from 'react-hot-toast'

const ToastProvider = () => {
  return (
    <Toaster
      position="top-right"
      toastOptions={{
        duration: 3000,
        style: {
          background: '#1C1C1C',
          color: '#fff',
          borderRadius: '12px',
          padding: '12px 20px',
        },
        success: {
          style: {
            background: '#8A9A86',
          },
        },
        error: {
          style: {
            background: '#DC2626',
          },
        },
      }}
    />
  )
}

export default ToastProvider