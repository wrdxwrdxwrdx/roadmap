import { Routes, Route } from 'react-router-dom'
import { useTheme } from './theme/useTheme'
import { MainLayout } from './components/layouts/MainLayout'
import HomePage from './pages/HomePage'
import NotFoundPage from './pages/NotFoundPage'

function App() {
  // Инициализируем тему при загрузке приложения
  useTheme()

  return (
    <MainLayout>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </MainLayout>
  )
}

export default App

