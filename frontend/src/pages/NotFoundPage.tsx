import { Link } from 'react-router-dom'
import { Button } from '../components/ui/Button'

function NotFoundPage() {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '60vh',
      textAlign: 'center',
    }}>
      <h1 style={{ fontSize: '6rem', margin: 0, color: 'var(--color-primary)' }}>404</h1>
      <h2 style={{ color: 'var(--color-text)' }}>Страница не найдена</h2>
      <p style={{ color: 'var(--color-text-secondary)', marginBottom: '2rem' }}>
        Запрашиваемая страница не существует или была перемещена.
      </p>
      <Link to="/" style={{ textDecoration: 'none' }}>
        <Button variant="primary">Вернуться на главную</Button>
      </Link>
    </div>
  )
}

export default NotFoundPage

