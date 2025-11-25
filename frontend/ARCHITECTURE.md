# Архитектура Frontend проекта

## Структура проекта

```
frontend/src/
├── components/          # Переиспользуемые компоненты
│   ├── ui/             # UI компоненты (Button, Card и т.д.)
│   └── layouts/        # Layout компоненты (MainLayout и т.д.)
├── pages/              # Страницы приложения
├── hooks/              # Кастомные React хуки
├── services/           # API клиент и endpoints
├── store/              # Zustand stores (state management)
├── theme/              # Система тем
│   ├── theme.ts        # Конфигурация тем
│   └── useTheme.ts     # Хук для работы с темами
├── styles/             # Глобальные стили
│   └── themes.css      # CSS переменные для тем
├── types/              # TypeScript типы
├── utils/              # Утилиты
├── App.tsx             # Главный компонент приложения
└── main.tsx            # Точка входа
```

## Система тем

### Описание

Проект поддерживает светлую и темную темы. Система тем построена на CSS переменных и автоматически применяется ко всему приложению.

### Использование

1. **Применение темы:**
   - Тема применяется через хук `useTheme()` в `App.tsx`
   - Тема синхронизирована с Zustand store
   - Изменения сохраняются в localStorage
   - Применяется к DOM через CSS классы и переменные

2. **Переключение темы:**
   ```tsx
   import { useTheme } from '../theme/useTheme'
   
   function Component() {
     const { theme, toggleTheme, setTheme } = useTheme()
     
     return (
       <button onClick={toggleTheme}>
         Текущая тема: {theme}
       </button>
     )
   }
   ```

3. **Использование CSS переменных:**
   ```css
   .my-component {
     background-color: var(--color-background);
     color: var(--color-text);
     border: 1px solid var(--color-border);
   }
   ```

### Доступные CSS переменные

- `--color-background` - основной фон
- `--color-surface` - фон для карточек и блоков
- `--color-primary` - основной цвет (кнопки, ссылки)
- `--color-text` - основной текст
- `--color-text-secondary` - вторичный текст
- `--color-text-muted` - приглушенный текст
- `--color-border` - границы
- `--color-success`, `--color-error`, `--color-warning`, `--color-info`
- И другие (см. `styles/themes.css`)

## Компоненты

### UI компоненты (`components/ui/`)

Переиспользуемые UI компоненты для всего приложения:

- **Button** - кнопка с вариантами (primary, secondary, success, error и т.д.)
- **Card** - карточка для группировки контента

### Layout компоненты (`components/layouts/`)

Компоненты для структуры страниц:

- **MainLayout** - основной layout с header, main, footer

## Добавление новых страниц

1. Создайте файл в `pages/`:
   ```tsx
   // pages/AboutPage.tsx
   import { Card } from '../components/ui/Card'
   
   export default function AboutPage() {
     return (
       <Card title="О нас">
         <p>Содержимое страницы</p>
       </Card>
     )
   }
   ```

2. Добавьте маршрут в `App.tsx`:
   ```tsx
   import AboutPage from './pages/AboutPage'
   
   <Route path="/about" element={<AboutPage />} />
   ```

## State Management

Используется Zustand для управления глобальным состоянием:

- `useAppStore` - основное состояние приложения (тема, пользователь, авторизация)

## API

- `services/api.ts` - настроенный axios клиент
- `services/apiEndpoints.ts` - типизированные endpoints
- `hooks/useApi.ts` - хук для удобной работы с API

## Масштабирование

Архитектура рассчитана на большой проект:

1. **Модульность** - четкое разделение на папки
2. **Переиспользование** - UI компоненты и layouts
3. **Типизация** - полная поддержка TypeScript
4. **Темы** - легко расширяемая система тем
5. **Масштабируемость** - легко добавлять новые страницы и компоненты

## Best Practices

1. Используйте UI компоненты вместо inline стилей где возможно
2. Применяйте CSS переменные для цветов
3. Размещайте переиспользуемые компоненты в `components/ui/`
4. Размещайте специфичные для страницы компоненты рядом со страницей
5. Используйте типы из `types/` для TypeScript
6. Используйте `useTheme()` для работы с темами

