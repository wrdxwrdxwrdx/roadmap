import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useApi } from '../hooks/useApi'
import { apiEndpoints, ProfileResponse } from '../services/apiEndpoints'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'

function ProfilePage() {
  const { t } = useTranslation()
  const navigate = useNavigate()

  // Загружаем профиль пользователя
  const { data: profile, loading, error, execute } = useApi<ProfileResponse>(
    apiEndpoints.getProfile,
    {
      immediate: true,
      onError: () => {
        // Если ошибка авторизации, перенаправляем на логин
        navigate('/login')
      },
    }
  )

  const handleLogout = () => {
    localStorage.removeItem('token')
    navigate('/login')
  }

  // Проверяем наличие токена при загрузке
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate('/login')
    }
  }, [navigate])

  if (loading) {
    return (
      <div style={{
        maxWidth: '600px',
        margin: '0 auto',
        padding: 'var(--spacing-lg)',
      }}>
        <Card>
          <p style={{ textAlign: 'center', color: 'var(--color-text-secondary)' }}>
            {t('profile.loading')}
          </p>
        </Card>
      </div>
    )
  }

  if (error || !profile) {
    return (
      <div style={{
        maxWidth: '600px',
        margin: '0 auto',
        padding: 'var(--spacing-lg)',
      }}>
        <Card>
          <div style={{
            padding: 'var(--spacing-md)',
            backgroundColor: 'rgba(220, 53, 69, 0.15)',
            border: '2px solid #dc3545',
            borderRadius: 'var(--radius-md)',
            color: '#dc3545',
            fontSize: '0.9rem',
            fontWeight: 500,
            marginBottom: 'var(--spacing-md)',
          }}>
            <strong style={{ display: 'block', marginBottom: 'var(--spacing-xs)' }}>
              {t('profile.error')}
            </strong>
            {error?.message || t('profile.errorMessage')}
          </div>
          <Button
            variant="primary"
            onClick={() => execute()}
            style={{ width: '100%' }}
          >
            {t('profile.retry')}
          </Button>
        </Card>
      </div>
    )
  }

  return (
    <div style={{
      maxWidth: '600px',
      margin: '0 auto',
      padding: 'var(--spacing-lg)',
    }}>
      <Card>
        <h1 style={{
          marginTop: 0,
          marginBottom: 'var(--spacing-lg)',
          textAlign: 'center',
          color: 'var(--color-text)',
        }}>
          {t('profile.title')}
        </h1>

        <div style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 'var(--spacing-md)',
        }}>
          {/* User ID */}
          <div>
            <label style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text-secondary)',
              fontSize: '0.9rem',
              fontWeight: 500,
            }}>
              {t('profile.userId')}
            </label>
            <div style={{
              padding: 'var(--spacing-sm)',
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
              borderRadius: 'var(--radius-md)',
              color: 'var(--color-text)',
              fontSize: '1rem',
            }}>
              {profile.user_id}
            </div>
          </div>

          {/* Username */}
          <div>
            <label style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text-secondary)',
              fontSize: '0.9rem',
              fontWeight: 500,
            }}>
              {t('profile.username')}
            </label>
            <div style={{
              padding: 'var(--spacing-sm)',
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
              borderRadius: 'var(--radius-md)',
              color: 'var(--color-text)',
              fontSize: '1rem',
            }}>
              {profile.username}
            </div>
          </div>

          {/* Email */}
          <div>
            <label style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text-secondary)',
              fontSize: '0.9rem',
              fontWeight: 500,
            }}>
              {t('profile.email')}
            </label>
            <div style={{
              padding: 'var(--spacing-sm)',
              backgroundColor: 'var(--color-surface)',
              border: '1px solid var(--color-border)',
              borderRadius: 'var(--radius-md)',
              color: 'var(--color-text)',
              fontSize: '1rem',
            }}>
              {profile.email}
            </div>
          </div>

          {/* Logout Button */}
          <Button
            variant="primary"
            onClick={handleLogout}
            style={{
              width: '100%',
              marginTop: 'var(--spacing-md)',
            }}
          >
            {t('navigation.logout')}
          </Button>
        </div>
      </Card>
    </div>
  )
}

export default ProfilePage

