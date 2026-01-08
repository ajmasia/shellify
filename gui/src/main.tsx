import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { Toaster } from 'sonner'
import { StoreProvider } from '@/stores'
import { ModalProvider } from '@/components/Modal'
import { App } from './App'
import './styles/global.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <StoreProvider>
        <ModalProvider>
          <App />
          <Toaster
            position="bottom-right"
            theme="dark"
            toastOptions={{
              style: {
                background: 'var(--ctp-surface0)',
                color: 'var(--text-primary)',
                border: '1px solid var(--ctp-surface1)',
              },
              classNames: {
                success: 'toast-success',
                error: 'toast-error',
              },
            }}
          />
        </ModalProvider>
      </StoreProvider>
    </BrowserRouter>
  </StrictMode>
)
