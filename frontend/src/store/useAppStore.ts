import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

// Типы для store
interface AppState {
  // Состояние приложения
  isAuthenticated: boolean
  user: {
    id: number | null
    name: string | null
    email: string | null
  } | null
  
  // UI состояние
  isLoading: boolean
  theme: 'light' | 'dark'
  
  // Действия
  setAuthenticated: (isAuthenticated: boolean) => void
  setUser: (user: AppState['user']) => void
  setLoading: (isLoading: boolean) => void
  setTheme: (theme: 'light' | 'dark') => void
  reset: () => void
}

// Начальное состояние
const initialState = {
  isAuthenticated: false,
  user: null,
  isLoading: false,
  theme: 'dark' as const,
}

// Создание store с поддержкой Redux DevTools и персистентности
export const useAppStore = create<AppState>()(
  devtools(
    persist(
      (set) => ({
        ...initialState,
        
        setAuthenticated: (isAuthenticated) =>
          set({ isAuthenticated }, false, 'setAuthenticated'),
        
        setUser: (user) =>
          set({ user }, false, 'setUser'),
        
        setLoading: (isLoading) =>
          set({ isLoading }, false, 'setLoading'),
        
        setTheme: (theme) =>
          set({ theme }, false, 'setTheme'),
        
        reset: () =>
          set(initialState, false, 'reset'),
      }),
      {
        name: 'app-storage', // ключ для localStorage
        partialize: (state) => ({
          // Сохраняем только нужные поля в localStorage
          theme: state.theme,
          user: state.user,
          isAuthenticated: state.isAuthenticated,
        }),
      }
    ),
    { name: 'AppStore' } // имя для Redux DevTools
  )
)

