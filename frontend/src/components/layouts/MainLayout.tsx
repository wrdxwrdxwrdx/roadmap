import { ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useTheme } from '../../theme/useTheme'
import { useLanguage } from '../../hooks/useLanguage'

interface MainLayoutProps {
  children: ReactNode
}

/**
 * Главный layout для приложения
 * Включает навигацию, хедер, футер и основной контент
 */
export function MainLayout({ children }: MainLayoutProps) {
  const { theme, toggleTheme } = useTheme()
  const { t } = useTranslation()
  const { currentLanguage, toggleLanguage } = useLanguage()

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
          <h1 style={{ margin: 0, fontSize: '1.5rem' }}>
            <Link to="/" style={{ color: 'var(--color-text)', textDecoration: 'none' }}>
              {t('common.appName')}
            </Link>
          </h1>
          
          <nav style={{
            display: 'flex',
            gap: 'var(--spacing-md)',
            alignItems: 'center',
          }}>
            <Link
              to="/register"
              style={{
                color: 'var(--color-text)',
                textDecoration: 'none',
                padding: 'var(--spacing-xs) var(--spacing-sm)',
                borderRadius: 'var(--radius-md)',
                fontSize: '0.9rem',
                transition: 'background-color 0.2s ease',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--color-hover)'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = 'transparent'
              }}
            >
              {t('navigation.register')}
            </Link>
          </nav>
          
          <div style={{
            display: 'flex',
            gap: 'var(--spacing-sm)',
            alignItems: 'center',
          }}>
            {/* Language Toggle */}
            <button
              onClick={(e) => {
                e.preventDefault()
                toggleLanguage()
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
                display: 'flex',
                alignItems: 'center',
                gap: 'var(--spacing-xs)',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--color-hover)'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--color-surface)'
              }}
              aria-label={t('language.switch')}
              title={t('language.switch')}
            >
              <span>{currentLanguage === 'ru' ? '🇷🇺' : '🇬🇧'}</span>
              <span>{currentLanguage === 'ru' ? t('language.russian') : t('language.english')}</span>
            </button>

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
              aria-label={t('theme.toggle')}
            >
              {theme === 'light' ? '🌙' : '☀️'} {theme === 'light' ? t('theme.dark') : t('theme.light')}
            </button>
          </div>
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
        <p style={{ margin: 0 }}>© 2024 {t('common.appName')}. All rights reserved.</p>
      </footer>
    </div>
  )
}

