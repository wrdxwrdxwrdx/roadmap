import { useState, useEffect } from 'react'
import { AxiosResponse, AxiosError } from 'axios'

// Кастомный хук для работы с API
interface UseApiOptions<T> {
  immediate?: boolean // Выполнить запрос сразу
  onSuccess?: (data: T) => void
  onError?: (error: AxiosError) => void
}

interface UseApiResult<T> {
  data: T | null
  loading: boolean
  error: AxiosError | null
  execute: (...args: any[]) => Promise<void>
  reset: () => void
}

export function useApi<T>(
  apiFunction: (...args: any[]) => Promise<AxiosResponse<T>>,
  options: UseApiOptions<T> = {}
): UseApiResult<T> {
  const { immediate = false, onSuccess, onError } = options
  
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState<boolean>(immediate)
  const [error, setError] = useState<AxiosError | null>(null)

  const execute = async (...args: any[]) => {
    try {
      setLoading(true)
      setError(null)
      const response = await apiFunction(...args)
      setData(response.data)
      onSuccess?.(response.data)
    } catch (err) {
      const axiosError = err as AxiosError
      setError(axiosError)
      onError?.(axiosError)
    } finally {
      setLoading(false)
    }
  }

  const reset = () => {
    setData(null)
    setError(null)
    setLoading(false)
  }

  useEffect(() => {
    if (immediate) {
      execute()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { data, loading, error, execute, reset }
}

