import { ReactNode } from 'react'

interface CardProps {
  children: ReactNode
  title?: string
  style?: React.CSSProperties
  className?: string
}

/**
 * Переиспользуемый компонент карточки
 */
export function Card({ children, title, style, className }: CardProps) {
  return (
    <div
      className={className}
      style={{
        backgroundColor: 'var(--color-surface)',
        border: '1px solid var(--color-border)',
        borderRadius: 'var(--radius-lg)',
        padding: 'var(--spacing-lg)',
        boxShadow: 'var(--shadow-md)',
        ...style,
      }}
    >
      {title && (
        <h3 style={{
          marginTop: 0,
          marginBottom: 'var(--spacing-md)',
          color: 'var(--color-text)',
        }}>
          {title}
        </h3>
      )}
      {children}
    </div>
  )
}

