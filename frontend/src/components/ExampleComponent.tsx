import { useAppStore } from '../store'
import { useApi } from '../hooks/useApi'
import { apiEndpoints } from '../services/apiEndpoints'

// Пример компонента, демонстрирующий использование всех настроенных инструментов
function ExampleComponent() {
  // Использование Zustand store
  const { theme, setTheme } = useAppStore()
  
  // Использование кастомного хука для API
  const { data, loading, error, execute } = useApi(apiEndpoints.healthCheck, {
    onSuccess: (data) => {
      console.log('Health check successful:', data)
    },
    onError: (error) => {
      console.error('Health check failed:', error)
    },
  })

  const handleThemeToggle = () => {
    setTheme(theme === 'light' ? 'dark' : 'light')
  }

  return (
    <div>
      <h2>Example Component</h2>
      <p>Current theme: {theme}</p>
      <button onClick={handleThemeToggle}>Toggle Theme</button>
      
      <div>
        <button onClick={() => execute()}>Check API Health</button>
        {loading && <p>Loading...</p>}
        {error && <p>Error: {error.message}</p>}
        {data && <p>API Response: {JSON.stringify(data)}</p>}
      </div>
    </div>
  )
}

export default ExampleComponent

