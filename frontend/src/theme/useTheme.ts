import { useEffect } from 'react'
import { useAppStore } from '../store'
import { getTheme } from './theme'

/**
 * Хук для работы с темой приложения
 * Применяет тему к DOM и синхронизирует с store
 */
export function useTheme() {
  const { theme, setTheme } = useAppStore()

  // Применяем тему к DOM при изменении
  useEffect(() => {
    const root = document.documentElement
    const themeConfig = getTheme(theme)
    
    // Удаляем все классы тем
    root.classList.remove('light', 'dark')
    
    // Добавляем класс текущей темы
    root.classList.add(theme)
    
    // Устанавливаем атрибут для совместимости
    root.setAttribute('data-theme', theme)
    
    // Устанавливаем color-scheme для лучшей поддержки
    root.style.colorScheme = theme
    
    // Обновляем или создаем meta тег для theme-color
    let metaTheme = document.querySelector('meta[name="theme-color"]') as HTMLMetaElement
    if (!metaTheme) {
      metaTheme = document.createElement('meta')
      metaTheme.name = 'theme-color'
      document.head.appendChild(metaTheme)
    }
    metaTheme.content = themeConfig.colors.background
  }, [theme])

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light')
  }

  return {
    theme,
    setTheme,
    toggleTheme,
    themeConfig: getTheme(theme),
  }
}

