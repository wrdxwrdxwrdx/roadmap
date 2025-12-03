import { useState, FormEvent, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useApi } from '../hooks/useApi'
import { apiEndpoints, LoginRequest, LoginResponse } from '../services/apiEndpoints'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'

interface LoginFormData {
  email: string
  password: string
}

function LoginPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
  })
  const [errors, setErrors] = useState<Partial<Record<keyof LoginFormData, string>>>({})
  const [touched, setTouched] = useState<Partial<Record<keyof LoginFormData, boolean>>>({})

  // Используем useApi для логина
  const { loading, error, execute, reset } = useApi<LoginResponse>(
    (loginData: LoginRequest) => apiEndpoints.login(loginData),
    {
      onSuccess: (response) => {
        // Сохраняем токен
        localStorage.setItem('token', response.token)
        // Перенаправляем на профиль
        navigate('/profile')
      },
    }
  )

  const validateField = (name: keyof LoginFormData, value: string): string => {
    switch (name) {
      case 'email':
        if (!value) return t('login.email.required')
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return t('login.email.invalid')
        return ''
      case 'password':
        if (!value) return t('login.password.required')
        return ''
      default:
        return ''
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
    
    // Валидация при изменении (если поле уже было тронуто)
    if (touched[name as keyof LoginFormData]) {
      const error = validateField(name as keyof LoginFormData, value)
      setErrors(prev => ({ ...prev, [name]: error }))
    }
  }

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setTouched(prev => ({ ...prev, [name]: true }))
    const error = validateField(name as keyof LoginFormData, value)
    setErrors(prev => ({ ...prev, [name]: error }))
  }

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof LoginFormData, string>> = {}
    let isValid = true

    Object.keys(formData).forEach(key => {
      const fieldName = key as keyof LoginFormData
      const error = validateField(fieldName, formData[fieldName])
      if (error) {
        newErrors[fieldName] = error
        isValid = false
      }
    })

    setErrors(newErrors)
    setTouched({
      email: true,
      password: true,
    })

    return isValid
  }

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    
    // Сбрасываем предыдущие ошибки
    reset()
    
    if (!validateForm()) {
      return
    }

    try {
      await execute({
        email: formData.email,
        password: formData.password,
      })
    } catch (err) {
      // Ошибка уже обработана в useApi
      console.error('Login error:', err)
    }
  }

  const getFieldError = (fieldName: keyof LoginFormData): string | undefined => {
    return touched[fieldName] ? errors[fieldName] : undefined
  }

  const isFormValid = Object.values(errors).every(err => !err) && 
                      Object.values(formData).every(val => val !== '')

  // Проверяем, не авторизован ли уже пользователь
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      navigate('/profile')
    }
  }, [navigate])

  return (
    <div style={{
      maxWidth: '500px',
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
          {t('login.title')}
        </h1>

        <form onSubmit={handleSubmit} style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 'var(--spacing-md)',
        }}>
          {/* Email */}
          <div>
            <label htmlFor="email" style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text)',
              fontWeight: 500,
            }}>
              {t('login.email.label')} *
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              onBlur={handleBlur}
              style={{
                width: '100%',
                padding: 'var(--spacing-sm)',
                fontSize: '1rem',
                border: `1px solid ${getFieldError('email') ? 'var(--color-error)' : 'var(--color-border)'}`,
                borderRadius: 'var(--radius-md)',
                backgroundColor: 'var(--color-background)',
                color: 'var(--color-text)',
                boxSizing: 'border-box',
              }}
              placeholder={t('login.email.placeholder')}
              disabled={loading}
            />
            {getFieldError('email') && (
              <p style={{
                margin: 'var(--spacing-xs) 0 0 0',
                color: 'var(--color-error)',
                fontSize: '0.875rem',
              }}>
                {getFieldError('email')}
              </p>
            )}
          </div>

          {/* Password */}
          <div>
            <label htmlFor="password" style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text)',
              fontWeight: 500,
            }}>
              {t('login.password.label')} *
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              onBlur={handleBlur}
              style={{
                width: '100%',
                padding: 'var(--spacing-sm)',
                fontSize: '1rem',
                border: `1px solid ${getFieldError('password') ? 'var(--color-error)' : 'var(--color-border)'}`,
                borderRadius: 'var(--radius-md)',
                backgroundColor: 'var(--color-background)',
                color: 'var(--color-text)',
                boxSizing: 'border-box',
              }}
              placeholder={t('login.password.placeholder')}
              disabled={loading}
            />
            {getFieldError('password') && (
              <p style={{
                margin: 'var(--spacing-xs) 0 0 0',
                color: 'var(--color-error)',
                fontSize: '0.875rem',
              }}>
                {getFieldError('password')}
              </p>
            )}
          </div>

          {/* Error message from API */}
          {error && !loading && (
            <div style={{
              padding: 'var(--spacing-md)',
              backgroundColor: 'rgba(220, 53, 69, 0.15)',
              border: '2px solid #dc3545',
              borderRadius: 'var(--radius-md)',
              color: '#dc3545',
              fontSize: '0.9rem',
              fontWeight: 500,
              marginTop: 'var(--spacing-sm)',
            }}>
              <strong style={{ display: 'block', marginBottom: 'var(--spacing-xs)' }}>
                {t('login.error')}
              </strong>
              {(() => {
                if (error.response?.data) {
                  const errorData = error.response.data as any
                  if (errorData.message) return errorData.message
                  if (typeof errorData === 'string') return errorData
                  if (errorData.error) return errorData.error
                }
                return error.message || t('login.error')
              })()}
            </div>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            variant="primary"
            size="lg"
            disabled={loading || !isFormValid}
            style={{
              width: '100%',
              marginTop: 'var(--spacing-sm)',
            }}
          >
            {loading ? t('login.submitting') : t('login.submit')}
          </Button>
        </form>

        {/* Link to register */}
        <div style={{
          marginTop: 'var(--spacing-lg)',
          textAlign: 'center',
          color: 'var(--color-text-secondary)',
          fontSize: '0.9rem',
        }}>
          {t('login.noAccount')}{' '}
          <Link
            to="/register"
            style={{
              color: 'var(--color-primary)',
              textDecoration: 'none',
              fontWeight: 500,
            }}
          >
            {t('login.registerLink')}
          </Link>
        </div>
      </Card>
    </div>
  )
}

export default LoginPage

