import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Button } from '../components/ui/Button'

function NotFoundPage() {
  const { t } = useTranslation()
  
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '60vh',
      textAlign: 'center',
    }}>
      <h1 style={{ fontSize: '6rem', margin: 0, color: 'var(--color-primary)' }}>{t('notFound.title')}</h1>
      <h2 style={{ color: 'var(--color-text)' }}>{t('notFound.message')}</h2>
      <p style={{ color: 'var(--color-text-secondary)', marginBottom: '2rem' }}>
        {t('notFound.description')}
      </p>
      <Link to="/" style={{ textDecoration: 'none' }}>
        <Button variant="primary">{t('notFound.goHome')}</Button>
      </Link>
    </div>
  )
}

export default NotFoundPage

