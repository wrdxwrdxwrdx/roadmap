// Система тем для приложения
export type Theme = 'light' | 'dark'

export interface ThemeColors {
  // Основные цвета
  background: string
  surface: string
  primary: string
  secondary: string
  
  // Текст
  text: string
  textSecondary: string
  textMuted: string
  
  // Границы
  border: string
  borderLight: string
  
  // Состояния
  success: string
  error: string
  warning: string
  info: string
  
  // Интерактивные элементы
  hover: string
  active: string
  disabled: string
}

export interface ThemeConfig {
  name: Theme
  colors: ThemeColors
}

// Светлая тема
export const lightTheme: ThemeConfig = {
  name: 'light',
  colors: {
    background: '#ffffff',
    surface: '#f8f9fa',
    primary: '#007bff',
    secondary: '#6c757d',
    text: '#212529',
    textSecondary: '#495057',
    textMuted: '#6c757d',
    border: '#dee2e6',
    borderLight: '#e9ecef',
    success: '#28a745',
    error: '#dc3545',
    warning: '#ffc107',
    info: '#17a2b8',
    hover: '#f1f3f5',
    active: '#e9ecef',
    disabled: '#adb5bd',
  },
}

// Темная тема
export const darkTheme: ThemeConfig = {
  name: 'dark',
  colors: {
    background: '#1a1a1a',
    surface: '#242424',
    primary: '#0d6efd',
    secondary: '#6c757d',
    text: '#ffffff',
    textSecondary: '#e0e0e0',
    textMuted: '#adb5bd',
    border: '#404040',
    borderLight: '#2d2d2d',
    success: '#198754',
    error: '#dc3545',
    warning: '#ffc107',
    info: '#0dcaf0',
    hover: '#2d2d2d',
    active: '#404040',
    disabled: '#6c757d',
  },
}

export const themes: Record<Theme, ThemeConfig> = {
  light: lightTheme,
  dark: darkTheme,
}

// Функция для получения темы
export const getTheme = (theme: Theme): ThemeConfig => themes[theme]

