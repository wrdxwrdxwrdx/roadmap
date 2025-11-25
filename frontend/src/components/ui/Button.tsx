import { ButtonHTMLAttributes, ReactNode } from 'react'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'success' | 'error' | 'warning' | 'info'
  size?: 'sm' | 'md' | 'lg'
  children: ReactNode
}

/**
 * Переиспользуемый компонент кнопки
 */
export function Button({
  variant = 'primary',
  size = 'md',
  children,
  style,
  ...props
}: ButtonProps) {
  const variantStyles = {
    primary: {
      backgroundColor: 'var(--color-primary)',
      color: '#ffffff',
    },
    secondary: {
      backgroundColor: 'var(--color-secondary)',
      color: '#ffffff',
    },
    success: {
      backgroundColor: 'var(--color-success)',
      color: '#ffffff',
    },
    error: {
      backgroundColor: 'var(--color-error)',
      color: '#ffffff',
    },
    warning: {
      backgroundColor: 'var(--color-warning)',
      color: '#000000',
    },
    info: {
      backgroundColor: 'var(--color-info)',
      color: '#ffffff',
    },
  }

  const sizeStyles = {
    sm: {
      padding: 'var(--spacing-xs) var(--spacing-sm)',
      fontSize: '0.875rem',
    },
    md: {
      padding: 'var(--spacing-sm) var(--spacing-md)',
      fontSize: '1rem',
    },
    lg: {
      padding: 'var(--spacing-md) var(--spacing-lg)',
      fontSize: '1.125rem',
    },
  }

  return (
    <button
      style={{
        ...variantStyles[variant],
        ...sizeStyles[size],
        border: 'none',
        borderRadius: 'var(--radius-md)',
        cursor: props.disabled ? 'not-allowed' : 'pointer',
        opacity: props.disabled ? 0.6 : 1,
        transition: 'opacity 0.2s ease, transform 0.1s ease',
        ...style,
      }}
      onMouseEnter={(e) => {
        if (!props.disabled) {
          e.currentTarget.style.opacity = '0.9'
          e.currentTarget.style.transform = 'scale(0.98)'
        }
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.opacity = props.disabled ? '0.6' : '1'
        e.currentTarget.style.transform = 'scale(1)'
      }}
      {...props}
    >
      {children}
    </button>
  )
}

