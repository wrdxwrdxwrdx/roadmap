import { ReactNode } from 'react'
import { useTheme } from '../../theme/useTheme'

interface MainLayoutProps {
  children: ReactNode
}

/**
 * Главный layout для приложения
 * Включает навигацию, хедер, футер и основной контент
 */
export function MainLayout({ children }: MainLayoutProps) {
  const { theme, toggleTheme } = useTheme()

  return (
    <div className="main-layout" style={{
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      backgroundColor: 'var(--color-background)',
      color: 'var(--color-text)',
    }}>
      {/* Header */}
      <header style={{
        padding: 'var(--spacing-md)',
        borderBottom: '1px solid var(--color-border)',
        backgroundColor: 'var(--color-surface)',
      }}>
        <div style={{
          maxWidth: '1200px',
          margin: '0 auto',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}>
          <h1 style={{ margin: 0, fontSize: '1.5rem' }}>Roadmap</h1>
          
          {/* Theme Toggle */}
          <button
            onClick={(e) => {
              e.preventDefault()
              toggleTheme()
            }}
            style={{
              padding: 'var(--spacing-sm) var(--spacing-md)',
              backgroundColor: 'var(--color-surface)',
              color: 'var(--color-text)',
              border: '1px solid var(--color-border)',
              borderRadius: 'var(--radius-md)',
              cursor: 'pointer',
              fontSize: '0.9rem',
              transition: 'all 0.2s ease',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = 'var(--color-hover)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = 'var(--color-surface)'
            }}
            aria-label={`Переключить на ${theme === 'light' ? 'темную' : 'светлую'} тему`}
          >
            {theme === 'light' ? '🌙' : '☀️'} {theme === 'light' ? 'Темная' : 'Светлая'}
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main style={{
        flex: 1,
        width: '100%',
        maxWidth: '1200px',
        margin: '0 auto',
        padding: 'var(--spacing-lg)',
      }}>
        {children}
      </main>

      {/* Footer */}
      <footer style={{
        padding: 'var(--spacing-md)',
        borderTop: '1px solid var(--color-border)',
        backgroundColor: 'var(--color-surface)',
        textAlign: 'center',
        fontSize: '0.9rem',
        color: 'var(--color-text-muted)',
      }}>
        <p style={{ margin: 0 }}>© 2024 Roadmap. All rights reserved.</p>
      </footer>
    </div>
  )
}

