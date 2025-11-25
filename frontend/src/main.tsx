import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App'
import './index.css'

// Применяем тему ДО рендера React, чтобы избежать мигания
function applyInitialTheme() {
  const root = document.documentElement
  
  // Получаем сохраненную тему из localStorage
  const savedTheme = localStorage.getItem('app-storage')
  let theme = 'dark' // значение по умолчанию
  
  if (savedTheme) {
    try {
      const parsed = JSON.parse(savedTheme)
      if (parsed.state?.theme) {
        theme = parsed.state.theme
      }
    } catch (e) {
      console.warn('Failed to parse saved theme:', e)
    }
  }
  
  // Применяем тему синхронно
  root.classList.remove('light', 'dark')
  root.classList.add(theme)
  root.setAttribute('data-theme', theme)
  root.style.colorScheme = theme
}

// Применяем тему сразу
applyInitialTheme()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
)

