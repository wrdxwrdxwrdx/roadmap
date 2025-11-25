import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useApi } from '../hooks/useApi'
import { apiEndpoints } from '../services/apiEndpoints'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { useTheme } from '../theme/useTheme'

function HomePage() {
  const { t } = useTranslation()
  // Демонстрация использования темы
  const { theme } = useTheme()
  
  // Демонстрация использования API хука
  const { data, loading, error, execute } = useApi(apiEndpoints.healthCheck)

  return (
    <div>
      <h1>{t('home.title')}</h1>
      <p style={{ color: 'var(--color-text-secondary)' }}>{t('common.welcome')}</p>
      
      <div style={{ marginTop: '2rem', display: 'flex', flexDirection: 'column', gap: 'var(--spacing-lg)' }}>
        <h2>{t('home.demoTitle')}</h2>
        
        <Card title={t('home.theme.title')}>
          <p style={{ color: 'var(--color-text-secondary)' }}>
            {t('home.theme.current')}: <strong style={{ color: 'var(--color-text)' }}>{theme}</strong>
          </p>
          <p style={{ color: 'var(--color-text-muted)', fontSize: '0.9rem' }}>
            {t('home.theme.description')}
          </p>
        </Card>

        <Card title={t('home.api.title')}>
          <p style={{ fontSize: '0.9em', color: 'var(--color-text-secondary)', marginBottom: '0.5em' }}>
            {t('home.api.description')}
          </p>
          <Button 
            onClick={() => execute()} 
            disabled={loading}
            variant="primary"
          >
            {loading ? t('home.api.checking') : t('home.api.checkButton')}
          </Button>
          
          {loading && (
            <div style={{ marginTop: '1em', color: 'var(--color-text-secondary)' }}>
              <p>{t('home.api.sending')}</p>
            </div>
          )}
          
          {error && !loading && (
            <div style={{ 
              marginTop: '1em', 
              padding: '1em', 
              backgroundColor: 'var(--color-error)', 
              color: '#fff',
              border: '1px solid var(--color-error)', 
              borderRadius: 'var(--radius-md)',
              opacity: 0.1,
            }}>
              <div style={{ 
                backgroundColor: 'rgba(220, 53, 69, 0.1)',
                padding: '1em',
                borderRadius: 'var(--radius-md)',
                border: '1px solid var(--color-error)',
              }}>
                <p style={{ color: 'var(--color-error)', fontWeight: 'bold', marginBottom: '0.5em' }}>{t('home.api.error.title')}</p>
                <p style={{ color: 'var(--color-text)' }}><strong>{t('home.api.error.message')}:</strong> {error.message}</p>
                {error.response && (
                  <p style={{ color: 'var(--color-text)', marginTop: '0.5em' }}>
                    <strong>{t('home.api.error.httpStatus')}:</strong> {error.response.status} - {error.response.statusText}
                  </p>
                )}
                <p style={{ fontSize: '0.9em', color: 'var(--color-text-muted)', marginTop: '0.5em' }}>
                  {t('home.api.error.hint')}
                </p>
              </div>
            </div>
          )}
          
          {data && !loading && !error && (
            <div style={{ 
              marginTop: '1em', 
              padding: '1em', 
              backgroundColor: 'rgba(40, 167, 69, 0.1)',
              border: '1px solid var(--color-success)', 
              borderRadius: 'var(--radius-md)',
            }}>
              <p style={{ color: 'var(--color-success)', fontWeight: 'bold', marginBottom: '0.5em' }}>{t('home.api.success.title')}</p>
              <div style={{ fontSize: '0.95em', marginTop: '0.5em', color: 'var(--color-text)' }}>
                <p><strong>{t('home.api.success.status')}:</strong> <span style={{ color: 'var(--color-success)' }}>{data.status}</span></p>
                <p><strong>{t('home.api.success.service')}:</strong> {data.service}</p>
              </div>
              <details style={{ marginTop: '0.75em', fontSize: '0.85em' }}>
                <summary style={{ cursor: 'pointer', color: 'var(--color-text-muted)', userSelect: 'none' }}>
                  {t('home.api.success.showJson')}
                </summary>
                <pre style={{ 
                  marginTop: '0.5em', 
                  padding: '0.75em', 
                  background: 'var(--color-surface)', 
                  border: '1px solid var(--color-border)',
                  borderRadius: 'var(--radius-md)', 
                  overflow: 'auto',
                  fontSize: '0.85em',
                  textAlign: 'left',
                  color: 'var(--color-text)',
                }}>
                  {JSON.stringify(data, null, 2)}
                </pre>
              </details>
            </div>
          )}
        </Card>

        <Card title={t('home.router.title')}>
          <nav style={{ display: 'flex', gap: 'var(--spacing-md)', flexWrap: 'wrap' }}>
            <Link 
              to="/" 
              style={{ 
                color: 'var(--color-primary)', 
                textDecoration: 'none',
                padding: 'var(--spacing-xs) var(--spacing-sm)',
                borderRadius: 'var(--radius-sm)',
              }}
            >
              {t('navigation.home')}
            </Link>
            <Link 
              to="/404" 
              style={{ 
                color: 'var(--color-primary)', 
                textDecoration: 'none',
                padding: 'var(--spacing-xs) var(--spacing-sm)',
                borderRadius: 'var(--radius-sm)',
              }}
            >
              {t('home.router.test404')}
            </Link>
          </nav>
        </Card>
      </div>
    </div>
  )
}

export default HomePage

