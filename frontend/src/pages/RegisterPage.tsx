import { useState, FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useApi } from '../hooks/useApi'
import { apiEndpoints } from '../services/apiEndpoints'
import { Card } from '../components/ui/Card'
import { Button } from '../components/ui/Button'

interface RegisterFormData {
  email: string
  username: string
  password: string
  confirmPassword: string
}

interface RegisterResponse {
  id: string
  username: string
  email: string
  token: string
  created_at: string
  updated_at: string
}

function RegisterPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [formData, setFormData] = useState<RegisterFormData>({
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
  })
  const [errors, setErrors] = useState<Partial<Record<keyof RegisterFormData, string>>>({})
  const [touched, setTouched] = useState<Partial<Record<keyof RegisterFormData, boolean>>>({})

  // Используем useApi для регистрации
  const { data, loading, error, execute, reset } = useApi<RegisterResponse>(
    (registerData: { email: string; username: string; password: string }) => 
      apiEndpoints.register(registerData),
    {
      onSuccess: (response) => {
        // Сохраняем токен
        localStorage.setItem('token', response.token)
        // Перенаправляем на профиль
        navigate('/profile')
      },
    }
  )

  const validateField = (name: keyof RegisterFormData, value: string): string => {
    switch (name) {
      case 'email':
        if (!value) return t('register.email.required')
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return t('register.email.invalid')
        return ''
      case 'username':
        if (!value) return t('register.username.required')
        if (value.length < 3) return t('register.username.minLength')
        if (value.length > 100) return t('register.username.maxLength')
        if (!/^[a-zA-Z0-9_]+$/.test(value)) return t('register.username.invalid')
        return ''
      case 'password':
        if (!value) return t('register.password.required')
        if (value.length < 8) return t('register.password.minLength')
        return ''
      case 'confirmPassword':
        if (!value) return t('register.confirmPassword.required')
        if (value !== formData.password) return t('register.confirmPassword.mismatch')
        return ''
      default:
        return ''
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
    
    // Валидация при изменении (если поле уже было тронуто)
    if (touched[name as keyof RegisterFormData]) {
      const error = validateField(name as keyof RegisterFormData, value)
      setErrors(prev => ({ ...prev, [name]: error }))
    }
  }

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setTouched(prev => ({ ...prev, [name]: true }))
    const error = validateField(name as keyof RegisterFormData, value)
    setErrors(prev => ({ ...prev, [name]: error }))
  }

  const validateForm = (): boolean => {
    const newErrors: Partial<Record<keyof RegisterFormData, string>> = {}
    let isValid = true

    Object.keys(formData).forEach(key => {
      const fieldName = key as keyof RegisterFormData
      const error = validateField(fieldName, formData[fieldName])
      if (error) {
        newErrors[fieldName] = error
        isValid = false
      }
    })

    setErrors(newErrors)
    setTouched({
      email: true,
      username: true,
      password: true,
      confirmPassword: true,
    })

    return isValid
  }

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    
    // Сбрасываем предыдущие ошибки и успешные сообщения
    reset()
    
    if (!validateForm()) {
      return
    }

    try {
      await execute({
        email: formData.email,
        username: formData.username,
        password: formData.password,
      })
    } catch (err) {
      // Ошибка уже обработана в useApi
      console.error('Registration error:', err)
    }
  }

  const getFieldError = (fieldName: keyof RegisterFormData): string | undefined => {
    return touched[fieldName] ? errors[fieldName] : undefined
  }

  const isFormValid = Object.values(errors).every(err => !err) && 
                      Object.values(formData).every(val => val !== '')

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
          {t('register.title')}
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
              {t('register.email.label')} *
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
              placeholder={t('register.email.placeholder')}
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

          {/* Username */}
          <div>
            <label htmlFor="username" style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text)',
              fontWeight: 500,
            }}>
              {t('register.username.label')} *
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={formData.username}
              onChange={handleChange}
              onBlur={handleBlur}
              style={{
                width: '100%',
                padding: 'var(--spacing-sm)',
                fontSize: '1rem',
                border: `1px solid ${getFieldError('username') ? 'var(--color-error)' : 'var(--color-border)'}`,
                borderRadius: 'var(--radius-md)',
                backgroundColor: 'var(--color-background)',
                color: 'var(--color-text)',
                boxSizing: 'border-box',
              }}
              placeholder={t('register.username.placeholder')}
              disabled={loading}
            />
            {getFieldError('username') && (
              <p style={{
                margin: 'var(--spacing-xs) 0 0 0',
                color: 'var(--color-error)',
                fontSize: '0.875rem',
              }}>
                {getFieldError('username')}
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
              {t('register.password.label')} *
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
              placeholder={t('register.password.placeholder')}
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

          {/* Confirm Password */}
          <div>
            <label htmlFor="confirmPassword" style={{
              display: 'block',
              marginBottom: 'var(--spacing-xs)',
              color: 'var(--color-text)',
              fontWeight: 500,
            }}>
              {t('register.confirmPassword.label')} *
            </label>
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              onBlur={handleBlur}
              style={{
                width: '100%',
                padding: 'var(--spacing-sm)',
                fontSize: '1rem',
                border: `1px solid ${getFieldError('confirmPassword') ? 'var(--color-error)' : 'var(--color-border)'}`,
                borderRadius: 'var(--radius-md)',
                backgroundColor: 'var(--color-background)',
                color: 'var(--color-text)',
                boxSizing: 'border-box',
              }}
              placeholder={t('register.confirmPassword.placeholder')}
              disabled={loading}
            />
            {getFieldError('confirmPassword') && (
              <p style={{
                margin: 'var(--spacing-xs) 0 0 0',
                color: 'var(--color-error)',
                fontSize: '0.875rem',
              }}>
                {getFieldError('confirmPassword')}
              </p>
            )}
          </div>

          {/* Error message from API - красное окно */}
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
                {t('register.error')}
              </strong>
              {(() => {
                if (error.response?.data) {
                  const errorData = error.response.data as any
                  if (errorData.message) return errorData.message
                  if (typeof errorData === 'string') return errorData
                  if (errorData.error) return errorData.error
                }
                return error.message || t('register.error')
              })()}
            </div>
          )}

          {/* Success message - зеленое окно */}
          {data && !loading && !error && (
            <div style={{
              padding: 'var(--spacing-md)',
              backgroundColor: 'rgba(40, 167, 69, 0.15)',
              border: '2px solid #28a745',
              borderRadius: 'var(--radius-md)',
              color: '#28a745',
              fontSize: '0.9rem',
              fontWeight: 500,
              marginTop: 'var(--spacing-sm)',
            }}>
              <strong style={{ display: 'block', marginBottom: 'var(--spacing-xs)' }}>
                {t('register.success')}
              </strong>
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
            {loading ? t('register.submitting') : t('register.submit')}
          </Button>
        </form>

        {/* Link to login */}
        <div style={{
          marginTop: 'var(--spacing-lg)',
          textAlign: 'center',
          color: 'var(--color-text-secondary)',
          fontSize: '0.9rem',
        }}>
          {t('register.hasAccount')}{' '}
          <Link
            to="/login"
            style={{
              color: 'var(--color-primary)',
              textDecoration: 'none',
              fontWeight: 500,
            }}
          >
            {t('register.loginLink')}
          </Link>
        </div>
      </Card>
    </div>
  )
}

export default RegisterPage

